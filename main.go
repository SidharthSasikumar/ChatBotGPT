package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type ChatGPTResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func main() {
	http.HandleFunc("/chat", chatHandler)
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	// Get user input from query parameters
	userInput := r.URL.Query().Get("input")
	if userInput == "" {
		http.Error(w, "Missing input query parameter", http.StatusBadRequest)
		return
	}
	fmt.Println("Input is " + userInput)

	// Call ChatGPT API
	responseText, err := callChatGPTAPI(userInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response back to user
	fmt.Fprintf(w, "ChatGPT: %s", responseText)
}

func callChatGPTAPI(input string) (string, error) {
	// Construct API request
	url := "https://api.openai.com/v1/engines/text-davinci-003/completions"
	body := fmt.Sprintf(`{"prompt": "%s", "max_tokens": 50}`, input)
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		fmt.Println("Error", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	// Send API request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}
	// Parse API response
	var chatGPTResponse ChatGPTResponse
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error", err)
		return "", err
	}
	fmt.Println("Raw API Response:", string(bodyBytes))
	err = json.Unmarshal(bodyBytes, &chatGPTResponse)
	if err != nil {
		fmt.Println("Error", err)
		return "", err
	}

	fmt.Println("chatGPTResponse:", chatGPTResponse)
	// Return the response text
	if len(chatGPTResponse.Choices) > 0 {
		return chatGPTResponse.Choices[0].Text, nil
	}
	return "", fmt.Errorf("no response from ChatGPT")
}
