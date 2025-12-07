package dao

import "seventh/database"

func Create[T any](obj *T) error {
	return database.DB.Create(obj).Error
}

func Read[T any](list *[]T) error {
	return database.DB.Find(list).Error
}

func Update[T any](id uint, values map[string]interface{}, obj *T) error {
	return database.DB.Model(obj).Where("id = ?", id).Updates(values).Error
}

func Delete[T any](id uint, obj *T) error {
	return database.DB.Delete(obj, id).Error
}
