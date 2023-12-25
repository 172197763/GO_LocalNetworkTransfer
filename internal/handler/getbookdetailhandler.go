package handler

import (
	"net/http"

	"GO_LOCALNETWORKTRANSFER/internal/logic"
	"GO_LOCALNETWORKTRANSFER/internal/svc"
	"GO_LOCALNETWORKTRANSFER/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetBookDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetBookDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetBookDetail(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
