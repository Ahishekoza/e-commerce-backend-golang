package controllers

import (
	"context"
	"strconv"
	"time"

	"fmt"
	"log"
	"mime/multipart"

	"github.com/AhishekOza/E-commerce/database"
	"github.com/AhishekOza/E-commerce/models"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cld *cloudinary.Cloudinary

func init() {
	var err error
	cld, err = cloudinary.NewFromParams("dwkwy8qis", "136254854955345", "c5gD3n0cckRqFCQXpA-xlNKmdvo")

	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}

}

func UploadToCloudinary(file *multipart.FileHeader) (string, error) {

	src, err := file.Open()
	if err != nil {
		return "", err
	}

	response, err := cld.Upload.Upload(context.Background(), src, uploader.UploadParams{})
	if err != nil {
		return "", err
	}

	return response.SecureURL, nil
}

func CreateProduct(c *fiber.Ctx) error {
	collection := database.DB.Collection("products")

	product := &models.Product{}

	// Parse form data into product struct
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Print(product)
	// Parse CategoryID from the form data
	categoryIDStr := c.FormValue("category_id")
	if categoryIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "CategoryID is required"})
	}

	categoryID, err := primitive.ObjectIDFromHex(categoryIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid CategoryID"})
	}

	// Set CategoryID in the product struct
	product.CategoryId = categoryID

	// Handle image upload
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to upload image"})
	}

	imageUrl, err := UploadToCloudinary(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload image to Cloudinary"})
	}

	// Set the image URL in the product struct
	product.Image = imageUrl

	// Insert product into the database
	result, err := collection.InsertOne(context.Background(), product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	product.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(fiber.StatusOK).JSON(product)
}

func GetProductByCategory(c *fiber.Ctx) error {
	categoryID := c.Query("category_id")
	if categoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Category ID is required"})
	}

	fmt.Println(categoryID)

	collection := database.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert the categoryID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Category ID"})
	}

	filter := bson.M{"category_id": objectID}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(products)

}

func GetAllProducts(c *fiber.Ctx) error {
	collection := database.DB.Collection("products")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ---default page
	pageStr := c.Query("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": "Invalid page number"})
	}
	// ---limit of products to be displayed
	// limitStr := c.Query("limit", "10")
	// limit, err := strconv.Atoi(limitStr)
	// if err != nil {
	// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": "Invalid limit number"})
	// }

	limit := 10

	skip := (page - 1) * limit

	var products []models.Product

	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	for cursor.Next(ctx) {
		var product models.Product

		if err := cursor.Decode(&product); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		products = append(products, product)

	}

	// Get total count for pagination
	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":       products,
		"page":       page,
		"totalCount": totalCount,
		"totalPages": (totalCount + int64(limit) - 1) / int64(limit),
	})

}

func GetSingleProduct(c *fiber.Ctx) error {
	collection := database.DB.Collection("products")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	ObjectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	filters := bson.M{"_id": ObjectID}
	var product models.Product
	err = collection.FindOne(ctx, filters).Decode(&product)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	collection := database.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	var upadateProduct models.Product
	if err := c.BodyParser(&upadateProduct); err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	id := c.Params("id")
	ObjectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{"error": err.Error()})
	}

	filters := bson.M{"_id": ObjectID}
	update := bson.M{"$set": upadateProduct}

	result, err := collection.UpdateOne(ctx, filters, update)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Println(result)

	return c.JSON(fiber.Map{"success": "Product Updated Successfully !"})
}
