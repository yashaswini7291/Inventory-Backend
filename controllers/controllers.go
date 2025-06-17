package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/yashaswini7291/Inventory/database"
	"github.com/yashaswini7291/Inventory/models"
	"github.com/yashaswini7291/Inventory/tokens"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var (
	UserCollection    *mongo.Collection = database.UserData(database.Client, "Users")
	ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
	Validate                            = validator.New()
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login or Password is incorrect"
		valid = false
	}
	return valid, msg
}

// SignUp godoc
// @Summary Register a new user
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param user body models.User true "User Info"
// @Success 201 {string} string "account created successfully"
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /register [post]
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request payload"})
			return
		}
		// Validate the user input
		if validationErr := Validate.Struct(user); validationErr != nil {
			c.JSON(400, gin.H{"error": validationErr.Error()})
			return
		}
		count, err := UserCollection.CountDocuments(ctx, bson.M{"username": user.UserName})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		log.Println("Count of user with same username", count)
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exist"})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password
		user.CreatedTime, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedTime, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()
		token, refreshtoken, _ := tokens.TokenGenerator(*user.UserName)
		user.Token = &token
		user.RefreshToken = &refreshtoken
		user.UserCart = make([]models.ProductUser, 0)
		_, err = UserCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user creation failed"})
		}

		defer cancel()

		c.JSON(http.StatusCreated, "account created successfully")
	}
}

// Login godoc
// @Summary Login and get JWT token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param user body models.User true "Credentials"
// @Success 200 {object} map[string]string
// @Failure 400,500 {object} map[string]string
// @Router /login [post]
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var founduser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		err := UserCollection.FindOne(ctx, bson.M{"username": user.UserName}).Decode(&founduser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No User Found"})
			return
		}
		passwordIsValid, msg := VerifyPassword(*user.Password, *founduser.Password)

		defer cancel()

		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}

		token, refreshToken, _ := tokens.TokenGenerator(*founduser.UserName)
		defer cancel()

		tokens.UpdateAllTokens(token, refreshToken, founduser.UserId)

		c.JSON(http.StatusOK, gin.H{"access_token": token})
	}
}

// UpdateProductQuantity godoc
// @Summary Update the quantity of a product
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param quantity body map[string]int true "Quantity to update"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string
// @Router /products/{id}/quantity [put]
func UpdateProductQuantity() gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("id")

		objID, err := primitive.ObjectIDFromHex(productID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		var req struct {
			Quantity int `json:"quantity"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.Quantity < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		filter := bson.M{"_id": objID}
		update := bson.M{"$set": bson.M{"quantity": req.Quantity}}

		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		var updatedProduct models.Product
		err = ProductCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedProduct)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found or update failed"})
			return
		}

		c.JSON(http.StatusOK, updatedProduct)
	}
}

// GetAllProducts godoc
// @Summary Get a list of all products
// @Tags Products
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Product
// @Failure 500 {object} map[string]string
// @Router /products [get]
func GetAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "something went wrong please try after sometime")
			return
		}

		err = cursor.All(ctx, &productList)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}

		defer cancel()
		c.IndentedJSON(200, productList)
	}
}

// AddProduct godoc
// @Summary Add a new product
// @Tags Products
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param product body models.Product true "Product to Add"
// @Success 201 {object} map[string]interface{}
// @Failure 400,500 {object} map[string]string
// @Router /products [post]
func AddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var products models.Product
		defer cancel()
		if err := c.BindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		products.ProductId = primitive.NewObjectID()
		_, err := ProductCollection.InsertOne(ctx, products)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated, gin.H{
			"message":    "Product added successfully",
			"product_id": products.ProductId.Hex(),
		})
	}
}
