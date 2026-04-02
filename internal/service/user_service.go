package service

import (
	"context"
	"finance-backend/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

type UpdateUserInput struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	IsActive *bool  `json:"is_active"`
}

type UpdateProfileInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (s *UserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) List(ctx context.Context, page, pageSize int) ([]domain.User, int64, error) {
	return s.userRepo.List(ctx, page, pageSize)
}

func (s *UserService) Update(ctx context.Context, id string, input UpdateUserInput) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Role != "" {
		if input.Role != domain.RoleViewer && input.Role != domain.RoleAnalyst && input.Role != domain.RoleAdmin {
			return nil, domain.ErrForbidden
		}
		user.Role = input.Role
	}
	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, id string, input UpdateProfileInput) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashed)
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.userRepo.SoftDelete(ctx, id)
}
