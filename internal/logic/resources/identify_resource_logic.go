package resources

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/task"
	"github.com/XXueTu/wise/internal/types"
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
		args := task.UrlMarkTaskArgs{
			Url: url,
		}
		argsJson, err := json.Marshal(args)

		err = task.CreateUrlMarkTask(l.ctx, l.svcCtx, task.Args{
			Params: string(argsJson),
		})
		if err != nil {
			return resp, err
		}
		resp.Urls = append(resp.Urls, url)
	}
	logx.Info("identify resource urls:", strings.Join(resp.Urls, ","))
	return resp, nil
}
