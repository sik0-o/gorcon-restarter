package server

import (
	"encoding/json"
	"fmt"

	"github.com/sik0-o/gorcon-restarter/v2/internal/service/action"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func (ass *Server) Announce(ctx *fasthttp.RequestCtx) {
	reqBody := struct {
		Message string
	}{}

	if err := json.Unmarshal(ctx.Request.Body(), &reqBody); err != nil {
		ass.logger.Debug("failed to parse request",
			zap.String("handler", "Announce"),
			zap.Error(err),
		)
		ctx.Error(fmt.Sprintf("failed to parse request: %s", err.Error()), fasthttp.StatusBadRequest)
		return
	}

	if err := ass.cs.Exec(action.NewAnnounce(reqBody.Message)); err != nil {
		ass.logger.Debug("failed to send announce",
			zap.String("handler", "Announce"),
			zap.Error(err),
		)
		ctx.Error(fmt.Sprintf("failed to send announce: %s", err.Error()), fasthttp.StatusInternalServerError)
		return
	}
}
