package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Hi-Im-Toyo/GO_Proj/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.UserData(database.Client, "user")
var productCollection *mongo.Collection = database.ProductData(database.Client, "product")	
var Validate = validator.New()


func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	
	}
	return string(bytes)

}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "invalid password"
		valid = false
	}
	return valid, msg
}

func Signup() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr()})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user email already exists"})
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user phone already exists"})

		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))	
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()

		token, refreshToken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshToken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)
		_, inserterr := UserCollection.InsertOne(ctx, user)

		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": inserterr})
			return
		}
		defer cancel()

		c.JSON(http.statuscreated, "Signup successful")
	}
}

func Login() gin.HandlerFunc {
	return function(c *gin.Context) {
		context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()
		var user models.User 
		c.BindJSON(&user)	

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)	
		defer cancel ()

		if founduser.Email == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
			return
		}

		PasswordisValid, msg := VerifyPassword(*user.Password, *founduser.Password)
		if !PasswordisValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}
		token, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, founduser.User_ID)
		defer cancel()

		generate.UpdateAllTokens(token, refreshToken, founduser.User_ID)

		c.JSON(http.StatusFound, founduser)
		
	}
}

func ProductViewerAdmin() gin.HandlerFunc {

}

func searchProduct() gin.HandlerFunc {

	return func(c *gin.Context) {

		var productlist []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)	
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "something went wrong please try again later")
			return
		} 

		err = cursor.All(ctx, &productlist)
		if err != nil {	
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError)
			return 
		}
		defer cursor.Close()

		if err := cursor.err(); err != nil {
			log.Println(err)
			c.AbortWithError(400, "invalid")
			return
		}
		defer cancel()
		c.IndentedJSON(200, productlist)

	}

}

func searchProductByQuery() gin.HandlerFunc {

}
