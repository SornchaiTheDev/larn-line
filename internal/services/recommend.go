package services

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func GetRecommend(message string) ([]string, error) {

	payload := map[string]any{
		"message": message,
	}

	marshalled, err := json.Marshal(payload)

	if err != nil {
		log.Fatal("Cannot marshal json message")
	}

	url := os.Getenv("LARN_API_URL") + "/ai/recommend"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(marshalled))

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

	type Response struct {
		Response []string `json:"response"`
	}

	var apiRes Response

	err = json.Unmarshal(body, &apiRes)

	if err != nil {
		return nil, err
	}

	var finalRecommends []string

	for _, recommend := range apiRes.Response {
		if len(recommend) != 0 {
			finalRecommends = append(finalRecommends, recommend)
		}
	}

	return finalRecommends, nil
}
