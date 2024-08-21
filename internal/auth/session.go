package auth

import (
	"context"

	"github.com/Slyngshot-Team/packages/auth"
)

type Session struct {
	UserID string
	Tenant string
	Roles  map[string]struct{}
}

func SessionFromContext(ctx context.Context) (s Session, err error) {
	s.UserID, err = auth.GetUserID(ctx)
	if err == nil {
		s.Roles, err = auth.GetRoles(ctx)
	}
	return
}

func (s Session) HasRole(role string) (ok bool) {
	_, ok = s.Roles[role]
	return
}

func HasRole(ctx context.Context, role string) bool {
	r, err := auth.GetRoles(ctx)
	if err != nil {
		return false
	}
	_, ok := r[role]
	return ok
}
