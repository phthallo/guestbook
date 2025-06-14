package main

import (
	"github.com/gin-gonic/gin"
)

type Entry struct {
	ID         string    `json:"id"`
	Message    string    `json:"message"`
	Author     string	 `json:"author"`
}


func GetEntries(c *gin.Context){
	var entry Entry
}


func API(){
	router := gin.Default()
	router.GET("/entries", GetEntries)
	router.Run("localhost:4142")
}