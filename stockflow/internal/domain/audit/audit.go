package audit

import (
	"context"
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID           string          `json:"id"`
	ActorID      string          `json:"actor_id"`
	ActorEmail   string          `json:"actor_email"`
	Action       string          `json:"action"`
	ResourceType string          `json:"resource_type"`
	ResourceID   string          `json:"resource_id"`
	OldValues    json.RawMessage `json:"old_values,omitempty"`
	NewValues    json.RawMessage `json:"new_values,omitempty"`
	IPAddress    string          `json:"ip_address,omitempty"`
	UserAgent    string          `json:"user_agent,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
}

type AuditRepository interface {
	Record(ctx context.Context, log *AuditLog) error
	Search(ctx context.Context, filter Filter) ([]*AuditLog, int64, error)
}

type Filter struct {
	ActorID      string
	ResourceType string
	ResourceID   string
	Action       string
	FromDate     *time.Time
	ToDate       *time.Time
	Limit        int
	Offset       int
}
