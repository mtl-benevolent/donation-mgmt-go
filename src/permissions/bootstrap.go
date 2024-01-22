package permissions

var permissionsService *PermissionsService

func Bootstrap() {
	permissionsService = NewPermissionsService()
}

func GetPermissionsService() *PermissionsService {
	if permissionsService == nil {
		panic("Permissions service not bootstrapped")
	}

	return permissionsService
}
