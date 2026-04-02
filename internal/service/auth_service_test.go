package service

import (
	"context"
	"finance-backend/internal/domain"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// mockUserRepo implements domain.UserRepository for testing
type mockUserRepo struct {
	users map[string]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*domain.User)}
}

func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) error {
	for _, u := range m.users {
		if u.Email == user.Email {
			return domain.ErrDuplicateEmail
		}
	}
	user.ID = "test-id-" + user.Email
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserRepo) List(ctx context.Context, page, pageSize int) ([]domain.User, int64, error) {
	var users []domain.User
	for _, u := range m.users {
		users = append(users, *u)
	}
	return users, int64(len(users)), nil
}

func (m *mockUserRepo) Update(ctx context.Context, user *domain.User) error {
	if _, ok := m.users[user.ID]; !ok {
		return domain.ErrNotFound
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) SoftDelete(ctx context.Context, id string) error {
	if _, ok := m.users[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.users, id)
	return nil
}

func TestRegister_Success(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo, "test-secret", time.Hour)

	resp, err := svc.Register(context.Background(), RegisterInput{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected token to be set")
	}
	if resp.User.Email != "test@example.com" {
		t.Fatalf("expected email test@example.com, got %s", resp.User.Email)
	}
	if resp.User.Role != domain.RoleViewer {
		t.Fatalf("expected role viewer, got %s", resp.User.Role)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo, "test-secret", time.Hour)

	_, _ = svc.Register(context.Background(), RegisterInput{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	})

	_, err := svc.Register(context.Background(), RegisterInput{
		Email:    "test@example.com",
		Password: "password456",
		Name:     "Another User",
	})

	if err != domain.ErrDuplicateEmail {
		t.Fatalf("expected ErrDuplicateEmail, got %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo, "test-secret", time.Hour)

	_, _ = svc.Register(context.Background(), RegisterInput{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	})

	resp, err := svc.Login(context.Background(), LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected token to be set")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo, "test-secret", time.Hour)

	_, _ = svc.Register(context.Background(), RegisterInput{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	})

	_, err := svc.Login(context.Background(), LoginInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})

	if err != domain.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo, "test-secret", time.Hour)

	_, err := svc.Login(context.Background(), LoginInput{
		Email:    "nobody@example.com",
		Password: "password123",
	})

	if err != domain.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_InactiveUser(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewAuthService(repo, "test-secret", time.Hour)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), 12)
	repo.users["inactive-user"] = &domain.User{
		ID:       "inactive-user",
		Email:    "inactive@example.com",
		Password: string(hashed),
		Name:     "Inactive",
		Role:     domain.RoleViewer,
		IsActive: false,
	}

	_, err := svc.Login(context.Background(), LoginInput{
		Email:    "inactive@example.com",
		Password: "password123",
	})

	if err != domain.ErrUserInactive {
		t.Fatalf("expected ErrUserInactive, got %v", err)
	}
}
