package review

import "rent_service/internal/repository/interfaces/review"

type IProvider interface {
	GetReviewRepository() review.IRepository
}

