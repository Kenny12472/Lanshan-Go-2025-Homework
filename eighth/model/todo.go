package model

type Todo struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	WhatTodo string `json:"what_todo"`
	TimeTodo string `json:"time_todo"`
}
