package logic

import (
	"context"

	"GO_LOCALNETWORKTRANSFER/internal/svc"
	"GO_LOCALNETWORKTRANSFER/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBookDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBookDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBookDetailLogic {
	return &GetBookDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBookDetailLogic) GetBookDetail(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
