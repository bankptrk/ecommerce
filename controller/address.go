package controller

import (
	"bank/config"
	"bank/models"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress(c *fiber.Ctx) error {
	user_id := c.Query("id")
	if user_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid Search Index"})
	}
	address_id, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	var address models.Address
	address.ID = primitive.NewObjectID()
	if err := c.BodyParser(&address); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid Input"})
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address_id}}}}
	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
	group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

	cursor, err := config.GetUserCollection().Aggregate(ctx, mongo.Pipeline{filter, unwind, group})
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var address_details []bson.M
	if err = cursor.All(ctx, &address_details); err != nil {
		panic(err)
	}
	var size int32
	for _, address_no := range address_details {
		count := address_no["count"]
		size = count.(int32)
	}
	if size < 2 {
		filter := bson.D{primitive.E{Key: "_id", Value: address_id}}
		update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: address}}}}
		_, err := config.GetUserCollection().UpdateOne(ctx, filter, update)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Not Allowed"})
	}
	ctx.Done()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Add address successfully!"})
}

func EditBillingAddress(c *fiber.Ctx) error {
	user_id := c.Query("id")
	if user_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid Search Index"})
	}
	user_no, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	var edit_address models.Address
	if err := c.BodyParser(&edit_address); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid input"})
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{primitive.E{Key: "_id", Value: user_no}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: edit_address.House},
		{Key: "address.0.street_name", Value: edit_address.Street},
		{Key: "address.0.city_name", Value: edit_address.City},
		{Key: "address.0.zip_code", Value: edit_address.Zipcode}}}}
	_, err = config.GetUserCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": "Something went wrong"})
	}
	ctx.Done()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Update address successful"})
}

func EditShippingAddress(c *fiber.Ctx) error {
	user_id := c.Query("id")
	if user_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid Search Index"})
	}
	user_no, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	var edit_address models.Address
	if err := c.BodyParser(&edit_address); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid input"})
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{primitive.E{Key: "_id", Value: user_no}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: edit_address.House},
		{Key: "address.1.street_name", Value: edit_address.Street},
		{Key: "address.1.city_name", Value: edit_address.City},
		{Key: "address.1.zip_code", Value: edit_address.Zipcode}}}}
	_, err = config.GetUserCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": "Something went wrong"})
	}
	ctx.Done()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Update address successful"})
}

func DeleteAddress(c *fiber.Ctx) error {
	user_id := c.Query("id")
	if user_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid Search Index"})
	}
	address := make([]models.Address, 0)
	user_no, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	filter := bson.D{primitive.E{Key: "_id", Value: user_no}}
	update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: address}}}}
	_, err = config.GetUserCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Error": "Not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Delete successful"})
}
