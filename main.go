package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		IdleTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		Prefork:      true,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World !!!!!")
	})

	if fiber.IsChild() {
		fmt.Println("This is child process")
	} else {
		fmt.Println("This is parent process")
	}

	err := app.Listen("localhost:3000")

	if err != nil {
		panic(err)
	}
}
