package dao

import (
	"encoding/json"
	"fmt"
	"time"

	"eighth/database"
	"eighth/model"
)

func todoIDsKey() string {
	return "todo:ids"
}

func todoItemKey(id uint) string {
	return fmt.Sprintf("todo:item:%d", id)
}

func GetTodoList() ([]model.Todo, error) {
	ids, err := database.RDB.SMembers(database.Ctx, todoIDsKey()).Result()

	var todoIDs []uint
	if err == nil && len(ids) > 0 {
		for _, idStr := range ids {
			var id uint
			fmt.Sscan(idStr, &id)
			todoIDs = append(todoIDs, id)
		}
	} else {
		var todos []model.Todo
		if err := database.DB.Find(&todos).Error; err != nil {
			return nil, err
		}

		for _, t := range todos {
			todoIDs = append(todoIDs, t.ID)
			database.RDB.SAdd(database.Ctx, todoIDsKey(), t.ID)
		}
		database.RDB.Expire(database.Ctx, todoIDsKey(), time.Minute*10)
	}

	var result []model.Todo

	for _, id := range todoIDs {
		key := todoItemKey(id)

		val, err := database.RDB.Get(database.Ctx, key).Result()
		if err == nil {
			var todo model.Todo
			_ = json.Unmarshal([]byte(val), &todo)
			result = append(result, todo)
			continue
		}

		var todo model.Todo
		if err := database.DB.First(&todo, id).Error; err != nil {
			continue
		}

		data, _ := json.Marshal(todo)
		database.RDB.Set(database.Ctx, key, data, time.Minute*5)

		result = append(result, todo)
	}

	return result, nil
}

func CreateTodo(todo *model.Todo) error {
	if err := database.DB.Create(todo).Error; err != nil {
		return err
	}

	database.RDB.SAdd(database.Ctx, todoIDsKey(), todo.ID)

	data, _ := json.Marshal(todo)
	database.RDB.Set(
		database.Ctx,
		todoItemKey(todo.ID),
		data,
		time.Minute*5,
	)

	return nil
}
