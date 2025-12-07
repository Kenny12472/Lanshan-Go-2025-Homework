package main

import (
	"fmt"

	"seventh2/dao"
	"seventh2/protected"
)

func main() {

	err := dao.InitDB()
	if err != nil {
		panic(err)
	}

	_ = dao.RegisterUser("zixuan", "123456")

	token, err := dao.LoginUser("zixuan", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("JWT：", token)

	err = protected.AddTask(token, "完成Go作业")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("任务添加成功")
	}

	tasks, _ := protected.GetTasks(token)
	fmt.Println("任务列表：", tasks)
}
