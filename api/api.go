package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/phthallo/guestbook/internal"
)

type Entry struct {
	Name       string	 `json:"name"`
	Message    string    `json:"message"`
}

func GetEntries(ctx *gin.Context, conn *pgx.Conn){
	var jsonBytes []byte
	err := conn.QueryRow(context.Background(), "SELECT json_agg(json_build_object('name', entries.name, 'message', entries.message)) from entries").Scan(&jsonBytes)
	var entries []Entry
	if err := json.Unmarshal(jsonBytes, &entries); err != nil {
		fmt.Printf("%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch entries"})
	}
	if (err != nil){
		fmt.Printf("%v", err)
	}
	ctx.IndentedJSON(http.StatusOK, entries)
}


func main(){
	conn, _ := internal.CreateDBConnection()

	router := gin.Default()
	router.GET("/entries", func (context *gin.Context){
		GetEntries(context, conn)
	})
	defer conn.Close(context.Background())

	router.Run("localhost:4142")

}