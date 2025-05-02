package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"
	"path/filepath"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024,
	})

	uploadFolder := "./upload"
	err := os.MkdirAll(uploadFolder, os.ModePerm)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/upload", func(c *fiber.Ctx) error {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "No file uploaded",
			})
		}

		fileName := filepath.Base(fileHeader.Filename)
		filePath := filepath.Join(uploadFolder, fileName)

		// Remove the file opening code here - just save the file directly
		if err := c.SaveFile(fileHeader, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success":  true,
			"filename": fileName,
		})
	})

	app.Get("/download/:filename", func(c *fiber.Ctx) error {
		fileName := c.Params("filename")
		filePath := filepath.Join(uploadFolder, fileName)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "File not found",
			})
		}

		return c.SendFile(filePath, false)
	})

	app.Delete("/:filename", func(c *fiber.Ctx) error {
		fileName := c.Params("filename")
		filePath := filepath.Join(uploadFolder, fileName)
		if err := os.Remove(filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
		})
	})

	log.Info("Server starting on :8084")
	err = app.Listen(":8084")
	if err != nil {
		log.Fatal(err.Error())
	}
}
