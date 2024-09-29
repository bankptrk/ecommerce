package controller

import (
	"bank/config"
	"bank/models"
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateProduct(c *fiber.Ctx) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var products models.Product
	defer cancel()
	if err := c.BodyParser(&products); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	if products.Product_Name == "" || products.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Product name and price are required"})
	}
	products.Product_ID = primitive.NewObjectID()
	_, err := config.GetProductCollection().InsertOne(ctx, products)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Can't create product"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Add product successfully!"})
}

func GetAllProducts(c *fiber.Ctx) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var productslist []models.Product
	cursor, err := config.GetProductCollection().Find(ctx, bson.D{{}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong"})
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &productslist); err != nil {
		log.Println("Error retrieving products:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve products"})
	}
	if err := cursor.Err(); err != nil {
		log.Println("Cursor error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid cursor"})
	}
	return c.Status(fiber.StatusOK).JSON(productslist)
}

func SearchProductByQuery(c *fiber.Ctx) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var searchproducts []models.Product
	queryParam := c.Query("name")
	if queryParam == "" {
		log.Println("query is empty")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid Search Index"})
	}

	searchquerydb, err := config.GetProductCollection().Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})
	if err != nil {
		log.Println("Error fetching from database:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": "Something went wrong in fetching the db query"})
	}
	defer searchquerydb.Close(ctx)

	if err := searchquerydb.All(ctx, &searchproducts); err != nil {
		log.Println("Error decoding cursor results:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid request"})
	}

	if err := searchquerydb.Err(); err != nil {
		log.Println("Cursor error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid request"})
	}

	return c.Status(fiber.StatusOK).JSON(searchproducts)
}
