package url_analyse

import (
	"context"

	"github.com/XXueTu/wise/pkg/spiders"
)

/*
request:
	{
		"url": "https://mp.weixin.qq.com/s/1234567890"
	}

response:
	{
		"url": "https://mp.weixin.qq.com/s/1234567890",
		"types": "wechat"
	}
*/

func CheckNodeHandler(ctx context.Context, param map[string]any) (map[string]any, error) {
	url := param["url"].(string)
	pattern := spiders.NewPattern().GetPatternTypes(url)
	return map[string]any{
		"url":   url,
		"types": pattern,
	}, nil
}
