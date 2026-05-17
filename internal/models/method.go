package models

type Method struct {
	ID          int64  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	CreatedAt   string `db:"created_at"`
}

// MethodWithUsage extends Method with a derived usage count,
// populated by a JOIN query in the store rather than stored in the DB.
type MethodWithUsage struct {
	Method
	UsageCount int `db:"usage_count"`
}
