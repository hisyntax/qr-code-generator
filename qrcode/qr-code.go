package qrcode

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	urlReq "github.com/hisyntax/domain/url"

	"github.com/hisyntax/qr-code-generator/uploader"
)

type QrCode struct {
	FrameName   string `json:"frame_name"`
	QRCodeText  string `json:"qr_code_text"`
	ImageFormat string `json:"image_format"`
	QRCodeLogo  string `json:"qr_code_logo"`
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func GenerateQRCode(payload QrCode) (string, error) {
	api_key := os.Getenv("API_KEY")
	requestUrl := fmt.Sprintf("https://api.qr-code-generator.com/v1/create?access-token=%s", api_key)
	method := "POST"

	//validate url
	texttype, err := urlReq.ValidateURL(payload.QRCodeText)
	if err != nil {
		return "", err
	}

	//generate qrcode
	resImg, _, err := NewRequest(method, requestUrl, payload)
	if err != nil {
		fmt.Printf("This is the server err: %v\n", err)
		return "", errors.New("error generating qr code")
	}

	var base64Encoding string
	mimeType := http.DetectContentType(resImg)
	fmt.Printf("This is the mimeType: %v\n", mimeType)

	var fileFormat string
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
		fileFormat = "jpeg"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
		fileFormat = "png"
	case "text/plain":
		base64Encoding += "data:text/plain;base64,"
		fileFormat = "svg"
	case "application/pdf":
		base64Encoding += "data:application/pdf;base64,"
		fileFormat = "pdf"
	}

	base64Encoding += toBase64(resImg)
	fileName := fmt.Sprintf("%s_%s_%v", texttype, payload.QRCodeText, fileFormat) // type(url/text), data type(google/ https://google.com), extention(file format, - png, jpg)

	imgUrl, err := uploader.FileUploader(base64Encoding, fileName)
	if err != nil {
		fmt.Printf("This is the cloudinary err: %v\n", err)
		return "", errors.New("qr code could not be uploaded")
	}

	return imgUrl, err
}

func NewRequest(method, url string, payload interface{}) ([]byte, int, error) {
	client := http.Client{}

	jsonReq, jsonReqErr := json.Marshal(&payload)
	if jsonReqErr != nil {
		return nil, 0, jsonReqErr
	}

	req, reqErr := http.NewRequest(method, url, bytes.NewBuffer(jsonReq))
	if reqErr != nil {
		return nil, 0, reqErr
	}

	req.Header.Add("Content-Type", "application/json")

	resp, respErr := client.Do(req)
	if respErr != nil {
		return nil, 0, respErr
	}

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	return resp_body, resp.StatusCode, nil
}
