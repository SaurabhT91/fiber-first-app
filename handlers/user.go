package handlers

import (
	"fiber-first-app/db"
	"fiber-first-app/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	rows, err := db.Conn.Query(c.Context(), "SELECT id, name, email FROM users")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return err
		}
		users = append(users, u)
	}
	return c.JSON(users)
}

func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	row := db.Conn.QueryRow(c.Context(), "SELECT id, name, email FROM users WHERE id=$1", id)

	var u models.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(u)
}

func CreateUser(c *fiber.Ctx) error {
	var u models.User
	if err := c.BodyParser(&u); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := db.Conn.QueryRow(
		c.Context(),
		"INSERT INTO users(name, email) VALUES($1, $2) RETURNING id",
		u.Name, u.Email,
	).Scan(&u.ID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(u)
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var u models.User
	if err := c.BodyParser(&u); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	i, _ := strconv.Atoi(id)
	u.ID = i

	_, err := db.Conn.Exec(
		c.Context(),
		"UPDATE users SET name=$1, email=$2 WHERE id=$3",
		u.Name, u.Email, u.ID,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(u)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := db.Conn.Exec(c.Context(), "DELETE FROM users WHERE id=$1", id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204) // No Content
}
