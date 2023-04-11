package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/disintegration/imaging"
	"net/http"
	"time"
)

func main() {
	t1 := time.Now()
	url := "http://localhost:7878"

	//open, err := imaging.Open("1.jpg")
	//if err != nil {
	//	panic(err)
	//}
	//buf := new(bytes.Buffer)
	//imaging.Encode(buf, open, imaging.JPEG)
	//toString := base64.StdEncoding.EncodeToString(buf.Bytes())
	payload := map[string]interface{}{
		//"init_images": []interface {
		//}{
		//	toString,
		//},
		"prompt":          "(masterpiece),(best quality),(ultra-detailed), (full body:1.2), 1girl,chibi,cute, smile, white Bob haircut, red eyes, earring, white shirt,black skirt, lace legwear, (sitting on red sofa), seductive posture, smile, A sleek black coffee table sits in front of the sofa and a few decorative items are placed on the shelves, (beautiful detailed face), (beautiful detailed eyes),",
		"negative_prompt": "(low quality:1.3), (worst quality:1.3)",
		"steps":           28,
		"width":           576,
		"height":          768,
		"sampler_name":    "Euler a",
		"cfg_scale":       7.5,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(url+"/sdapi/v1/txt2img", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var response struct {
		Images []string `json:"images"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	for _, b64data := range response.Images {
		imgBytes, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			panic(err)
		}

		img, err := imaging.Decode(bytes.NewReader(imgBytes))
		if err != nil {
			panic(err)
		}

		err = imaging.Save(img, "output.png")
		if err != nil {
			panic(err)
		}
	}
	since := time.Since(t1)
	fmt.Println("Done", since.Seconds())
}
