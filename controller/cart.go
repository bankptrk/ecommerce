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
	"go.mongodb.org/mongo-driver/mongo"
)

func AddToCart(c *fiber.Ctx) error {
	productQueryId := c.Query("pid")

	if productQueryId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid Search Index"})
	}

	userQueryId := c.Query("uid")
	if userQueryId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid Search Index"})
	}
	product_id, err := primitive.ObjectIDFromHex(productQueryId)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	err = config.AddProductToCart(c, product_id, userQueryId)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Add to cart suscessful"})
}

func RemoveItemFromCart(c *fiber.Ctx) error {
	productQueryId := c.Query("pid")
	if productQueryId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Product id is empty"})
	}
	userQueryId := c.Query("uid")
	if userQueryId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "User id is empty"})
	}
	ProductID, err := primitive.ObjectIDFromHex(productQueryId)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	err = config.RemoveCartItem(c, ProductID, userQueryId)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Remove item from cart successful"})
}

func GetItemsFromCart(c *fiber.Ctx) error {
	user_id := c.Query("id")
	if user_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "User id is empty"})
	}
	user_no, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var fill_cart models.User
	err = config.GetUserCollection().FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: user_no}}).Decode(&fill_cart)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: user_no}}}}
	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$user_cart"}}}}
	group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$user_cart.price"}}}}}}

	cursor, err := config.GetUserCollection().Aggregate(ctx, mongo.Pipeline{filter, unwind, group})
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var list []bson.M
	if err = cursor.All(ctx, &list); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var totalPrice float64
	if len(list) > 0 {
		totalPrice = list[0]["total"].(float64)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total":    totalPrice,
		"userCart": fill_cart.UserCart,
	})
}

func BuyFromCart(c *fiber.Ctx) error {
	userQueryID := c.Query("id")
	if userQueryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "User id is empty"})
	}

	err := config.BuyItemFromCart(c, userQueryID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Place order successful"})
}

func InstantBuy(c *fiber.Ctx) error {

	userQueryID := c.Query("uid")
	if userQueryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "User id is empty"})
	}
	productQueryID := c.Query("pid")
	if productQueryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Product id is empty"})
	}
	productID, err := primitive.ObjectIDFromHex(productQueryID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	err = config.InstantBuyer(c, productID, userQueryID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Purchased goods successfully"})
}
