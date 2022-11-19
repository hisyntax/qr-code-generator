package qrcode

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

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

	//if the text passed contains a dot(.), it is assumed to be a link
	if strings.Contains(payload.QRCodeText, ".") {
		//validate the url
		_, err := url.ParseRequestURI(payload.QRCodeText)
		if err != nil {
			return "", errors.New("invalid url")
		}
	}

	//generate qrcode
	resImg, _, err := NewRequest(method, requestUrl, payload)
	if err != nil {
		fmt.Printf("This is tthe server err: %v\n", err)
		return "", err
	}

	var base64Encoding string
	mimeType := http.DetectContentType(resImg)
	fmt.Printf("This is the mimeType: %v\n", mimeType)

	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	case "text/plain":
		base64Encoding += "data:text/plain;base64,"
	case "application/pdf":
		base64Encoding += "data:application/pdf;base64,"
	}

	base64Encoding += toBase64(resImg)
	// fmt.Println(base64Encoding)

	imgUrl, err := uploader.FileUploader(base64Encoding)
	if err != nil {
		fmt.Printf("This is tthe cloudinary err: %v\n", err)
		return "", err
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
