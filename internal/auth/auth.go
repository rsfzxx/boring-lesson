package auth

import "context"

const (
	RoleAdmin  = "Admin"
	RoleLeader = "Leader"
	RoleMember = "Member"
)

type User struct {
	ID       string
	RoleID   string
	RoleName string
	GroupID  string
}

type contextKey struct{}

func WithUser(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, contextKey{}, u)
}

func UserFrom(ctx context.Context) (User, bool) {
	u, ok := ctx.Value(contextKey{}).(User)
	return u, ok
}