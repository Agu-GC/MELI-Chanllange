package models

type TableSchema struct {
	Name    string
	Columns []ColumnSchema
}

type ColumnSchema struct {
	Name     string
	Type     string
	Nullable bool
	Key      string
	Default  string
	Extra    string
}
