package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type entry struct {
	ID         string    `json:"id"`
	Message    string    `json:"message"`
	Author     string	 `json:"author"`
}

func ConnectToDB(){
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
}

func GetEntries(c *gin.Context){


}

func PostEntries(c *gin.Context){
	var newEntry entry
	if err := c.BindJSON(&newEntry); err != nil {
		return
	}
}

func API(){
	router := gin.Default()
	router.GET("/entries", GetEntries)
	router.POST("/entries", PostEntries)
	router.Run("localhost:4142")
}