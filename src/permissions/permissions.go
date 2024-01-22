package permissions

type Entity string

const (
	EntityOrganization Entity = "organization"
	EntityDonation     Entity = "donation"
	EntityRoles        Entity = "roles"
)

type Action string

const (
	ActionRead   Action = "read"
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

func (entity Entity) String() string {
	return string(entity)
}

func (entity Entity) Capability(action Action) string {
	return string(entity) + ":" + string(action)
}
