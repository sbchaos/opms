package internal

import (
	"fmt"

	"github.com/aliyun/aliyun-odps-go-sdk/odps/data"
	"github.com/aliyun/aliyun-odps-go-sdk/sqldriver"
)

func ToString(r any) string {
	rr, ok := r.(sqldriver.NullAble)
	if !ok {
		return ForceString(r)
	}
	if rr.IsNull() {
		return "NULL"
	}

	switch r := r.(type) {
	case *sqldriver.NullInt8:
		return fmt.Sprintf("%d", r.Int8)

	case *sqldriver.NullInt16:
		return fmt.Sprintf("%d", r.Int16)

	case *sqldriver.NullInt32:
		return fmt.Sprintf("%d", r.Int32)

	case *sqldriver.NullInt64:
		return fmt.Sprintf("%d", r.Int64)

	case *sqldriver.NullFloat32:
		return fmt.Sprintf("%f", r.Float32)

	case *sqldriver.NullFloat64:
		return fmt.Sprintf("%f", r.Float64)

	case *sqldriver.NullString:
		return r.String

	case *sqldriver.NullBool:
		return fmt.Sprintf("%v", r.Bool)

	case *sqldriver.Binary:
		return r.String()

	case *data.Decimal:
		return r.String()

	case *data.Map:
		return r.String()

	case *data.Array:
		return r.String()

	case *data.Struct:
		return r.String()

	case *data.Json:
		return r.String()

	default:
		s, ok1 := r.(string)
		if !ok1 {
			return ForceString(r)
		}
		return s
	}
}

func ForceString(v any) string {
	fmt.Printf("null typecast %T %v\n", v, v)
	return fmt.Sprintf("%v", v)
}
