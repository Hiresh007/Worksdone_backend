package controllers

import (
	"context"
	"log"
	"net/http"
	"server/database"
	"server/helpers"
	"server/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var validate = validator.New()

func HasPassword(password string) string {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		log.Panic(err)
	}

	return string(pass)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	msg := "Password Correct"
	check := true
	if err != nil {
		msg = "Password Incorrect"
		check = false
		log.Panic(err)
	}
	return check, msg

}
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		countE, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while checking email"})
		}
		countU, err := userCollection.CountDocuments(ctx, bson.M{"username": user.Username})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while checking email"})
		}
		password := HasPassword(*user.Password)
		user.Password = &password
		if countE > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email Alreadt Exists"})
		}
		if countU > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Username already exists"})
		}
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refreshToken, _ := helpers.GenerateTokens(*user.Email, *user.Username, user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken
		res, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			msg := "User item not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		defer cancel()
		c.JSON(http.StatusOK, res)

	}
}
func Signin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Username or Password Incorrect"})
			return
		}

		passwordValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !passwordValid {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": msg})
			return
		}
		if foundUser.Username == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "User not found"})
			return
		}

		token, refreshToken, _ := helpers.GenerateTokens(*foundUser.Email, *foundUser.Username, foundUser.User_id)

		helpers.UpdateTokens(token, refreshToken, foundUser.User_id)
		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}
}
