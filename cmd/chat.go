/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"net/http"

	"github.com/spf13/cobra"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		parseQuestion(args)
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func parseQuestion(args []string) {
	askQuestion(args)
}

const hfAPIURL = "https://api-inference.huggingface.co/models/bigscience/bloom"

type Payload struct {
	Inputs string `json:"inputs"`
}

func getToken() (string, error) {
	apiToken := os.Getenv("HF_API_TOKEN")

	if apiToken == "" {
		return "", errors.New("HF_API_TOKEN environment variable was not passed")
	}
	return apiToken, nil
}

func askQuestion(question []string) {
	apiToken, _ := getToken()

	payload := Payload{
		Inputs: question[0],
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error serializing payload: %s\n", err)
		return
	}

	req, err := http.NewRequest("POST", hfAPIURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Printf("Error creating request: %s\n", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %s\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %s\n", err)
		return
	}

	fmt.Printf("Response: %s\n", body)
}
