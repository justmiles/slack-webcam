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
	devicePtr = flag.String("device", "0", "camera device index")
	debugPtr  = flag.Bool("debug", false, "enable debug logging")
	slackResp SlackResponse
)

type tokensFlag []string

func (i *tokensFlag) String() string {
	return ""
}

func (i *tokensFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var tokens tokensFlag

func main() {
	flag.Var(&tokens, "token", "slack token for auth. Repeat for each Slack account you want to update")
	flag.Parse()
	validateInput(tokens[0], "please provide a slack token")

	deviceID, _ := strconv.Atoi(*devicePtr)

	// Get the image
	debug("opening webcam")
	webcam, err := gocv.VideoCaptureDevice(int(deviceID))
	if err != nil {
		fmt.Printf("error opening video capture device: %v\n", deviceID)
		return
	}

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

	buf, err := gocv.IMEncode(".jpg", img)
	webcam.Close()
	// Upload image
	for i, token := range tokens {
		debug(fmt.Sprintf("token %d of %d: uploading image to Slack", i+1, len(tokens)))

		bodyBuf := &bytes.Buffer{}
		bodyWriter := multipart.NewWriter(bodyBuf)

		bodyWriter.WriteField("token", token)

		fileWriter, err := bodyWriter.CreateFormFile("image", "avatar.jpg")
		check(err)

		contentType := bodyWriter.FormDataContentType()
		_, err = fileWriter.Write(buf)
		check(err)

		bodyWriter.Close()

		resp, err := http.Post("https://slack.com/api/users.setPhoto", contentType, bodyBuf)
		check(err)

		defer resp.Body.Close()
		resp_body, err := ioutil.ReadAll(resp.Body)
		check(err)
		debug(fmt.Sprintf("token %d of %d: %s", i+1, len(tokens), resp.Status))
		err = json.Unmarshal(resp_body, &slackResp)
		check(err)
		if slackResp.Ok {
			debug(fmt.Sprintf("token %d of %d: successfully uploaded", i+1, len(tokens)))
		} else {
			if slackResp.Error != "" {

				switch slackResp.Error {
				case "invalid_auth":
					fmt.Println(fmt.Sprintf("token %d of %d: invalid authentication token provided", i+1, len(tokens)))
				case "not_authed":
					fmt.Println(fmt.Sprintf("token %d of %d: no authentication token provided", i+1, len(tokens)))
				default:
					fmt.Println(slackResp.Error)
				}
			} else {
				fmt.Println(fmt.Sprintf("token %d of %d: an unknown error occurred while communicating with the Slack API", i+1, len(tokens)))
			}
			os.Exit(1)
		}

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
