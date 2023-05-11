package service

import (
	"context"
	"github.com/dreaminglwj/sage/internal/storage/model"
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/dreaminglwj/sage/internal/conf"
	"github.com/dreaminglwj/sage/internal/storage"
	"github.com/dreaminglwj/sage/internal/storage/repository"
	sage "github.com/dreaminglwj/sage/proto/sage"
)

type SchemaService struct {
	sage.UnimplementedSchemaServer
	config  *conf.Config
	storage *storage.Storage
	logger  *log.Helper

	schemaRepository repository.SchemaRepository
}

func NewSchemaService(
	config *conf.Config,
	storage *storage.Storage,
	logger *log.Helper,

	schemaRepository repository.SchemaRepository,
) *SchemaService {
	svc := SchemaService{
		config:           config,
		storage:          storage,
		logger:           logger,
		schemaRepository: schemaRepository,
	}
	go func() {
		if err := svc.initSystemSchema(); err != nil {
			svc.logger.Error("init system schema failed :", err)
		}
	}()
	return &svc
}

func (s *SchemaService) initSystemSchema() error {
	s.logger.Infof("lwj===> %+v", s.config.DBTables)

	columnNames := []string{"id", "created_by", "updated_by", "created_at", "updated_at", "deleted_at"}
	columnNameMap := map[string]string{
		"id":         "id",
		"created_by": "created_by",
		"updated_by": "updated_by",
		"created_at": "created_at",
		"updated_at": "updated_at",
		"deleted_at": "deleted_at",
	}

	content := ""
	db := s.storage.DB(context.Background())
	for database, tables := range s.config.DBTables.Tables {
		for _, table := range tables {
			rltMap, err := s.schemaRepository.GetSchema(db, database, table)
			if err != nil {
				return err
			}

			content += table + "\n"
			content += "字段,类型,是否允许为空,默认值,备注" + "\n"

			for _, columnName := range columnNames {
				column, ok := rltMap[columnName]
				if ok {
					content += columnToString(column) + "\n"
				}
			}

			for _, rlt := range rltMap {
				colName := rlt.ColumnName
				if _, ok := columnNameMap[colName]; ok {
					continue
				}

				content += columnToString(rlt) + "\n"
			}

			content += "\n\n\n"
		}
	}

	out, err := os.Create(s.config.OutFileName + ".csv")
	if err != nil {
		s.logger.Error(err)
		return err
	}
	defer out.Close()

	_, err = out.WriteString(content)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	//s.logger.Infof("lwj===>content:%s", content)
	s.logger.Info("export finish!!!")

	return nil
}

func columnToString(column *model.Schema) string {
	colName := column.ColumnName
	colType := column.ColumnType
	colIsNull := "否"
	if strings.ToLower(column.IsNullable) == "yes" {
		colIsNull = "是"
	}
	colDefault := column.ColumnDefault
	colComment := column.ColumnComment
	rlt := colName + "," +
		colType + "," +
		colIsNull + "," +
		colDefault + "," +
		colComment
	return rlt
}

func (s *SchemaService) GetSchema(ctx context.Context, req *sage.GetSchemaRequest) (*sage.GetSchemaReply, error) {
	return nil, nil
}
