package controllers

import (
	"context"
	"errors"
	"time"

	"github.com/AhishekOza/E-commerce/database"
	"github.com/AhishekOza/E-commerce/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getProductPrice(productId primitive.ObjectID) (int, error) {
	collection := database.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filters := bson.M{"_id": productId}
	var product models.Product
	if err := collection.FindOne(ctx, filters).Decode(&product); err != nil {
		return 0, errors.New(err.Error())
	}

	return product.Price, nil

}

func CreateOrder(c *fiber.Ctx) error {
	collection := database.DB.Collection("orders")
	ctx, causal := context.WithTimeout(context.Background(), 10*time.Second)
	defer causal()

	id := c.Locals("userId").(string)
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var order models.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err.Error()})
	}

	totalPrice := 0
	for index, item := range order.Cart {
		productPrice, err := getProductPrice(item.ProductId)
		// --Adding the product price
		order.Cart[index].Price = productPrice
		order.Cart[index].ID = primitive.NewObjectID()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err.Error()})
		}
		totalPrice += productPrice * item.Quantity
	}

	order.TotalPrice = totalPrice
	order.UserId = userId
	order.CreatedAt = time.Now()

	result, err := collection.InsertOne(ctx, order)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err.Error()})
	}

	order.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(fiber.StatusOK).JSON(order)
}

func GetAllOrders(c *fiber.Ctx) error {
	collection := database.DB.Collection("orders")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var orders []models.Order

	filters := bson.M{}
	sort := bson.D{{Key: "created_at", Value: -1}}

	cursor, err := collection.Find(ctx, filters, options.Find().SetSort(sort))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"erorr": err.Error()})
	}

	for cursor.Next(ctx) {
		var order models.Order

		if err := cursor.Decode(&order); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"erorr": err.Error()})
		}

		orders = append(orders, order)
	}

	return c.JSON(orders)

}

func GetOrderById(c *fiber.Ctx) error {
	collection := database.DB.Collection("orders")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	ObjectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	filters := bson.M{"_id": ObjectId}
	var order models.Order
	if err := collection.FindOne(ctx, filters).Decode(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": err.Error()})
	}

	return c.JSON(order)
}

func DeleteOrder(c *fiber.Ctx) error {
	collection := database.DB.Collection("orders")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	ObjectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	filters := bson.M{"_id": ObjectId}

	if _, err := collection.DeleteOne(ctx, filters); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Message": "Order Deleted Successfully"})
}

func UpdateOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var order models.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	collection := database.DB.Collection("orders")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	filters := bson.M{"_id": objID}

	update := bson.M{
		"$set": order,
	}
	_, err = collection.UpdateOne(ctx, filters, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Order updated successfully"})
}
