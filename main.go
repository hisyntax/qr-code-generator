package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/hisyntax/qr-code-generator/docs"
	"github.com/hisyntax/qr-code-generator/qrcode"

	"github.com/gin-gonic/gin"
	urlReq "github.com/hisyntax/domain/url"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("no .env found")
	}

	urlReq.Token = os.Getenv("TOKEN")
}

// @title           qrcode generator API
// @version         1.0
// @description     This is the API docs for testing
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host                       customqr.herokuapp.com
// @BasePath
// @schemes                    https
// @query.collection.format    multi
// @securityDefinitions.basic  BasicAuth
func main() {
	r := gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	r.POST("/generate-qrcode", GenerateQrCode)

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":" + port)
}

// generate qrcode 	godoc
// @Summary      generate qrcode
// @Description  use this endpoint to generate a qr code . This is an example request payload "frame_name": "no-frame",  "qr_code_logo": "scan-me-square"(this is be optional),  "image_format": "PNG, PDF, JPG",  "qr_code_text": "https://google.com",
// @Tags         qr-code
// @Accept       json
// @Produce      json
// @param        qrcode  body  qrcode.QrCode  true  "qrcode"
// @Success      200
// @Router       /generate-qrcode [post]
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
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": res,
	})
}
