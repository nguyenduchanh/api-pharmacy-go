package common

import "database/sql/driver"

type BitBool bool

// Scan implements the Scanner interface.
func (b *BitBool) Scan(value interface{}) error {
	if value == nil {
		*b = false
		return nil
	}

	if v, ok := value.([]byte); ok && len(v) > 0 {
		*b = v[0] == 1
	} else if v, ok := value.(int64); ok {
		*b = v == 1
	} else {
		*b = false
	}
	return nil
}

// Value implements the driver Valuer interface.
func (b BitBool) Value() (driver.Value, error) {
	if b {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}
