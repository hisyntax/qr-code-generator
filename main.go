package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hisyntax/qr-code-generator/qrcode"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("no .env found")
	}
}

func main() {
	r := gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	r.POST("/generate-qrcode", GenerateQrCode)

	r.Run(":" + port)
}

func GenerateQrCode(c *gin.Context) {
	var qr qrcode.QrCode
	if err := c.Bind(&qr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error binding json",
		})
		return
	}

	res, err := qrcode.GenerateQRCode(qr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error generating qr code",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": res,
	})
}
