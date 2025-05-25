package resources

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/model"
	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
	"github.com/XXueTu/wise/pkg/spiders"
)

type CreateAiResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建AI资源
func NewCreateAiResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAiResourceLogic {
	return &CreateAiResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAiResourceLogic) CreateAiResource(req *types.CreateAiResourceRequest) (resp *types.Resource, err error) {
	// 检查是否是微信公众号链接
	title, content, err := spiders.NewPattern().GetPattern(req.URL)
	if err != nil {
		return resp, err
	}

	resource := model.Resource{
		URL:     req.URL,
		Title:   title,
		Content: content,
		Type:    "微信公众号",
		Tags:    "default",
	}
	err = l.svcCtx.ResourceModel.Create(l.ctx, &resource)
	if err != nil {
		return resp, err
	}
	resp = &types.Resource{
		Id:        resource.ID,
		URL:       req.URL,
		Title:     title,
		Content:   content,
		Type:      "微信公众号",
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	return
}
