package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	UserName     *string            `json:"username" bson:"username" validate:"required,min=2,max=30"`
	Password     *string            `json:"password" bson:"password" validate:"required"`
	Token        *string            `json:"access_token" bson:"access_token"`
	RefreshToken *string            `json:"refreshToken" bson:"refreshToken"`
	CreatedTime  time.Time          `json:"createdTime" bson:"createdTime"`
	UpdatedTime  time.Time          `json:"updatedTime" bson:"updatedTime"`
	UserId       string             `json:"userId" bson:"userId"`
	UserCart     []ProductUser      `json:"usercart" bson:"usercart"`
}

type Product struct {
	ProductId   primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Type        string             `json:"type" bson:"type"`
	SKU         string             `json:"sku" bson:"sku"`
	ImageURL    string             `json:"image_url" bson:"image_url"`
	Description string             `json:"description" bson:"description"`
	Quantity    int                `json:"quantity" bson:"quantity"`
	Price       float64            `json:"price" bson:"price"`
}

type ProductUser struct {
	ProductId   primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Type        string             `json:"type" bson:"type"`
	SKU         string             `json:"sku" bson:"sku"`
	ImageURL    string             `json:"image_url" bson:"image_url"`
	Description string             `json:"description" bson:"description"`
	Quantity    int                `json:"quantity" bson:"quantity"`
	Price       float64            `json:"price" bson:"price"`
}
