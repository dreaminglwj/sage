package repository

import (
	"xorm.io/xorm"
)

type BaseRepository[T any] struct{}

// Get 获取模型
func (r *BaseRepository[T]) Get(db *xorm.Session, id string) (*T, error) {
	var t T
	has, err := db.ID(id).Get(&t)
	if err != nil {
		return nil, err
	}
	if has {
		return &t, nil
	}
	return nil, nil
}

// Update 更新模型
func (r *BaseRepository[T]) Update(db *xorm.Session, id string, model *T) (rowsAffected int64, err error) {
	rowsAffected, err = db.ID(id).Update(model)
	return
}

// UpdateByMap 按需更新模型
func (r *BaseRepository[T]) UpdateByMap(db *xorm.Session, id string, data map[string]any) (rowsAffected int64, err error) {
	rowsAffected, err = db.Table(new(T)).ID(id).Update(data)
	return
}

// Insert 新增
func (r *BaseRepository[T]) Insert(db *xorm.Session, model *T) (rowsAffected int64, err error) {
	rowsAffected, err = db.Insert(model)
	return
}

// InsertAll 新增
func (r *BaseRepository[T]) InsertAll(db *xorm.Session, models []*T) (rowsAffected int64, err error) {
	rowsAffected, err = db.Insert(models)
	return
}

// Delete 删除
func (r *BaseRepository[T]) Delete(db *xorm.Session, id string) (rowsAffected int64, err error) {
	rowsAffected, err = db.ID(id).Delete(new(T))
	return
}
