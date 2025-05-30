package api

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/model"
	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/task"
	"github.com/XXueTu/wise/internal/types"
	"github.com/XXueTu/wise/pkg/spiders"
)

type UrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// URL链接识别
func NewUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UrlLogic {
	return &UrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UrlLogic) Url(req *types.URLRequest) (resp *types.URLResponse, err error) {
	resp = &types.URLResponse{}

	// 检查是否是微信公众号链接
	title, content, err := spiders.NewPattern().GetPattern(req.URL)
	if err != nil {
		return resp, err
	}

	err = l.svcCtx.ResourceModel.Create(l.ctx, &model.Resource{
		URL:     req.URL,
		Title:   title,
		Content: content,
		Type:    "微信公众号",
	})
	if err != nil {
		return resp, err
	}
	resp.Title = title
	resp.Description = content
	resp.Link = req.URL
	resp.Tag = []string{"微信公众号"}
	logx.Info("url response:", resp)

	args := task.UrlMarkTaskArgs{
		Url: req.URL,
	}
	argsJson, err := json.Marshal(args)
	err = task.CreateUrlMarkTask(l.ctx, l.svcCtx, task.Args{
		Params: string(argsJson),
	})
	if err != nil {
		return resp, err
	}

	return resp, nil
}
