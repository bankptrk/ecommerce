package config

import (
	"bank/models"
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddProductToCart(c *fiber.Ctx, productID primitive.ObjectID, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	searchFromDb, err := GetProductCollection().Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return c.SendString("Can't find product")
	}
	var productcart []models.ProductUser
	err = searchFromDb.All(ctx, &productcart)
	if err != nil {
		log.Println(err)
		return c.SendString("Can't find product")
	}
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return c.SendString("User is not valid")
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "user_cart", Value: bson.D{{Key: "$each", Value: productcart}}}}}}
	_, err = GetUserCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return c.SendString("Can't add product to cart")
	}
	return nil
}

func RemoveCartItem(c *fiber.Ctx, productID primitive.ObjectID, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return c.SendString("User is not valid")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"user_cart": bson.M{"_id": productID}}}
	_, err = GetUserCollection().UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return c.SendString("Can't remove item from cart")
	}
	return nil
}

func BuyItemFromCart(c *fiber.Ctx, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("User is not valid")
	}

	var getcart_item models.User
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	if err = GetUserCollection().FindOne(ctx, filter).Decode(&getcart_item); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve user cart items")
	}

	if len(getcart_item.UserCart) == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("No items in cart to purchase")
	}

	order_cart := models.Order{
		Order_ID:       primitive.NewObjectID(),
		Order_At:       time.Now(),
		Order_Cart:     make([]models.ProductUser, 0),
		Payment_Method: models.Payment{COD: true},
	}

	var total_price float64
	for _, product := range getcart_item.UserCart {
		order_cart.Order_Cart = append(order_cart.Order_Cart, product)
		total_price += product.Price
	}
	order_cart.Price = total_price

	update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "orders", Value: []models.Order{order_cart}}}}}
	if _, err = GetUserCollection().UpdateOne(ctx, filter, update); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update user orders")
	}

	usercart_empty := make([]models.ProductUser, 0)
	update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "user_cart", Value: usercart_empty}}}}
	if _, err = GetUserCollection().UpdateOne(ctx, filter, update3); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to empty user cart")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Purchase successful"})
}

func InstantBuyer(c *fiber.Ctx, productID primitive.ObjectID, UserID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("User is not valid")
	}

	order_detail := models.Order{
		Order_ID:       primitive.NewObjectID(),
		Order_At:       time.Now(),
		Order_Cart:     make([]models.ProductUser, 0),
		Payment_Method: models.Payment{COD: true},
	}

	var product_details models.ProductUser
	err = GetProductCollection().FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&product_details)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve product details")
	}

	order_detail.Price = product_details.Price

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: order_detail}}}}

	if _, err = GetUserCollection().UpdateOne(ctx, filter, update); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update user orders")
	}

	update2 := bson.M{"$push": bson.M{"orders.$.order_list": product_details}}
	if _, err = GetUserCollection().UpdateOne(ctx, filter, update2); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update order items")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Order placed successfully"})
}
