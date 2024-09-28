package provide

import "rent_service/internal/logic/services/interfaces/provide"

type IFactory interface {
	CreateProvisionService() provide.IService
}

type IRequestFactory interface {
	CreateProvisionRequestService() provide.IRequestService
}

type IRevokeFactory interface {
	CreateProvisionRevokeService() provide.IRevokeService
}

