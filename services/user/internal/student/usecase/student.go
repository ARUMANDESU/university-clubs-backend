package usecase

import (
	"context"
	"github.com/ARUMANDESU/university-clubs-backend/services/user/internal/domain"
)

type StudentUsecase interface {
	// SignUp TODO: доделать
	SignUp(ctx context.Context)
	GetByEmail(ctx context.Context, email string) (*domain.Student, error)
	GetByID(ctx context.Context, id int64) (*domain.Student, error)
	Update(ctx context.Context, student *domain.Student) error
	Delete(ctx context.Context, id int64) error
}
