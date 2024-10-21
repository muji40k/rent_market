package review

import "rent_service/internal/repository/interfaces/review"

type IFactory interface {
	CreateReviewRepository() review.IRepository
}

