package main

import (
	"log"
	"os"

	"github.com/Hi-Im-Toyo/GO_Proj/controllers"
	"github.com/Hi-Im-Toyo/GO_Proj/database"
	"github.com/Hi-Im-Toyo/GO_Proj/middleware"
	"github.com/Hi-Im-Toyo/GO_Proj/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantlybuy", app.InstantlyBuy())

	log.Fatal(router.Run(":" + port))
}
