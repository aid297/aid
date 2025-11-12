package rbac

type (
	RoleAttributer interface {
		Register(role *Role)
	}

	AttrRoleName struct{ name string }
)

func RoleName(name string) RoleAttributer {
	return AttrRoleName{name: name}
}

func (my AttrRoleName) Register(role *Role) {
	role.Name = my.name
}
