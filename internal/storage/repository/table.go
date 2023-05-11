package repository

import (
	//"github.com/dreaminglwj/sage/internal/storage/repository"
	"xorm.io/xorm"

	"github.com/dreaminglwj/sage/internal/storage/model"
)

type SchemaRepository interface {
	Insert(db *xorm.Session, schema *model.Schema) (rowsAffected int64, err error)
	Update(db *xorm.Session, id string, schema *model.Schema) (rowsAffected int64, err error)
	UpdateByMap(db *xorm.Session, id string, data map[string]any) (rowsAffected int64, err error)
	Delete(db *xorm.Session, id string) (rowsAffected int64, err error)
	Get(db *xorm.Session, id string) (*model.Schema, error)
	GetSchema(db *xorm.Session, dbName, tableName string) (map[string]*model.Schema, error)
}

var _ SchemaRepository = (*schemaRepository)(nil)

func NewSchemaRepository() *schemaRepository {
	return &schemaRepository{
		BaseRepository: BaseRepository[model.Schema]{},
	}
}

type schemaRepository struct {
	BaseRepository[model.Schema]
}

func (r *schemaRepository) GetSchema(db *xorm.Session, dbName, tableName string) (map[string]*model.Schema, error) {
	sql := "select COLUMN_NAME , COLUMN_TYPE , DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, IS_NULLABLE , COLUMN_DEFAULT ,COLUMN_COMMENT " +
		"FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema = ? and table_name = ?"

	rlts, err := db.Query(sql, dbName, tableName)
	if err != nil {
		return nil, err
	}

	schemaMap := make(map[string]*model.Schema, 0)

	for _, rlt := range rlts {
		sc := &model.Schema{
			TableName:              tableName,
			ColumnName:             string(rlt["COLUMN_NAME"]),
			ColumnType:             string(rlt["COLUMN_TYPE"]),
			DataType:               string(rlt["DATA_TYPE"]),
			CharacterMaximumLength: string(rlt["CHARACTER_MAXIMUM_LENGTH"]),
			IsNullable:             string(rlt["IS_NULLABLE"]),
			ColumnDefault:          string(rlt["COLUMN_DEFAULT"]),
			ColumnComment:          string(rlt["COLUMN_COMMENT"]),
		}
		schemaMap[sc.ColumnName] = sc
	}

	//fmt.Printf("%+v", rlts)

	return schemaMap, nil
}
