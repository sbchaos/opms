package internal

import (
	"fmt"

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

	switch r.(type) {
	case *sqldriver.NullInt8:
		v, ok := r.(*sqldriver.NullInt8)
		if !ok {
			return ForceString(v)
		}
		return fmt.Sprintf("%d", v.Int8)

	case *sqldriver.NullInt16:
		v, ok := r.(*sqldriver.NullInt16)
		if !ok {
			return ForceString(v)
		}
		return fmt.Sprintf("%d", v.Int16)

	case *sqldriver.NullInt32:
		v, ok := r.(*sqldriver.NullInt32)
		if !ok {
			return ForceString(v)
		}
		return fmt.Sprintf("%d", v.Int32)

	case *sqldriver.NullInt64:
		v, ok := r.(*sqldriver.NullInt64)
		if !ok {
			return ForceString(v)
		}
		return fmt.Sprintf("%d", v.Int64)

	case *sqldriver.NullFloat32:
		v, ok := r.(*sqldriver.NullFloat32)
		if !ok {
			return ForceString(v)
		}
		return fmt.Sprintf("%f", v.Float32)

	case *sqldriver.NullFloat64:
		v, ok := r.(*sqldriver.NullFloat64)
		if !ok {
			return ForceString(v)
		}
		return fmt.Sprintf("%f", v.Float64)

	case *sqldriver.NullString:
		v, ok := r.(*sqldriver.NullString)
		if !ok {
			return ForceString(v)
		}
		return fmt.Sprintf("%s", v.String)

	case *sqldriver.NullBool:
		v, ok := r.(*sqldriver.NullBool)
		if !ok {
			return ForceString(v)
		}
		return fmt.Sprintf("%v", v.Bool)

	default:
		s, ok1 := r.(string)
		if !ok1 {
			return ForceString(r)
		}
		return s
	}
}

func ForceString(v any) string {
	return fmt.Sprintf("%v", v)
}
