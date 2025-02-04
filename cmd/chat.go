/*
Copyright © 2024 Anita Bendelja @anitabee
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	huggingface "github.com/hupe1980/go-huggingface"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var modelName = viper.GetString("model")

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
		if len(args) == 0 {
			log.Fatalf("Please provide a question")
		}
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

func getClient() (*huggingface.InferenceClient, error) {
	apiToken := os.Getenv("HF_API_TOKEN")

	if apiToken == "" {
		return nil, errors.New("HF_API_TOKEN environment variable was not passed")
	}
	client := huggingface.NewInferenceClient(apiToken)
	return client, nil
}

func parseQuestion(question []string) {
	client, err := getClient()
	if err != nil {
		log.Fatalf("Error reaching Hugging Face API or model %s: %v", modelName, err)

	}

	maxNewTokens := 30
	topP := 0.8
	repetitionPenalty := 1.2
	temperature := 0.7
	numReturnSequences := 1

	request := &huggingface.TextGenerationRequest{
		Inputs: question[0],
		Parameters: huggingface.TextGenerationParameters{
			MaxNewTokens:       &maxNewTokens,
			TopP:               &topP,
			RepetitionPenalty:  &repetitionPenalty,
			Temperature:        &temperature,
			NumReturnSequences: &numReturnSequences,
		},
		Model: modelName,
	}
	response, err := client.TextGeneration(context.Background(), request)
	if err != nil {
		log.Fatalf("Error generating text: %v", err)
	}
	fmt.Println(response[0].GeneratedText)

}
