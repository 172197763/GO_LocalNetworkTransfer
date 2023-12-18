package handler

import (
	"net/http"

	"GOZEROTEST/internal/logic"
	"GOZEROTEST/internal/svc"
	"GOZEROTEST/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetBookListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetBookListLogic(r.Context(), svcCtx)
		resp, err := l.GetBookList(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
