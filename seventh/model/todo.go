package model

type Todo struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	WhatTodo string `gorm:"comment:'事项内容'" json:"what_todo"`
	TimeTodo string `gorm:"comment:'事项时间'" json:"time_todo"`
}
