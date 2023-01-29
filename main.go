package main

import (
	"bufio"
	"context"
	"danieljurek/kids-gpt/config"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	gogpt "github.com/sashabaranov/go-gpt3"
)

var sessionConfig config.Config

func readString() string {
	reader := bufio.NewReader(os.Stdin)

	var input string
	var err error
	for {
		input, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		input = strings.TrimSuffix(input, "\n")
		if input != "" {
			break
		}
	}

	// Remove delimiter
	input = strings.TrimSuffix(input, "\n")
	return input
}

func sayText(text string) {
	say(text, "Samantha")
}

func say(text string, voice string) {
	// Don't write the user's string directly into exec, use a file to prevent
	// injections. There may be happier ways to do this that don't involve as
	// many i/o round trips.
	os.Mkdir("say", os.ModePerm)
	f, err := os.CreateTemp("say", "say-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.WriteString(text)
	f.Close()

	// TODO: run this asynchronously so the UI doesn't appear to hang
	cmd := exec.Command(
		"say",
		"--input-file", f.Name(),
		"--voice", voice,
		"--interactive=/blue",
		"--rate", fmt.Sprint(sessionConfig.Speed))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func moderateGpt3(prompt string, client *gogpt.Client, ctx *context.Context) (bool, error) {
	textModerationLatest := "text-moderation-latest"
	req := gogpt.ModerationRequest{
		Input: prompt,
		Model: &textModerationLatest,
	}

	resp, err := client.Moderations(*ctx, req)
	if err != nil {
		say("There was an error. Please go get dad.", "Zarvox")
		log.Fatal(err)
	}

	for _, result := range resp.Results {
		if result.Flagged {
			return false, errors.New("I will not respond to mean things. Please start over.")
		}
	}
	return true, nil
}

func completeGpt3(prompt string, client *gogpt.Client, ctx *context.Context) (string, error) {
	if _, err := moderateGpt3(prompt, client, ctx); err != nil {
		return "", err
	}

	// TODO: handle error
	config, _ := config.GetConfig()

	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3TextDavinci003,
		MaxTokens: 1024,
		Prompt:    prompt,
		Stop:      []string{config.UserName},
	}
	resp, err := client.CreateCompletion(*ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Text, nil
}

func logConversation(conversation *string) {
	os.Mkdir("conversations", os.ModePerm)
	filename := fmt.Sprintf("conversation-%d.txt", time.Now().Unix())
	f, err := os.Create("conversations/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.WriteString(*conversation)
}

func main() {
	gptClient := gogpt.NewClient(os.Getenv("OPENAI_KEY"))
	ctx := context.Background()

	sessionConfig, err := config.GetConfig()
	if err != nil {
		log.Fatalf("There was an error: %s", err)
	}

	conversation := sessionConfig.InitialPrompt
	stopSequence := sessionConfig.StopSequence
	userName := sessionConfig.UserName
	gptName := sessionConfig.GptName

	// Log conversation in a normal exit
	defer logConversation(&conversation)

	// Log conversation in the event of Ctrl+C (SIGNIT)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		logConversation(&conversation)
		os.Exit(1)
	}()

	var gptResponse, trimmedResponse string

	spinner := spinner.New(
		spinner.CharSets[33],
		250*time.Millisecond)

	for {
		spinner.Start()
		gptResponse, err = completeGpt3(conversation, gptClient, &ctx)
		spinner.Stop()

		if err != nil {
			say("Something went wrong. Go get dad.", "Zarvox")
			log.Fatalf("There was an error: %s", err)
		}

		trimmedResponse = strings.TrimSuffix(
			strings.Trim(gptResponse, "\n"),
			stopSequence)
		fmt.Printf("%s: ", gptName)
		conversation += gptResponse + "\n"
		sayText(trimmedResponse)

		if strings.Contains(gptResponse, stopSequence) {
			break
		}

		fmt.Print("User: ")
		input := readString()
		conversation += fmt.Sprintf(
			"\n%s: %s\n%s: ",
			userName,
			input,
			gptName)
	}

	fmt.Println("Conversation complete.")
}
