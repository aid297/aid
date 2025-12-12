package rbac

type (
	PermissionAttributer interface {
		Register(permission *Permission)
	}

	AttrPermissionName     struct{ name string }
	AttrPermissionIdentity struct{ identity string }
)

func PermissionName(name string) PermissionAttributer         { return AttrPermissionName{name: name} }
func (my AttrPermissionName) Register(permission *Permission) { permission.Name = my.name }

func PermissionIdentity(identity string) PermissionAttributer {
	return AttrPermissionIdentity{identity: identity}
}
func (my AttrPermissionIdentity) Register(permission *Permission) { permission.Identity = my.identity }
