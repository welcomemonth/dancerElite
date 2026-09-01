package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// StringSlice 字符串数组，以 JSON 文本形式存储，兼容 SQLite / Postgres。
// 用于活动级别组合、舞种等多选字段。
type StringSlice []string

// Value 实现 driver.Valuer，写入时序列化为 JSON 文本。
func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan 实现 sql.Scanner，读取时从 JSON 文本反序列化。
func (s *StringSlice) Scan(value any) error {
	if value == nil {
		*s = nil
		return nil
	}
	var raw []byte
	switch v := value.(type) {
	case []byte:
		raw = v
	case string:
		raw = []byte(v)
	default:
		return fmt.Errorf("unsupported StringSlice value: %T", value)
	}
	if len(raw) == 0 {
		*s = nil
		return nil
	}
	return json.Unmarshal(raw, s)
}
