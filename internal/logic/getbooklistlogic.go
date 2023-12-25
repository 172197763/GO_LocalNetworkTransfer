package logic

import (
	"context"

	"GO_LOCALNETWORKTRANSFER/internal/svc"
	"GO_LOCALNETWORKTRANSFER/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBookListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBookListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBookListLogic {
	return &GetBookListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBookListLogic) GetBookList(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
