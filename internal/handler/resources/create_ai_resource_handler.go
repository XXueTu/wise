package resources

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/XXueTu/wise/internal/logic/resources"
	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/types"
	"github.com/XXueTu/wise/response"
)

func CreateAiResourceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateAiResourceRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := resources.NewCreateAiResourceLogic(r.Context(), svcCtx)
		resp, err := l.CreateAiResource(&req)
		response.Response(w, resp, err)

	}
}
