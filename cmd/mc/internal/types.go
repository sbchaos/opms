package internal

import (
	"fmt"

	"github.com/aliyun/aliyun-odps-go-sdk/sqldriver"
)

func ToString(r any) string {
	rr, ok := r.(sqldriver.NullAble)
	if !ok {
		return fmt.Sprintf("%v", r)
	}
	if rr.IsNull() {
		return "NULL"
	}

	return fmt.Sprintf("%v", r)
}
