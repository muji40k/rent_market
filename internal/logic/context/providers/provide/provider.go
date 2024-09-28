package provide

import "rent_service/internal/logic/services/interfaces/provide"

type IProvider interface {
	GetProvisionService() provide.IService
}

type IRequestProvider interface {
	GetProvisionRequestService() provide.IRequestService
}

type IRevokeProvider interface {
	GetProvisionRevokeService() provide.IRevokeService
}

