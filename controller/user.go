package controller

import (
	"bank/config"
	"bank/models"
	"context"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var Validate = validator.New()

func RegisterUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error invalid input"})
	}
	return createUser(c, user)
}

func LoginUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error invalid input"})
	}
	token, err := AuthUser(c, user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
	})
	return c.JSON(fiber.Map{"message": "Login successful"})
}

func createUser(c *fiber.Ctx, user *models.User) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	validation := Validate.Struct(user)
	if validation != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	if user.Email == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password are required"})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Can't hash password"})
	}
	user.Password = string(hashedPassword)

	var existingEmail models.User
	err = config.GetUserCollection().FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingEmail)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email already exists"})
	}

	user.ID = primitive.NewObjectID()
	user.Create_At = time.Now().Unix()
	user.Update_At = time.Now().Unix()
	user.UserCart = make([]models.ProductUser, 0)
	user.Address_Details = make([]models.Address, 0)
	user.Order_Status = make([]models.Order, 0)
	user.Role = "user"

	_, err = config.GetUserCollection().InsertOne(ctx, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Can't create user"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully!"})
}

func generateToken(email string) (string, error) {
	jwtSecret := []byte(os.Getenv("jwtSecret"))
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func AuthUser(c *fiber.Ctx, user *models.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	selectedUser := new(models.User)
	err := config.GetUserCollection().FindOne(ctx, bson.M{"email": user.Email}).Decode(&selectedUser)
	if err != nil {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(selectedUser.Password), []byte(user.Password))
	if err != nil {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}
	token, err := generateToken(selectedUser.Email)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Could not create token")
	}
	return token, nil
}

func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	cursor, err := config.GetUserCollection().Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return err
		}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		return err
	}
	return c.JSON(users)
}
