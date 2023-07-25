package db

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/sqlite3"
)

var (
	// db 对象
	storage fiber.Storage
)

func init() {
	storage = sqlite3.New(
		sqlite3.Config{
			Database: "./storage.sqlite3",
		},
	)
}

func Db() fiber.Storage {
	return storage
}
