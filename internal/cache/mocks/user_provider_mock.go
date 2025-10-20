// cache/mocks/user_provider_mock.go
package mocks

import (
	"PersonalAccountAPI/internal/models"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockUserProvider struct {
	mock.Mock
}

func (m *MockUserProvider) GetIDByLogin(ctx context.Context, user models.UserRequest) (models.UserResponse, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(models.UserResponse), args.Error(1)
}

func (m *MockUserProvider) GetUserByID(ctx context.Context, user models.UserRequest) (models.UserResponse, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(models.UserResponse), args.Error(1)
}

func (m *MockUserProvider) AddingUser(ctx context.Context, user models.UserRequest) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserProvider) UpdateUser(ctx context.Context, user models.UserRequest) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
