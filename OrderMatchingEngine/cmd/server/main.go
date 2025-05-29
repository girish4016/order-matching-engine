package main

import (
	"OrderMatchingEngine/internal/api"
	"OrderMatchingEngine/internal/db"

	"github.com/gin-gonic/gin"
)

func main() {

	db := db.InitDB()
	defer db.Close()

	router := api.NewRouter(db)
	engine := gin.Default()
	router.RegisterRoutes(engine)
	engine.Run(":8080")

}
