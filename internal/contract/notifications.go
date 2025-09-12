package contract

import (
	"context"
)

type Emailer interface {
	SendResetPasswordEmail(ctx context.Context, email, token string) error
	Send2FACodeEmail(ctx context.Context, email, code, action string) error
}
