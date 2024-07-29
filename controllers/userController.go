package controllers

import (
	"context"
	"time"

	"github.com/AhishekOza/E-commerce/database"
	"github.com/AhishekOza/E-commerce/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {

	// get the user info
	// parse the info
	// check if the user is already registered
	// if yes then send an error message

	// make a password bcrypt
	// insert the user
	// send the json response
	collection := database.DB.Collection("users")
	user := &models.User{}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	filters := bson.M{"email": user.Email}

	count, err := collection.CountDocuments(context.Background(), filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if count > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email already exists"})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(password)

	result, err := collection.InsertOne(context.Background(), user)

	if err != nil {
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(200).JSON(user)
}

func Login(c *fiber.Ctx) error {
	// get the user info
	// decode the json data in struct
	// if not able to decode throw an error
	// now find the user
	// check the password and then  create a jwt token
	// send then token with user

	collection := database.DB.Collection("users")
	user := &models.User{}

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	filters := bson.M{"email": user.Email}

	dbUser := models.User{}

	err := collection.FindOne(context.Background(), filters).Decode(&dbUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Password"})
	}

	claims := jwt.MapClaims{
		"userId": dbUser.ID.Hex(),
		"admin":  dbUser.Admin,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(database.JWTSecret)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"token": tokenString, "user": dbUser})
}
