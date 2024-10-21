package role

import "rent_service/internal/repository/interfaces/role"

type IAdministratorFactory interface {
	CreateRoleAdministratorRepository() role.IAdministratorRepository
}

type IRenterFactory interface {
	CreateRoleRenterRepository() role.IRenterRepository
}

type IStorekeeperFactory interface {
	CreateRoleStorekeeperRepository() role.IStorekeeperRepository
}

