// cache/mocks/user_provider_mock.go
package mocks

import (
	"PersonalAccountAPI/internal/models"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRepoProvider struct {
	mock.Mock
}

func (m *MockRepoProvider) GetIDByLogin(ctx context.Context, user models.UserDTO) (models.UserDTO, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(models.UserDTO), args.Error(1)
}

func (m *MockRepoProvider) GetUserByID(ctx context.Context, user models.UserDTO) (models.UserDTO, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(models.UserDTO), args.Error(1)
}

func (m *MockRepoProvider) AddingUser(ctx context.Context, user models.UserDTO) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepoProvider) UpdateUser(ctx context.Context, user models.UserDTO) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
