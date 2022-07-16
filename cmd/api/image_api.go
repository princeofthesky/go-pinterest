package main

import (
	"github.com/gin-gonic/gin"
	"go-pinterest/db"
	"net/http"
)

var (
	folderNFT = "/data/nft/images"
)

func main() {
	db.Init()
	defer db.Close()
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

}

func nftByUsers(c *gin.Context) {
}

