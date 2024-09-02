package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type LLMRequest struct {
	Prompt string `json:"prompt"`
}

type LLMResponse struct {
	Response string `json:"response"`
}

func callLLMAPI(prompt string) (string, error) {
	apiKey := os.Getenv("LLM_API_KEY")
	apiURL := os.Getenv("LLM_API_URL")

	reqBody, err := json.Marshal(LLMRequest{Prompt: prompt})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var llmResp LLMResponse
	err = json.NewDecoder(resp.Body).Decode(&llmResp)
	if err != nil {
		return "", err
	}

	return llmResp.Response, nil
}
