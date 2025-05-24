package wechat

import (
	"context"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type PublicAccounts struct {
}

const (
	WechatUrl = "https://mp.weixin.qq.com/"
)

func Init() *PublicAccounts {
	return &PublicAccounts{}
}

func (l *PublicAccounts) Identification(url string) bool {
	return strings.HasPrefix(url, WechatUrl)
}

// 获取微信公众号文章内容
func (l *PublicAccounts) GetData(url string) (string, string, error) {
	// 创建新的Chrome实例
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 设置超时时间
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var title, content string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		// 等待页面加载完成
		chromedp.WaitVisible(`.rich_media_title`, chromedp.ByQuery),
		chromedp.WaitVisible(`.rich_media_content`, chromedp.ByQuery),
		// 获取标题
		chromedp.Text(`.rich_media_title`, &title, chromedp.ByQuery),
		// 获取内容
		chromedp.Text(`.rich_media_content`, &content, chromedp.ByQuery),
	)
	if err != nil {
		return "", "", err
	}

	// 清理内容
	content = l.cleanContent(content)

	return strings.TrimSpace(title), content, nil
}

// 清理内容
func (l *PublicAccounts) cleanContent(content string) string {
	// 移除多余的空白字符
	content = strings.TrimSpace(content)

	// 替换常见的HTML实体
	content = strings.ReplaceAll(content, "&nbsp;", " ")
	content = strings.ReplaceAll(content, "&lt;", "<")
	content = strings.ReplaceAll(content, "&gt;", ">")
	content = strings.ReplaceAll(content, "&amp;", "&")
	content = strings.ReplaceAll(content, "&quot;", "\"")

	// 移除图片描述
	content = strings.ReplaceAll(content, "图片", "")

	// 移除多余的空行
	lines := strings.Split(content, "\n")
	var cleanedLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanedLines = append(cleanedLines, line)
		}
	}

	return strings.Join(cleanedLines, "\n")
}
