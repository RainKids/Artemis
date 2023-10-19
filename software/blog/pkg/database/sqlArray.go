package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
)

//type Array[T string | int32 | int8] []T func (a *Array[T])

type Array []any

type Tes struct {
	a Array
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (a *Array) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to scan Array value:", value))
	}
	if len(bytes) > 0 {
		return json.Unmarshal(bytes, a)
	}
	*a = make([]any, 0)
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (a Array) Value() (driver.Value, error) {
	if a == nil {
		return "[]", nil
	}
	return convertor.ToString(a), nil
}
