package handlers

import (
	"fiber-first-app/db"
	"fiber-first-app/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// GetUsers - Get all users
func GetUsers(c *fiber.Ctx) error {
	logrus.Info("📥 Retrieving all users...")

	rows, err := db.Conn.Query(c.Context(), "SELECT id, name, email, password FROM users")
	if err != nil {
		logrus.Error("❌ Error querying database for users: ", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password); err != nil {
			logrus.Error("❌ Error scanning user row: ", err)
			return err
		}
		users = append(users, u)
	}
	logrus.Infof("✅ Successfully retrieved %d users.", len(users))
	return c.JSON(users)
}

// GetUserByID - Get a user by their ID
func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	logrus.Infof("📥 Retrieving user by ID: %s", id)

	row := db.Conn.QueryRow(c.Context(), "SELECT id, name, email, password FROM users WHERE id=$1", id)

	var u models.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password); err != nil {
		logrus.Error("❌ User not found: ", err)
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	logrus.Infof("✅ Successfully retrieved user: %s", u.Name)
	return c.JSON(u)
}

// CreateUser - Create a new user
func CreateUser(c *fiber.Ctx) error {
	logrus.Info("📥 Creating a new user...")

	var u models.User

	// Parse the incoming request
	if err := c.BodyParser(&u); err != nil {
		logrus.Error("❌ Failed to parse request body: ", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	logrus.Infof("📥 Parsed user input: Name=%s, Email=%s, Password=%s", u.Name, u.Email, u.Password)

	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error("❌ Failed to hash password: ", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	logrus.Infof("🔐 Final hashed password to store: %s", string(hashedPassword))

	// Set timestamps
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	logrus.Infof("🕒 Timestamps set: CreatedAt=%s, UpdatedAt=%s", u.CreatedAt.Format(time.RFC3339), u.UpdatedAt.Format(time.RFC3339))

	// Insert user into the database
	logrus.Info("📤 Inserting user into the database...")
	// err = db.Conn.QueryRow(
	// 	c.Context(),
	// 	`INSERT INTO users(name, email, password, created_at, updated_at)
	// 	 VALUES($1, $2, $3, $4, $5) RETURNING id`,
	// 	u.Name, u.Email, string(hashedPassword), u.CreatedAt, u.UpdatedAt,
	// ).Scan(&u.ID)

	err = db.Conn.QueryRow(
		c.Context(),
		`INSERT INTO users(name, email, password, created_at, updated_at) 
     VALUES($1, $2, $3, $4, $5) RETURNING id`,
		u.Name, u.Email, string(hashedPassword), u.CreatedAt, u.UpdatedAt,
	).Scan(&u.ID)

	if err != nil {
		logrus.Error("❌ Failed to insert user into DB: ", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	logrus.Infof("✅ User inserted with ID: %d", u.ID)

	// Do not return the password in the response
	u.Password = "" // Remove password from response

	return c.Status(201).JSON(u)
}

// UpdateUser - Update user details
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	logrus.Infof("📥 Updating user by ID: %s", id)

	var u models.User
	if err := c.BodyParser(&u); err != nil {
		logrus.Error("❌ Failed to parse request body: ", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	i, _ := strconv.Atoi(id)
	u.ID = i

	// Hash the password (if updating)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error("❌ Failed to hash password: ", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	u.Password = string(hashedPassword)
	logrus.Infof("🔐 Password hashed successfully for update.")

	_, err = db.Conn.Exec(
		c.Context(),
		"UPDATE users SET name=$1, email=$2, password=$3 WHERE id=$4",
		u.Name, u.Email, u.Password, u.ID,
	)

	if err != nil {
		logrus.Error("❌ Failed to update user: ", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	logrus.Infof("✅ Successfully updated user ID: %d", u.ID)
	return c.JSON(u)
}

// DeleteUser - Delete a user by their ID
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	logrus.Infof("📥 Deleting user by ID: %s", id)

	_, err := db.Conn.Exec(c.Context(), "DELETE FROM users WHERE id=$1", id)

	if err != nil {
		logrus.Error("❌ Failed to delete user: ", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	logrus.Infof("✅ Successfully deleted user ID: %s", id)
	return c.SendStatus(204) // No Content
}
