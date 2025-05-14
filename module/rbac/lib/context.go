package lib

import (
	"context"

	"github.com/euiko/webapp/module/rbac/lib/role"
)

type (
	contextKeyType struct{}
)

var (
	contextKey = contextKeyType{}
)

func ContextWithPermissions(ctx context.Context, manager role.PermissionManager) context.Context {
	return context.WithValue(ctx, contextKey, manager)
}

func PermissionsFromContext(ctx context.Context) role.PermissionManager {
	if v, ok := ctx.Value(contextKey).(role.PermissionManager); ok {
		return v
	}
	return nil
}
