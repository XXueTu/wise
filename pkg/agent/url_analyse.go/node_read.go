package url_analyse

import (
	"context"

	"github.com/XXueTu/wise/internal/model"
	"github.com/XXueTu/wise/pkg/spiders"
)

/*
request:

	{
		"url": "https://mp.weixin.qq.com/s/1234567890"
	}

response:

	{
		"resource_id": 1,
		"url": "https://mp.weixin.qq.com/s/1234567890",
		"types": "wechat",
		"title": "标题",
		"content": "内容"
	}
*/
func ReadNodeHandler(ctx context.Context, param map[string]any) (map[string]any, error) {
	url := param["url"].(string)
	types := param["types"].(string)
	title, content, err := spiders.NewPattern().GetPattern(param["url"].(string))
	if err != nil {
		return nil, err
	}
	resource := &model.Resource{
		URL:     url,
		Title:   title,
		Content: content,
		Type:    types,
	}
	err = svcCtx.ResourceModel.Create(ctx, resource)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"resource_id": resource.ID,
		"url":         param["url"],
		"types":       param["types"],
		"title":       title,   // 标题
		"content":     content, // 内容
	}, nil
}
