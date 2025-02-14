package parse

import (
	"errors"
	"strconv"
	"unicode"

	h "github.com/sbchaos/consume/comb"
	b "github.com/sbchaos/consume/par"
	"github.com/sbchaos/consume/par/char"
	sp "github.com/sbchaos/consume/par/strings"

	"github.com/sbchaos/opms/cmd/optimus/internal/resource"
)

var (
	ErrNotReq = errors.New("not required")
)

type DDLParser struct {
	l  b.Logger
	p1 b.Parser[rune, resource.ExternalTable]
}

func NewDDLParser(l b.Logger, required map[string]string) (*DDLParser, error) {
	p1 := ddlParser(l, required)
	return &DDLParser{
		l:  l,
		p1: p1,
	}, nil
}

func (d *DDLParser) ParseExternalTable(content string) (*resource.ExternalTable, error) {
	ext, err := b.ParseString(content, d.p1)
	if err != nil {
		return &ext, err
	}
	return &ext, nil
}

func StringWithOptionalQuotes(l b.Logger) b.Parser[rune, string] {
	return b.Trace(l, "stringWithOptionalQuotes", h.Choice(
		sp.QuotedString(0, sp.Quotes...),
		sp.CustomString(func(a rune) bool {
			return unicode.IsLetter(a) || unicode.IsDigit(a) || a == '.' || a == '_' || a == '-'
		}),
	))
}

func SpaceAndString(l b.Logger) b.Parser[rune, string] {
	return b.Trace(l, "SpaceAndString", h.Skip(char.WhiteSpaces(), StringWithOptionalQuotes(l)))
}

func OptValue(l b.Logger) b.Parser[rune, string] {
	return h.Choice(
		h.Between(sp.String("[(", sp.Equals), sp.CustomString(func(a rune) bool {
			return a != ']'
		}), sp.String("]", sp.Equals)),
		h.Between(char.Single('['), SpaceAndString(l), char.Single(']')),
		SpaceAndString(l),
	)
}

func ddlParser(l b.Logger, required map[string]string) b.Parser[rune, resource.ExternalTable] {
	nameAndSchema := h.And(
		nameParser(l, required),
		h.Optional(schemaParser(l), resource.Schema{}),
		func(a string, sch resource.Schema) resource.ExternalTable {
			return resource.ExternalTable{
				Name:   a,
				Schema: sch,
			}
		})
	return h.And(nameAndSchema, sourceParser(l), func(a resource.ExternalTable, b resource.ExternalSource) resource.ExternalTable {
		a.Source = &b
		return a
	})
}

func nameParser(l b.Logger, required map[string]string) b.Parser[rune, string] {
	header := b.Trace(l, "header", sp.Sequence([]string{"CREATE", "EXTERNAL", "TABLE"}, sp.EqualIgnoreCase))
	return h.Skip(
		header,
		h.FMap1(func(a string) (string, error) {
			if len(required) != 0 {
				if _, ok := required[a]; !ok {
					return "", ErrNotReq
				}
			}
			return a, nil
		}, b.Trace(l, "string", StringWithOptionalQuotes(l))))
}

func fieldParser(l b.Logger) b.Parser[rune, *resource.Field] {
	optionalQual := h.And(
		char.WhiteSpaces(),
		h.Optional(sp.Sequence([]string{"NOT", "NULL"}, sp.EqualIgnoreCase), ""),
		func(a rune, b string) string {
			return ""
		})
	emptyMap := map[string]string{}
	strt := h.And(
		SpaceAndString(l),
		h.SkipAfter(SpaceAndString(l), optionalQual),
		func(nm string, typ string) *resource.Field {
			return &resource.Field{
				Name: nm,
				Type: typ,
			}
		})
	return h.And(strt, h.Optional(parseOptions(l), emptyMap), func(f *resource.Field, mp map[string]string) *resource.Field {
		if v, ok := mp["description"]; ok {
			f.Description = v
		}
		return f
	})
}

func schemaParser(l b.Logger) b.Parser[rune, resource.Schema] {
	fields := h.Optional(h.SepBy(char.Single(','), b.Trace(l, "field_Parser", fieldParser(l))), []*resource.Field{})

	return b.Trace(l, "schemaParser", h.Skip(char.WhiteSpaces(), h.Between(
		char.Single('('),
		h.FMap(func(f []*resource.Field) resource.Schema {
			return f
		}, fields),
		h.Skip(char.WhiteSpaces(), char.Single(')')),
	)))
}

func parseOptions(l b.Logger) b.Parser[rune, map[string]string] {
	start := h.Skip(char.WhiteSpaces(),
		sp.Sequence([]string{"OPTIONS", "("}, sp.EqualIgnoreCase))
	parseAsMap := b.Trace(l, "mapParse", h.ToMap(SpaceAndString(l), char.Single('='), OptValue(l), char.Single(',')))
	end := h.Skip(char.WhiteSpaces(), char.Single(')'))
	return h.Between(start, parseAsMap, end)
}

func sourceParser(l b.Logger) b.Parser[rune, resource.ExternalSource] {
	return h.FMap(func(mp map[string]string) resource.ExternalSource {
		es := resource.ExternalSource{
			SourceType: mp["format"],
			SourceURIs: []string{mp["uris"]},
			Config:     map[string]interface{}{},
		}
		if v, ok := mp["sheet_range"]; ok {
			es.Config["range"] = v
		}

		if v2, ok := mp["skip_leading_rows"]; ok {
			v3, err := strconv.ParseInt(v2, 10, 64)
			if err == nil {
				es.Config["skip_leading_rows"] = v3
			}
		}

		return es
	}, parseOptions(l))
}
