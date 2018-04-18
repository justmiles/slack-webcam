package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"gocv.io/x/gocv"
)

var (
	tokenPtr  = flag.String("token", "", "slack token for auth")
	devicePtr = flag.String("device", "0", "camera device index")
	debugPtr  = flag.Bool("debug", false, "enable debug logging")
	slackResp SlackResponse
)

func main() {
	flag.Parse()
	validateInput(*tokenPtr, "please provide a slack token")

	deviceID, _ := strconv.Atoi(*devicePtr)

	// Get the image
	debug("opening webcam")
	webcam, err := gocv.VideoCaptureDevice(int(deviceID))
	if err != nil {
		fmt.Printf("error opening video capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	img := gocv.NewMat()
	defer img.Close()

	debug("reading webcam")
	if ok := webcam.Read(&img); !ok {
		fmt.Printf("cannot read device %d\n", deviceID)
		return
	}
	if img.Empty() {
		fmt.Printf("no image on device %d\n", deviceID)
		return
	}

	// Upload image
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	bodyWriter.WriteField("token", *tokenPtr)

	fileWriter, err := bodyWriter.CreateFormFile("image", "avatar.jpg")
	check(err)

	contentType := bodyWriter.FormDataContentType()
	buf, err := gocv.IMEncode(".jpg", img)
	_, err = fileWriter.Write(buf)
	check(err)

	bodyWriter.Close()

	debug("uploading image to slack")
	resp, err := http.Post("https://slack.com/api/users.setPhoto", contentType, bodyBuf)
	check(err)

	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	check(err)
	debug(resp.Status)
	debug(string(resp_body))
	err = json.Unmarshal(resp_body, &slackResp)
	check(err)
	if slackResp.Ok {
		debug("successfully uploaded")
	} else {
		if slackResp.Error != "" {
			fmt.Println(slackResp.Error)
		} else {
			fmt.Println("an unknown error occurred when communicating with the Slack API")
		}
		os.Exit(1)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func validateInput(s string, msg string) {
	if s == "" {
		fmt.Println(msg)
		os.Exit(1)
	}
}

func debug(msg string) {
	if *debugPtr {
		fmt.Printf("[DEBUG] %s\n", msg)
	}
}

type SlackResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}
