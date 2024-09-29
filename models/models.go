package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	First_Name      string             `json:"first_name" validate:"required"`
	Last_Name       string             `json:"last_name" validate:"required"`
	Email           string             `json:"email" validate:"required"`
	Phone           string             `json:"phone" validate:"required"`
	Password        string             `json:"password" validate:"required"`
	Create_At       time.Time          `json:"create_at"`
	Update_At       time.Time          `json:"update_at"`
	UserCart        []ProductUser      `json:"user_cart" bson:"user_cart"`
	Address_Details []Address          `json:"address" bson:"address"`
	Order_Status    []Order            `json:"orders" bson:"orders"`
	Role            string             `json:"role" bson:"role"`
}

type Address struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	House   string             `json:"house_name" bson:"house_name"`
	Street  string             `json:"street_name" bson:"street_name"`
	City    string             `json:"city_name" bson:"city_name"`
	Zipcode string             `json:"zip_code" bson:"zip_code"`
}

type Product struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name string             `json:"product_name"`
	Price        float64            `json:"price"`
	Rating       uint               `json:"rating"`
	Image        string             `json:"image"`
}

type ProductUser struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name string             `json:"product_name" bson:"product_name"`
	Price        float64            `json:"price" bson:"price"`
	Rating       uint               `json:"rating" bson:"rating"`
	Image        string             `json:"image"`
}

type Order struct {
	Order_ID       primitive.ObjectID `bson:"_id"`
	Order_Cart     []ProductUser      `json:"order_list" bson:"order_list"`
	Order_At       time.Time          `json:"order_at" bson:"order_at"`
	Price          float64            `json:"total_price" bson:"total_price"`
	Discount       float64            `json:"discount" bson:"discount"`
	Payment_Method Payment            `json:"payment_method" bson:"payment_method"`
}

type Payment struct {
	Digital bool `json:"digital" bson:"digital"`
	COD     bool `json:"cod" bson:"cod"`
}
