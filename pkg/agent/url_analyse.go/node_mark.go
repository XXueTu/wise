package url_analyse

import (
	"context"
	"strings"

	"github.com/cloudwego/eino/schema"

	"github.com/XXueTu/wise/pkg/model"
)

/*
request:
	{
		"resource_id": 1,
		"segments": [
			"逻辑处理器"，对 G 来说，P 相当于 CPU 核，G 只有绑定到 P 才能被调度。",
			"对 M 来说，P 提供了相关的执行环境(Context)"
		]
	}

response:
	{
		"tags": [
			"golang",
			"algorithm"
		],
		"summarize": "golang 是一种编程语言，算法是一种解决问题的思路",
	}
*/

type Tags struct {
	Tags []string `json:"tags"`
}

func MarkNodeHandler(ctx context.Context, param map[string]any) (map[string]any, error) {
	segments := param["segments"].([]string)
	resourceId := param["resource_id"].(int64)
	summarize, err := llmMark(ctx, segments)
	if err != nil {
		return nil, err
	}
	err = updateResource(ctx, resourceId, summarize["tags"].([]string), summarize["summarize"].(string))
	if err != nil {
		return nil, err
	}
	return summarize, nil
}

func llmMark(ctx context.Context, segments []string) (map[string]any, error) {

	ctModel, err := model.NewChatModel(ctx)
	if err != nil {
		return nil, err
	}
	limitWords := 1000
	totalWords := 0
	historyMessages := []*schema.Message{}
	userMessage := []*schema.Message{}
	for i, segment := range segments {
		totalWords += len(segment)
		userMessage = append(userMessage, schema.UserMessage(segment))
		if totalWords > limitWords || i == len(segments)-1 {
			messages, err := model.ChatPromptSummarize(ctx, userMessage, historyMessages)
			if err != nil {
				return nil, err
			}
			respond, err := ctModel.Generate(ctx, messages)
			if err != nil {
				return nil, err
			}
			historyMessages = append(historyMessages, respond)
			// 清空计数
			totalWords = 0
			userMessage = []*schema.Message{}
		}
	}
	sumarize := historyMessages[len(historyMessages)-1]
	// 根据最后一条历史消息总结,并生成 5 个左右的标签
	labelMessages, err := model.ChatPromptLabel(ctx, []*schema.Message{sumarize})
	if err != nil {
		return nil, err
	}
	tags, err := ctModel.Generate(ctx, labelMessages)
	if err != nil {
		return nil, err
	}
	tagsEntity, err := schema.NewMessageJSONParser[Tags](&schema.MessageJSONParseConfig{
		ParseFrom: schema.MessageParseFromContent,
	}).Parse(ctx, tags)
	if err != nil {
		return nil, err
	}

	// 每达到1000字左右就行总结一次
	return map[string]any{
		"tags":      tagsEntity.Tags,
		"summarize": sumarize.Content,
	}, nil
}

func updateResource(ctx context.Context, resourceId int64, tags []string, describe string) error {
	// 更新 resource 的 tags 和 describe
	resource, err := svcCtx.ResourceModel.Get(ctx, resourceId)
	if err != nil {
		return err
	}
	resource.Tags = strings.Join(tags, ",")
	// 创建 tag

	resource.Describe = describe
	err = svcCtx.ResourceModel.Update(ctx, resource)
	if err != nil {
		return err
	}
	return nil
}
