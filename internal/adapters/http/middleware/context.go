package middleware

import (
	"net/http"
	"os"

	"github.com/google/uuid"
	zapctx "github.com/saltpay/go-zap-ctx"
	"go.uber.org/zap"
)

// ContextMiddleware adds service context to zapcontext.
func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = zapctx.WithFields(ctx,
			zap.String("env", os.Getenv("ENV_NAME")),
			zap.String("service", os.Getenv("SERVICE")),
			zap.Stringer("trace.id", uuid.New()),
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
