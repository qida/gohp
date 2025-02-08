package gormx

import (
	"fmt"

	"gorm.io/gorm"
)

func Paginate(page, limit int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Offset(int(offset)).Limit(int(limit))
	}
}

// 拼接模糊查询
func Wildcard[T any](key T) string {
	return fmt.Sprintf("%%%v%%", key)
}
