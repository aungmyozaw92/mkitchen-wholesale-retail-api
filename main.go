package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/cmd"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/models"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/routes"
)

func main(){
	models.ConnectDatabase()

	cmd.Execute()

	r := gin.Default()

	// Router
	routes.SetupRoutes(r)

	r.NoRoute(customNotFoundHandler)

    r.Run(":8000")
}


func customNotFoundHandler(c *gin.Context) {
    c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
}