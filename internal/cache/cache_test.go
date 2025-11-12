package cache

import (
	"PersonalAccountAPI/internal/cache/mocks"
	"PersonalAccountAPI/internal/models"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCacheDecorator_GetIDByLogin_CacheHit(t *testing.T) {
	mockProvider := &mocks.MockRepoProvider{}
	mockProvider.On("GetIDByLogin", mock.Anything, mock.Anything).Return(&models.UserDTO{ID: 123, Login: "test"}, nil)

	cache := New(mockProvider, 10*time.Second)

	req := models.UserDTO{Login: "test", Password: "pass"}

	// Первый вызов — должен вызвать провайдера
	_, err := cache.GetIDByLogin(context.Background(), &req)
	require.NoError(t, err)
	mockProvider.AssertNumberOfCalls(t, "GetIDByLogin", 1)

	// Второй вызов — должен взять из кэша
	_, err = cache.GetIDByLogin(context.Background(), &req)
	require.NoError(t, err)
	mockProvider.AssertNumberOfCalls(t, "GetIDByLogin", 1) // всё ещё 1 раз
}

func TestCacheDecorator_TTL_Expiry(t *testing.T) {
	mockProvider := &mocks.MockRepoProvider{}
	mockProvider.On("GetIDByLogin", mock.Anything, mock.Anything).Return(&models.UserDTO{ID: 123, Login: "test"}, nil).Twice()

	cache := New(mockProvider, 100*time.Millisecond)

	req := models.UserDTO{Login: "test", Password: "pass"}

	_, err := cache.GetIDByLogin(context.Background(), &req)
	require.NoError(t, err)
	mockProvider.AssertNumberOfCalls(t, "GetIDByLogin", 1)

	time.Sleep(150 * time.Millisecond)

	_, err = cache.GetIDByLogin(context.Background(), &req)
	require.NoError(t, err)
	mockProvider.AssertNumberOfCalls(t, "GetIDByLogin", 2)
}

func TestCacheDecorator_RunCleaner(t *testing.T) {
	mockProvider := &mocks.MockRepoProvider{}
	cache := New(mockProvider, 50*time.Millisecond)
	cache.set(1, &models.UserDTO{ID: 1})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go cache.RunCleaner(ctx, 10*time.Millisecond)

	time.Sleep(100 * time.Millisecond)
	cancel()

	assert.Equal(t, 0, cache.GetCacheSize())
}
