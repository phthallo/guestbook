package api

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Entry struct {
	Name       string	 `json:"name"`
	Message    string    `json:"message"`
    Timestamp  string    `json:"timestamp"`
}

func GetEntries(ctx *gin.Context, dbpool *pgxpool.Pool, limit string){
	// Returns entries in the guestbook
	// Limit, as the name suggests, limits the number of entries returned.
	var jsonBytes []byte
	if err := dbpool.QueryRow(context.Background(), "SELECT json_agg(json_build_object('name', limited_entries.name, 'message', limited_entries.message, 'timestamp', limited_entries.timestamp)) from (SELECT id, name, message, timestamp FROM entries ORDER BY id DESC LIMIT $1) AS limited_entries", limit).Scan(&jsonBytes); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf(`Failed to select entries %f`, err)})
		return
	}
	
	if len(jsonBytes) == 0 {
		ctx.IndentedJSON(http.StatusOK, make([]string, 0))
		return
	}

	var entries []Entry
	if err := json.Unmarshal(jsonBytes, &entries); err != nil {
		fmt.Printf("%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to jsonify entries"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, entries)
}


func StartAPIService(dbpool *pgxpool.Pool){
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.GET("/entries", func (context *gin.Context){
		limit := context.DefaultQuery("limit", "10")
		GetEntries(context, dbpool, limit)
	})
	router.Run(fmt.Sprintf(":%s",os.Getenv("API_PORT")))

}