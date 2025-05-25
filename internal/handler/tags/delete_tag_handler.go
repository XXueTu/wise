package tags

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/XXueTu/wise/internal/svc"
	"github.com/XXueTu/wise/internal/logic/tags"
	"github.com/XXueTu/wise/internal/types"
	"github.com/XXueTu/wise/response"
)

func DeleteTagHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteTagRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := tags.NewDeleteTagLogic(r.Context(), svcCtx)
		resp, err := l.DeleteTag(&req)
		response.Response(w, resp, err)

	}
}
