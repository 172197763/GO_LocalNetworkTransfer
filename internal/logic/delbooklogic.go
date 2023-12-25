package logic

import (
	"context"

	"GO_LOCALNETWORKTRANSFER/internal/svc"
	"GO_LOCALNETWORKTRANSFER/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelBookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelBookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelBookLogic {
	return &DelBookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelBookLogic) DelBook(req *types.Request) (resp *types.ResponseCommont, err error) {
	// todo: add your logic here and delete this line

	return
}
