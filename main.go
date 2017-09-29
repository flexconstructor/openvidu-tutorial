package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", index)
	router.Run()
}

func index(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Hallo gin")
}
