package permissions

type Entity string

const (
	Organization Entity = "organization"
	Sandbox      Entity = "sandbox"
	Donation     Entity = "donation"
	Roles        Entity = "roles"
)

type Action string

const (
	Read   Action = "read"
	Create Action = "create"
	Update Action = "update"
	Delete Action = "delete"
)

func (entity Entity) String() string {
	return string(entity)
}

func (entity Entity) Capability(action Action) string {
	return string(entity) + ":" + string(action)
}
