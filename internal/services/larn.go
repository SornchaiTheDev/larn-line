package services

import (
	"bytes"
	"encoding/json"
	"io"
	"larn-line/internal/models"
	"log"
	"net/http"
	"os"
)

func GetLarn(message string) (*models.Message, error) {
	payload := map[string]string{
		"message": message,
	}

	marshalled, err := json.Marshal(payload)

	if err != nil {
		log.Fatal("Cannot marshal json message")
	}

	req, err := http.NewRequest("POST", os.Getenv("LARN_API_URL"), bytes.NewBuffer(marshalled))

	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Fatal("Request Failed")
	}

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	var response models.Message

	err = json.Unmarshal(body, &response)

	if err != nil {
		return nil, err
	}

	return &response, nil

}
