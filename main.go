package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	gogpt "github.com/sashabaranov/go-gpt3"
	"gopkg.in/yaml.v2"
)

func readString() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	// Remove delimiter
	input = strings.TrimSuffix(input, "\n")
	return input
}

func sayText(text string) {
	say(text, "Samantha")
}

func say(text string, voice string) {
	// TODO: put this in a temp folder
	f, err := os.CreateTemp("/tmp", "say-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.WriteString(text)
	f.Close()

	// TODO: run this asynchronously so the UI doesn't appear to hang
	cmd := exec.Command("say", "--input-file", f.Name(), "--voice", voice)

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func completeGpt3(prompt string, client *gogpt.Client, ctx *context.Context) string {
	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3TextDavinci003,
		MaxTokens: 1024,
		Prompt:    prompt,
	}
	resp, err := client.CreateCompletion(*ctx, req)
	if err != nil {
		say("There was an error. Please go get dad.", "Zarvox")
		log.Fatal(err)
	}

	return resp.Choices[0].Text
}

type Config struct {
	InitialPrompt  string
	SpinnerCharset int
	StopSequence   string
	UserName       string
	GptName        string
}

func logConversation(conversation string) {
	filename := fmt.Sprintf("conversation-%d.txt", time.Now().Unix())
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.WriteString(conversation)
}

func main() {
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)

	gptClient := gogpt.NewClient(os.Getenv("OPENAI_KEY"))
	ctx := context.Background()

	var conversation, gptResponse, trimmedResponse string
	conversation = config.InitialPrompt
	defer logConversation(conversation)

	spinner := spinner.New(
		spinner.CharSets[config.SpinnerCharset],
		250*time.Millisecond)

	for true {
		spinner.Start()
		gptResponse = completeGpt3(conversation, gptClient, &ctx)
		spinner.Stop()

		trimmedResponse = strings.TrimSuffix(
			strings.Trim(gptResponse, "\n"),
			config.StopSequence)
		fmt.Printf("%s: %s\n", config.GptName, trimmedResponse)
		conversation += gptResponse + "\n"
		sayText(trimmedResponse)

		if strings.Contains(gptResponse, config.StopSequence) {
			break
		}

		fmt.Print("User: ")
		input := readString()
		conversation += fmt.Sprintf(
			"\n%s: %s\n%s: ",
			config.UserName,
			input,
			config.GptName)
	}

	fmt.Println("Conversation complete.")
}
