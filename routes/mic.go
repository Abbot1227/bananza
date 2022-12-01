package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

var client = &http.Client{}

func LoadAudio(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	//languageParam := c.Params.ByName("lang")
	//language := languageParam[5:]

	// the FormFile function takes in the POST input id file
	c.Request.ParseMultipartForm(32 << 20)

	questionId := c.Request.MultipartForm.Value["id"]
	languageId := c.Request.MultipartForm.Value["languageId"]
	level := c.Request.MultipartForm.Value["level"]
	file, _, err := c.Request.FormFile("mp3")
	if err != nil {
		fmt.Println("Error when requesting file" + err.Error())
		defer cancel()
		return
	}
	defer file.Close()

	defer cancel()

	// TODO remove temp structure
	//var temp map[string]interface{}
	temp := make(map[string]interface{})

	// Forwarding file to AI part
	//if err = sendPostRequest(file, &temp, language); err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	// TODO не забудь убрать
	temp["text"] = "Froehliche Weihnachten"

	var answerStruct bson.D
	question, _ := primitive.ObjectIDFromHex(questionId[0])
	filter := bson.D{{"_id", question}}
	opts := options.FindOne().SetProjection(bson.D{{"_id", 0}, {"answer", 1}})

	if err := tempExercisesCollection.FindOne(ctx, filter, opts).Decode(&answerStruct); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		defer cancel()
		return
	}
	defer cancel()
	rightAnswer := answerStruct.Map()

	exp, _ := strconv.Atoi(level[0])
	expToAdd := calculateGainExp(exp)

	fmt.Println(questionId)
	fmt.Println("Right:", rightAnswer["answer"])
	fmt.Println("User:", temp["text"])

	if temp["text"] == rightAnswer["answer"] {
		c.JSON(http.StatusOK, gin.H{"correct": "true", "answer": rightAnswer["answer"], "exp": expToAdd})
	} else {
		c.JSON(http.StatusOK, gin.H{"correct": "false", "answer": rightAnswer["answer"], "exp": 0})
		return
	}

	langId, _ := primitive.ObjectIDFromHex(languageId[0])

	result, err := userProgressCollection.UpdateByID(ctx, langId, bson.D{
		{"$inc", bson.D{{"level", expToAdd}}},
	})
	if err != nil {
		fmt.Println("Could not add points to user")
	}
	fmt.Println(result)
}

func sendPostRequest(file multipart.File, temp *map[string]interface{}, language string) error {
	b := &bytes.Buffer{}
	w := multipart.NewWriter(b)

	filename := createAudioFilename(language)
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

	cnt, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(cnt, &temp); err != nil {
		fmt.Println("wrong request")
		return err
	}

	return nil
}

// createAudioFilename is a function to generate audio file name
func createAudioFilename(language string) string {
	// language - de, kr
	filename := language + "_audio" + time.Now().Format("01022006150405") + ".mp3"
	return filename
}
