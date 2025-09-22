//nolint:revive // it's ok here
package context

import (
	"context"
	"net/http"

	"github.com/rom8726/etoggle/internal/domain"
)

type contextKey string

const (
	ctxKeyProjectID  contextKey = "project_id"
	ctxKeyUserID     contextKey = "user_id"
	ctxKeyIsSuper    contextKey = "is_superuser"
	ctxKeyRawRequest contextKey = "raw_request"
	ctxKeyRequestID  contextKey = "request_id"
	ctxKeyActorID    contextKey = "actor"
	ctxKeyUsername   contextKey = "username"
)

func WithProjectID(ctx context.Context, id domain.ProjectID) context.Context {
	return context.WithValue(ctx, ctxKeyProjectID, id)
}

func ProjectID(ctx context.Context) domain.ProjectID {
	return ctx.Value(ctxKeyProjectID).(domain.ProjectID)
}

func WithUserID(ctx context.Context, userID domain.UserID) context.Context {
	return context.WithValue(ctx, ctxKeyUserID, userID)
}

func WithIsSuper(ctx context.Context, isSuper bool) context.Context {
	return context.WithValue(ctx, ctxKeyIsSuper, isSuper)
}

func IsSuper(ctx context.Context) bool {
	v, ok := ctx.Value(ctxKeyIsSuper).(bool)

	return ok && v
}

func UserID(ctx context.Context) domain.UserID {
	id, ok := ctx.Value(ctxKeyUserID).(domain.UserID)
	if !ok {
		return 0
	}

	return id
}

func WithRawRequest(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, ctxKeyRawRequest, req)
}

func RawRequest(ctx context.Context) *http.Request {
	return ctx.Value(ctxKeyRawRequest).(*http.Request)
}

func WithRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, ctxKeyRequestID, reqID)
}

func RequestID(ctx context.Context) string {
	v, ok := ctx.Value(ctxKeyRequestID).(string)
	if ok {
		return v
	}

	return ""
}

func WithActor(ctx context.Context, actor domain.AuditActor) context.Context {
	return context.WithValue(ctx, ctxKeyActorID, actor)
}

func Actor(ctx context.Context) domain.AuditActor {
	v, ok := ctx.Value(ctxKeyActorID).(domain.AuditActor)
	if ok {
		return v
	}

	return domain.AuditActorSystem
}

func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, ctxKeyUsername, username)
}

func Username(ctx context.Context) string {
	v, ok := ctx.Value(ctxKeyUsername).(string)
	if ok {
		return v
	}

	return ""
}
