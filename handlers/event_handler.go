package handlers

import (
	"fiber-first-app/db"
	"fiber-first-app/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// CreateEvent creates a new event
func CreateEvent(c *fiber.Ctx) error {
	var e models.Event
	if err := c.BodyParser(&e); err != nil {
		logrus.Error("Failed to parse event: ", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()

	err := db.Conn.QueryRow(
		c.Context(),
		`INSERT INTO events (title, description, location, start_time, end_time, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		e.Title, e.Description, e.Location, e.StartTime, e.EndTime, e.CreatedAt, e.UpdatedAt,
	).Scan(&e.ID)

	if err != nil {
		logrus.Error("Failed to insert event: ", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(e)
}

// GetEvents returns all events
func GetEvents(c *fiber.Ctx) error {
	rows, err := db.Conn.Query(c.Context(), "SELECT id, title, description, location, start_time, end_time, created_at, updated_at FROM events")
	if err != nil {
		logrus.Error("Failed to fetch events: ", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.StartTime, &e.EndTime, &e.CreatedAt, &e.UpdatedAt); err != nil {
			logrus.Error("Scan error: ", err)
			return err
		}
		events = append(events, e)
	}
	return c.JSON(events)
}

// GetEventByID returns a single event
func GetEventByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var e models.Event
	err := db.Conn.QueryRow(
		c.Context(),
		"SELECT id, title, description, location, start_time, end_time, created_at, updated_at FROM events WHERE id=$1",
		id,
	).Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.StartTime, &e.EndTime, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		logrus.Error("Event not found: ", err)
		return c.Status(404).JSON(fiber.Map{"error": "Event not found"})
	}
	return c.JSON(e)
}

// UpdateEvent updates an event
func UpdateEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	var e models.Event
	if err := c.BodyParser(&e); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	e.UpdatedAt = time.Now()

	_, err := db.Conn.Exec(
		c.Context(),
		`UPDATE events SET title=$1, description=$2, location=$3, start_time=$4, end_time=$5, updated_at=$6 WHERE id=$7`,
		e.Title, e.Description, e.Location, e.StartTime, e.EndTime, e.UpdatedAt, id,
	)
	if err != nil {
		logrus.Error("Failed to update event: ", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(e)
}

// DeleteEvent deletes an event
func DeleteEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := db.Conn.Exec(c.Context(), "DELETE FROM events WHERE id=$1", id)
	if err != nil {
		logrus.Error("Failed to delete event: ", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
