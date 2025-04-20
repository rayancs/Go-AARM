package repo

import (
	"app/configs"
	"app/logger"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type AiRepo[T any] struct {
	client *openai.Client
}

type OpenAIRepo struct {
	client *openai.Client
}

func NewStructAIRepo[T any](connStr string) *AiRepo[T] {
	cl := openai.NewClient(
		option.WithAPIKey(connStr),
	)

	return &AiRepo[T]{
		client: &cl,
	}
}

func NewOpenAi(connStr string) *OpenAIRepo {
	cl := openai.NewClient(
		option.WithAPIKey(connStr),
	)

	return &OpenAIRepo{
		client: &cl,
	}
}

func (a *OpenAIRepo) PromptText(prompt string) (string, error) {

	chatCompletion, err := a.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: openai.ChatModelGPT4o,
	})

	if err != nil {
		return "", err
	}

	if len(chatCompletion.Choices) == 0 {
		return "", nil
	}

	// Extract the response text from the completion
	responseText := chatCompletion.Choices[0].Message.Content

	return responseText, nil
}

func (o *OpenAIRepo) TranscribeAudio(fPath string) (string, error) {
	file, err := os.Open(fPath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	transcription, err := o.client.Audio.Transcriptions.New(context.Background(),
		openai.AudioTranscriptionNewParams{
			Model:  openai.AudioModelWhisper1,
			File:   file,
			Prompt: openai.String("Conver this into indepth notes"),
		})
	if err != nil {
		return "", err
	}
	return transcription.Text, nil
}

func (o *OpenAIRepo) Speech(data string, out string) error {
	log := logger.New()
	res, err := o.client.Audio.Speech.New(context.Background(), openai.AudioSpeechNewParams{
		Model: openai.SpeechModelGPT4oMiniTTS,
		Input: data,
		Voice: openai.AudioSpeechNewParamsVoiceAsh,
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Create the output file
	outputPath := out + ".mp3"

	log.Info("saving_audio@" + outputPath)

	outputFile, err := os.Create(outputPath)

	if err != nil {
		return err
	}

	defer outputFile.Close()
	// Copy the audio data from the response to the file
	_, err = io.Copy(outputFile, res.Body)
	if err != nil {
		return err
	}
	return nil
}

// params
//
//	[]openai.ChatCompletionMessageParamUnion
func StructuredText[T any](a *OpenAIRepo, prompt []openai.ChatCompletionMessageParamUnion, schemaName string) (interface{}, error) {

	schmeaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:   schemaName,
		Schema: configs.GenerateSchema[T](),
		Strict: openai.Bool(true),
	}

	chat, err := a.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schmeaParam,
			},
		},
		Messages: prompt,
		Model:    openai.ChatModelGPT4o2024_08_06,
	})

	if err != nil {
		return nil, err
	}

	if chat == nil {
		return nil, errors.New("nil_chat")
	}

	var resVar T
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &resVar)
	if err != nil {
		return nil, err
	}

	return resVar, nil
}

func (a *AiRepo[T]) StructuredOut(prompt []openai.ChatCompletionMessageParamUnion, schemaName string, temp float64) (*T, error) {
	schmeaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:   schemaName,
		Schema: configs.GenerateSchema[T](),
		Strict: openai.Bool(true),
	}

	chat, err := a.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schmeaParam,
			},
		},
		Messages:    prompt,
		Model:       openai.ChatModelGPT4o2024_08_06,
		Temperature: openai.Float(temp),
	})

	if err != nil {
		return nil, err
	}

	if chat == nil {
		return nil, errors.New("nil_chat")
	}

	var resVar T
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &resVar)
	if err != nil {
		return nil, err
	}

	return &resVar, nil
}
