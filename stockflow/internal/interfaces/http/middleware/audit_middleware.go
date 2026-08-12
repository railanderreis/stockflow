package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/railanderreis/stockflow/stockflow/internal/domain/audit"
)

type responseWriterRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *responseWriterRecorder) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriterRecorder) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func AuditMiddleware(auditRepo audit.AuditRepository, action string, resourceType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only audit mutating HTTP methods
			if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			var reqBody []byte
			if r.Body != nil {
				reqBody, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			}

			rec := &responseWriterRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				body:           bytes.NewBuffer(nil),
			}

			next.ServeHTTP(rec, r)

			// Record audit log entry if operation succeeded (2xx / 3xx)
			if rec.statusCode >= 200 && rec.statusCode < 400 {
				actorID := r.Header.Get("X-User-ID")
				actorEmail := r.Header.Get("X-User-Email")
				resourceID := r.Header.Get("X-Resource-ID")

				logEntry := &audit.AuditLog{
					ActorID:      actorID,
					ActorEmail:   actorEmail,
					Action:       action,
					ResourceType: resourceType,
					ResourceID:   resourceID,
					NewValues:    json.RawMessage(reqBody),
					IPAddress:    r.RemoteAddr,
					UserAgent:    r.UserAgent(),
					CreatedAt:    time.Now(),
				}

				_ = auditRepo.Record(r.Context(), logEntry)
			}
		})
	}
}
