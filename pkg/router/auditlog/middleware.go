package audit

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"zq-xu/warehouse-admin/internal/webserver/types"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	bc chan []byte
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.bc <- b
	return w.ResponseWriter.Write(b)
}

type AuditorMiddleware struct {
	parser map[string]AuditLog
}

func (aw *AuditorMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		aw.middlewareImpl(ctx)
	}
}

func (aw *AuditorMiddleware) middlewareImpl(ctx *gin.Context) {
	skipAudit := aw.isSkipAudit(ctx)
	if skipAudit {
		ctx.Next()
		return
	}

	bc := make(chan []byte)
	ctx.Writer = &bodyLogWriter{
		ResponseWriter: ctx.Writer,
		bc:             bc,
	}
	ctx.Next()

	go record(ctx, bc)
}

func (aw *AuditorMiddleware) isSkipAudit(c *gin.Context) bool {
	method := c.Request.Method
	if method == http.MethodGet {
		return true
	}

	username := c.GetString(types.AuthUserNameToken)
	if username == "" {
		return true
	}

	return false
}

func record(ctx *gin.Context, bc <-chan []byte) {
	bs := <-bc
	url := ctx.Request.URL.Path

	message := fmt.Sprintf("%v", returnJson.Message)

	adminInfo := admin.(**service.RunningClaims)
	adminId := (*adminInfo).ID
	adminName := (*adminInfo).Account

	var log = model.AdminLog{
		AdminId:   adminId,
		AdminName: adminName,
		Method:    method,
		Url:       url,
		Ip:        util.RemoteIP(c.Request),
		Code:      returnJson.Code,
		Message:   message,
	}

	model.CreateLog(log)
}
