package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AhishekOza/E-commerce/database"
	"github.com/AhishekOza/E-commerce/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCategory(c *fiber.Ctx) error {
	collection := database.DB.Collection("categories")
	category := &models.Category{}

	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	filters := bson.M{"name": category.Name}

	count, err := collection.CountDocuments(context.Background(), filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if count > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Category already exists"})
	}

	result, err := collection.InsertOne(context.Background(), category)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	category.ID = result.InsertedID.(primitive.ObjectID)
	return c.Status(200).JSON(category)
}

func GetAllCategories(c *fiber.Ctx) error {
	collection := database.DB.Collection("categories")
	categories := []models.Category{}

	filters := bson.M{}

	cursor, err := collection.Find(context.Background(), filters)

	if err != nil {
		log.Fatal(err)
		return err
	}

	// close the connection with the cursor as it frees the resources
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		category := models.Category{}

		err := cursor.Decode(&category)
		if err != nil {
			log.Fatal(err)
			return err
		}

		categories = append(categories, category)
	}

	return c.Status(200).JSON(categories)
}

func GetSingleCategory(c *fiber.Ctx) error {
	collection := database.DB.Collection("categories")

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id is empty"})
	}

	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	filters := bson.M{"_id": ObjectID}

	var category models.Category

	err = collection.FindOne(context.Background(), filters).Decode(&category)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(category)
}

func UpdateCategory(c *fiber.Ctx) error {
	collection := database.DB.Collection("categories")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var category models.Category

	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	id := c.Params("id")

	ObjectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Object ID"})
	}

	filters := bson.M{"_id": ObjectId}
	options := bson.M{"$set": category}

	_, err = collection.UpdateOne(ctx, filters, options)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Message": "Category updated successfully !"})

}

func DeleteCategory(c *fiber.Ctx) error {
	collection := database.DB.Collection("categories")

	id := c.Params("id")
	categoryId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	filters := bson.M{"_id": categoryId}

	result, err := collection.DeleteOne(context.Background(), filters)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Println(result)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Category deleted successfully!"})
}
