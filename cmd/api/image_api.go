package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-pinterest/db"
	"go-pinterest/model"
	"net/http"
	"strconv"
	"time"
)
var (
	httpPort = flag.String("http_port", "9090", "http_port listen")
)
func main() {
	flag.Parse()
	db.Init()
	defer db.Close()
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/categories", GetAllCategorys)
	r.GET("/categories/images", GetImagesByCategory)
	r.GET("/home/rands", GetHomeRandomImage)
	r.Run(":"+*httpPort)
}

func GetAllCategorys(c *gin.Context) {
	categorys, err := db.GetAllCategory()
	if err != nil {
		fmt.Println(c.Request.RequestURI)
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SERVER_ERROR, Message: "Error when read data from server"})
		return
	}
	c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SUCCESS, Data: categorys})
}

func GetImagesByCategory(c *gin.Context) {
	category, exit := c.GetQuery("category")
	if !exit {
		fmt.Println(c.Request.RequestURI)
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_CLIENT_ERROR, Message: "Error query param category , must set query ?category=&offset=&length= "})
		return
	}
	offsetText, exit := c.GetQuery("offset")
	defaultOffset:= time.Now().UnixNano()
	var offset int64 = 0
	var err error
	if exit {
		offset, err = strconv.ParseInt(offsetText, 10, 64)
		if err != nil {
			fmt.Println(c.Request.RequestURI)
			c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_CLIENT_ERROR, Message: "Error when offset query"})
			return
		}
	}
	if offset==0 {
		offset=defaultOffset
	}
	lengthText, exit := c.GetQuery("length")
	var length int64 = 20
	if exit {
		length, err = strconv.ParseInt(lengthText, 10, 64)
		if err != nil {
			fmt.Println(c.Request.RequestURI)
			c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_CLIENT_ERROR, Message: "Error when length query"})
			return
		}
	}
	if length > 20{
		length=20
	}
	data, err := db.GetImageByCategory(category, offset, length)
	if err != nil {
		fmt.Println(c.Request.RequestURI)
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SERVER_ERROR, Message: "Error when read data from server"})
		return
	}
	c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SUCCESS, Data: data})
}


func GetHomeRandomImage(c *gin.Context) {
	offsetText, exit := c.GetQuery("offset")
	defaultOffset:= time.Now().UnixNano()
	var offset int64 = 0
	var err error
	if exit {
		offset, err = strconv.ParseInt(offsetText, 10, 64)
		if err != nil {
			fmt.Println(c.Request.RequestURI)
			c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_CLIENT_ERROR, Message: "Error when offset query"})
			return
		}
	}
	if offset==0 {
		offset=defaultOffset
	}
	lengthText, exit := c.GetQuery("length")
	var length int64 = 20
	if exit {
		length, err = strconv.ParseInt(lengthText, 10, 64)
		if err != nil {
			fmt.Println(c.Request.RequestURI)
			c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_CLIENT_ERROR, Message: "Error when length query"})
			return
		}
	}
	if length > 20{
		length=20
	}
	data, err := db.GetHomeImages(offset, length)
	if err != nil {
		fmt.Println(c.Request.RequestURI)
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SERVER_ERROR, Message: "Error when read data from server"})
		return
	}
	c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SUCCESS, Data: data})
}