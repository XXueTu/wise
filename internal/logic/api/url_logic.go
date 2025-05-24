package api

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/model"
	"github.com/XXueTu/wise/internal/svc"
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

	resource := model.NewResource(req.URL, title, content, "微信公众号")
	err = resource.Create(l.ctx, resource)
	if err != nil {
		return resp, err
	}
	resp.Title = title
	resp.Description = content
	resp.Link = req.URL
	resp.Tag = []string{"微信公众号"}
	logx.Info("url response:", resp)
	return resp, nil
}
