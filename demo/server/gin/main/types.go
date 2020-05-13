package main

import "github.com/jinzhu/gorm"

type (
	// 定义原始的数据库字段
	todoModel struct {
		gorm.Model
		Title     string `json:"title"`
		Completed int    `json:"completed"`
	}
	// 处理返回的字段
	transformedTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
)
