package qrcode

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
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
	res, _, err := NewRequest(method, requestUrl, payload)
	if err != nil {
		fmt.Printf("This is tthe server err: %v\n", err)
		return "", err
	}

	//write image byte to disk
	img, _, err := image.Decode(bytes.NewReader(res))
	if err != nil {
		return "", err
	}
	out, err := os.Create("./QRImg.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = png.Encode(out, img)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	f, err := os.Open("./QRImg.png")
	if err != nil {
		return "", err
	}

	fmt.Printf("This is the file opened: %v\n", f)
	// files, err := f.Readdir(0)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	imgUrl, err := uploader.FileUploader(f)
	if err != nil {
		fmt.Printf("This is tthe cloudinary err: %v\n", err)
		return "", err
	}

	//upload the image to cloudinary

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
