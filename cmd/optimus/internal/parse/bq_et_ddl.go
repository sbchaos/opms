package parse

import (
	"strconv"
	"unicode"

	h "github.com/sbchaos/consume/comb"
	b "github.com/sbchaos/consume/par"
	"github.com/sbchaos/consume/par/char"
	sp "github.com/sbchaos/consume/par/strings"
	"github.com/sbchaos/consume/stream"

	"github.com/sbchaos/opms/cmd/optimus/internal/resource"
)

type DDLParser struct {
	l  b.Logger
	p1 b.Parser[rune, resource.ExternalTable]
}

func NewDDLParser(l b.Logger, fn func(name string) error) (*DDLParser, error) {
	if l == nil {
		l = &b.FmtLog{}
	}

	p1 := ddlParser(l, fn)
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

func ddlParser(l b.Logger, fn func(string) error) b.Parser[rune, resource.ExternalTable] {
	return func(ss stream.SimpleStream[rune]) (resource.ExternalTable, error) {
		name, err := b.Parse(ss, nameParser(l))
		if err != nil {
			return resource.ExternalTable{}, err
		}

		err = fn(name)
		if err != nil {
			return resource.ExternalTable{
				Name: name,
			}, err
		}

		schema, err := b.Parse(ss, schemaParser(l))
		if err != nil {
			schema = resource.Schema{}
		}

		et := resource.ExternalTable{
			Name:   name,
			Schema: schema,
		}

		source, err := b.Parse(ss, sourceParser(l))
		if err != nil {
			return et, err
		}

		et.Source = &source
		return et, nil
	}
}

func nameParser(l b.Logger) b.Parser[rune, string] {
	header := b.Trace(l, "header", sp.Sequence([]string{"CREATE", "EXTERNAL", "TABLE"}, sp.EqualIgnoreCase))
	return h.Skip(
		header,
		b.Trace(l, "string", StringWithOptionalQuotes(l)))
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
