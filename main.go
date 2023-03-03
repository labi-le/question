package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
)

var (
	model       string
	question    string
	token       string
	temperature int
)

func parse() {
	flag.StringVar(&model, "model", "gpt-3.5-turbo-0301", "model")
	flag.StringVar(&question, "q", "", "question")
	flag.StringVar(&token, "token", "", "token")
	flag.IntVar(&temperature, "temperature", 0, "temperature")
	flag.Parse()

	if question == "" || token == "" || model == "" {
		flag.Usage()
		return
	}
}

func main() {
	parse()

	response, err := request(Data{
		Model: model,
		Message: []Message{
			{
				Role:    "user",
				Content: question,
			},
		},
		Temperature: temperature,
	}, token)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, choice := range response.Choices {
		fmt.Println(choice.Message.Content)
	}
}

func request(d Data, token string) (Response, error) {
	var r Response

	jsonData, err := json.Marshal(d)
	if err != nil {
		return r, err
	}

	req, reqErr := http.NewRequest(
		http.MethodPost,
		"https://api.openai.com/v1/chat/completions",
		strings.NewReader(string(jsonData)),
	)
	if reqErr != nil {
		return r, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+string(token))

	client := &http.Client{}
	resp, clientErr := client.Do(req)
	if clientErr != nil {
		return r, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var e ErrorResponse
		if decodeErr := json.NewDecoder(resp.Body).Decode(&e); decodeErr != nil {
			panic(decodeErr)
		}

		return r, fmt.Errorf("error: %s, type: %s, code: %s", e.Error.Message, e.Error.Type, e.Error.Code)
	}

	return r, json.NewDecoder(resp.Body).Decode(&r)
}
