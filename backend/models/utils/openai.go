package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Create the coderequest type for request
type CodeRequest struct {
	Topic string `json:"topic"`
}

// Create the type of the answer that we want to get
type CodeImplementation struct {
	Language    string `json:"language"`
	Code        string `json:"code"`
	Explanation string `json:"explanation"`
}

// Create teh type of Openai Response
type OpenAIResponse struct {
	Topic           string               `json:"topic"`
	Implementations []CodeImplementation `json:"implementations"`
}

func CallOpenAI(topic string) (*OpenAIResponse, error) {
	prompt := fmt.Sprintf(` You are a helpful AI assistant. Please provide implementation of the topic"%s"in C++, Python, Go, and Javascript. For each language, return the full code that can be run directly and a detailed explanation in JSON format.`, topic)
	payload := map[string]interface{}{
		"model": "gpt-4o-2024-08-06",
		"messages": []map[string]string{
			{"role": "system", "content": "you are a code generation assistant"},
			{"role": "user", "content": prompt},
		},
		"response_format": "json",
	}
	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.charanywhere.tech/v1/chat/completions", bytes.NewBuffer(jsonData))

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil.fmt.Errorf("missing OPENAI_API_KEY environment variable")
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.DefaultClient{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	var pased OpenAIResponse
	err = json.Unmarshal([]byte(result.Choices[0].Message.Content), &parsed)
	return &pased, nil
}
