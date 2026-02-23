package main

import (
	"project/internal/db"
	"project/internal/router"
)

func main() {
	db.InitDB()
	r := router.SetupRouter()
	r.Run(":8080")
}
