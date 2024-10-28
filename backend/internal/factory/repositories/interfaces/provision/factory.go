package provision

import "rent_service/internal/repository/interfaces/provision"

type IFactory interface {
	CreateProvisionRepository() provision.IRepository
}

type IRequestFactory interface {
	CreateProvisionRequestRepository() provision.IRequestRepository
}

type IRevokeFactory interface {
	CreateProvisionRevokeRepository() provision.IRevokeRepository
}

