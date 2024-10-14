package role

import "rent_service/internal/repository/interfaces/role"

type IAdministratorProvider interface {
	GetAdministratorRepository() role.IAdministratorRepository
}

type IRenterProvider interface {
	GetRenterRepository() role.IRenterRepository
}

type IStorekeeperProvider interface {
	GetStorekeeperRepository() role.IStorekeeperRepository
}

