package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

var client = &http.Client{}

func MicroShit(c *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	// the FormFile function takes in the POST input id file
	c.Request.ParseMultipartForm(32 << 20)

	file, _, err := c.Request.FormFile("mp3")
	if err != nil {
		fmt.Println("Error when requesting file" + err.Error())
		defer cancel()
		return
	}
	defer file.Close()

	defer cancel()

	// Forwarding file to AI part
	if err = sendPostRequest(file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "file uploaded successfully"})
}

func sendPostRequest(file multipart.File) error {
	b := &bytes.Buffer{}
	w := multipart.NewWriter(b)

	filename := createAudioFilename()
	fw, err := w.CreateFormFile("uploaded_file", filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(fw, file)
	if err != nil {
		return err
	}
	w.Close()

	req, err := http.NewRequest("POST", "http://localhost:4040/predict", bytes.NewReader(b.Bytes()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	res, _ := client.Do(req)
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		return err
	}
	defer res.Body.Close()

	// TODO remove temp structure
	var temp map[string]interface{}

	cnt, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(cnt, &temp); err != nil {
		fmt.Println("wrong request")
		return err
	}
	fmt.Println(temp["text"])

	return nil
}

// createAudioFilename is a function to generate audio file name
func createAudioFilename() string {
	filename := "audio" + time.Now().Format("01022006150405") + ".mp3"
	return filename
}
