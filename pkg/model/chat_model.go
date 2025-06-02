package model

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type ChatModelImpl struct {
	config *ChatModelConfig
	model  *openai.ChatModel
}

type ChatModelConfig struct {
	BaseURL string
	APIKey  string
	Model   string
}

// newChatModel component initialization function of node 'CustomChatModel1' in graph 'dev'
func NewChatModel(ctx context.Context) (cm model.ToolCallingChatModel, err error) {
	config := &ChatModelConfig{
		BaseURL: os.Getenv("DEFAULT_BASE_URL"),
		APIKey:  os.Getenv("DEFAULT_API_KEY"),
		Model:   os.Getenv("DEFAULT_MODEL"),
	}
	opcm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: config.BaseURL,
		APIKey:  config.APIKey,
		Model:   config.Model,
	})
	if err != nil {
		return nil, err
	}
	cm = &ChatModelImpl{config: config, model: opcm}
	return cm, nil
}

// Generate implements model.ToolCallingChatModel.
func (c *ChatModelImpl) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	return c.model.Generate(ctx, input, opts...)
}

// Stream implements model.ToolCallingChatModel.
func (c *ChatModelImpl) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	return c.model.Stream(ctx, input, opts...)
}

// WithTools implements model.ToolCallingChatModel.
func (c *ChatModelImpl) WithTools(tools []*schema.ToolInfo) (model.ToolCallingChatModel, error) {
	return c.model.WithTools(tools)
}
