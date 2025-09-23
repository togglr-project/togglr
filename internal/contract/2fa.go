package contract

import (
	"github.com/togglr-project/togglr/internal/domain"
)

type TwoFARateLimiter interface {
	Inc(userID domain.UserID) (attempts int, blocked bool)
	Reset(userID domain.UserID)
	IsBlocked(userID domain.UserID) bool
}
