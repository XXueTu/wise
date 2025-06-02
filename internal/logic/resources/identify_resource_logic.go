package resources

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/task"
	"github.com/XXueTu/wise/internal/types"
	"github.com/XXueTu/wise/pkg/agent/url_analyse.go"
)

type IdentifyResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// AI识别资源
func NewIdentifyResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IdentifyResourceLogic {
	return &IdentifyResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IdentifyResourceLogic) IdentifyResource(req *types.IdentifyResourceRequest) (resp *types.IdentifyResourceResponse, err error) {
	// 去除空格,\n,\r
	urls := strings.Split(strings.TrimSpace(req.URL), ",")
	resp = &types.IdentifyResourceResponse{
		Urls: make([]string, 0),
	}
	for _, url := range urls {
		if url == "" {
			continue
		}
		_ = task.CreateTask(l.ctx, l.svcCtx, url, "解析URL", "URL_ANALYSE", int64(len(url_analyse.UrlAnalyseSteps)))
		resp.Urls = append(resp.Urls, url)
	}
	logx.Info("identify resource urls:", strings.Join(resp.Urls, ","))
	return resp, nil
}
