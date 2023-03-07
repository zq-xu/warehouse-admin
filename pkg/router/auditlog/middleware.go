package auditlog

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"zq-xu/warehouse-admin/pkg/log"
	"zq-xu/warehouse-admin/pkg/router/auth"
	"zq-xu/warehouse-admin/pkg/store"
)

var (
	Middleware = NewAuditorMiddleware()
)

type bodyWriter struct {
	gin.ResponseWriter
	bodyBuf *bytes.Buffer
}

type auditLogMiddleware struct{}

func NewBodyWriter(ctx *gin.Context) *bodyWriter {
	return &bodyWriter{
		ResponseWriter: ctx.Writer,
		bodyBuf:        bytes.NewBufferString(""),
	}
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.bodyBuf.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *bodyWriter) Bytes() []byte {
	return w.bodyBuf.Bytes()
}

func NewAuditorMiddleware() *auditLogMiddleware {
	return &auditLogMiddleware{}
}

func (aw *auditLogMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		aw.middlewareImpl(ctx)
	}
}

func (aw *auditLogMiddleware) middlewareImpl(ctx *gin.Context) {
	skipAudit := aw.isSkipAudit(ctx)
	if skipAudit {
		ctx.Next()
		return
	}

	reqBody := aw.extractRequestBody(ctx)

	w := NewBodyWriter(ctx)
	ctx.Writer = w
	ctx.Next()

	aw.recordAuditLog(NewModelAuditLog(ctx, reqBody, w.Bytes()))
}

func (aw *auditLogMiddleware) extractRequestBody(ctx *gin.Context) []byte {
	if ctx.Request.Body == nil {
		return []byte{}
	}

	bodyBytes, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes
}

func (aw *auditLogMiddleware) isSkipAudit(ctx *gin.Context) bool {
	method := ctx.Request.Method
	if method != http.MethodPost &&
		method != http.MethodPut &&
		method != http.MethodDelete {
		return true
	}

	id := ctx.GetString(auth.AuthUserIDToken)
	if id == "" {
		return true
	}

	return false
}

func (aw *auditLogMiddleware) recordAuditLog(m *ModelAuditLog) {
	if m.StatusCode >= http.StatusBadRequest {
		return
	}

	db := store.DB(context.Background())
	err := store.Create(db, m)
	if err != nil {
		log.Logger.Errorf("Failed to create the audit log. %v", err)
	}
}
