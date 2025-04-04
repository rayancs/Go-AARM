package repo

import (
	"app/configs"
	"context"
	"encoding/json"
	"errors"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type IOpenAIRepo interface {
}
type OpenAIRepo struct {
	client *openai.Client
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
