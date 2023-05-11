package model

// 数据集
type Schema struct {
	TableName              string
	ColumnName             string
	ColumnType             string
	DataType               string
	CharacterMaximumLength string
	IsNullable             string
	ColumnDefault          string
	ColumnComment          string
}
