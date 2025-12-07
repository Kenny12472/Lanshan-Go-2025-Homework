package protected

import (
	"errors"

	"seventh2/utils"
)

var taskList = []string{}

func AddTask(token string, task string) error {

	claims, err := utils.ParseToken(token)
	if err != nil {
		return errors.New("无效 token，请先登录")
	}

	username := claims.Username

	taskList = append(taskList, username+" 的任务："+task)

	return nil
}

func GetTasks(token string) ([]string, error) {

	_, err := utils.ParseToken(token)
	if err != nil {
		return nil, errors.New("请先登录")
	}

	return taskList, nil
}
