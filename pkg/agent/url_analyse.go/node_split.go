package url_analyse

import (
	"context"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/recursive"
	"github.com/cloudwego/eino/schema"
)

/*
request:
	{
		"resource_id": 1,
		"content": "逻辑处理器，对 G 来说，P 相当于 CPU 核，G 只有绑定到 P 才能被调度。对 M 来说，P 提供了相关的执行环境(Context)",
	}

response:
	{
		"resource_id": 1,
		"segments": [
			"逻辑处理器"，对 G 来说，P 相当于 CPU 核，G 只有绑定到 P 才能被调度。",
			"对 M 来说，P 提供了相关的执行环境(Context)"
		]
	}
*/

func SplitNodeHandler(ctx context.Context, param map[string]any) (map[string]any, error) {
	content := param["content"].(string)
	resourceId := param["resource_id"].(int64)
	// 初始化分割器
	splitter, err := recursive.NewSplitter(ctx, &recursive.Config{
		ChunkSize:   1000,
		OverlapSize: 200,
		Separators:  []string{"\n\n", "\n", "。", "！", "？", ";"},
		KeepType:    recursive.KeepTypeEnd,
	})
	if err != nil {
		panic(err)
	}

	// 准备要分割的文档
	docs := []*schema.Document{
		{
			ID:      "doc",
			Content: content,
		},
	}

	// 执行分割
	results, err := splitter.Transform(ctx, docs)
	if err != nil {
		panic(err)
	}
	segments := make([]string, len(results))
	// 处理分割结果
	for i, doc := range results {
		segments[i] = doc.Content
	}
	return map[string]any{
		"resource_id": resourceId,
		"segments":    segments,
	}, nil
}
