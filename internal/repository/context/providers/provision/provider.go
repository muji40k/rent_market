package provision

import "rent_service/internal/repository/interfaces/provision"

type IProvider interface {
	GetProvisionRepository() provision.IRepository
}

type IRequestProvider interface {
	GetProvisionRequestRepository() provision.IRequestRepository
}

type IRevokeProvider interface {
	GetRevokeProvisionRepository() provision.IRevokeRepository
}

