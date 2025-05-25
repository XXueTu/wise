package tags

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
)

type DeleteTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除标签
func NewDeleteTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTagLogic {
	return &DeleteTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTagLogic) DeleteTag(req *types.DeleteTagRequest) (resp *types.DeleteTagResponse, err error) {
	// 查找是否存在
	tag, err := l.svcCtx.TagsModel.GetUid(l.ctx, req.Uid)
	if err != nil {
		return nil, errors.New("查询标签失败")
	}
	if tag == nil {
		return nil, errors.New("标签不存在")
	}
	err = l.svcCtx.TagsModel.Delete(l.ctx, tag.ID)
	if err != nil {
		return nil, errors.New("删除标签失败")
	}
	resp = &types.DeleteTagResponse{
		Result: "删除标签成功",
	}
	return
}
