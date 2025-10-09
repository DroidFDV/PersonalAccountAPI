package cache

import (
	"PersonalAccountAPI/internal/models"
	"PersonalAccountAPI/internal/usecase"
	"context"
	"mime/multipart"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type CacheDecorator struct {
	userProvider usecase.UserProvider
	mu           sync.RWMutex
	userMap      map[int]models.WrapUser
	ttl          time.Duration
}

func New(user *usecase.UserUsecase, ttl time.Duration) *CacheDecorator {
	cache := &CacheDecorator{
		userProvider: user,
		mu:           sync.RWMutex{},
		userMap:      make(map[int]models.WrapUser),
		ttl:          ttl,
	}

	return cache
}

func (c *CacheDecorator) delete() {
	for key, wrapped := range c.userMap {
		if time.Now().After(wrapped.TTL) {
			c.mu.Lock()
			delete(c.userMap, key)
			c.mu.Unlock()
		}
	}
}

func (c *CacheDecorator) RunCleaner(checkInterval time.Duration) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.delete()
		//WARNING: not work
		case <-context.Background().Done():
			return
		}
	}
}

func (c *CacheDecorator) getUserMapValue(key int) (models.UserRequest, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	wrapped := c.userMap[key]
	if time.Now().After(wrapped.TTL) {
		return models.UserRequest{}, false
	}

	wrapped, exists := c.userMap[key]
	return wrapped.User, exists
}

func (c *CacheDecorator) setUserMapValue(key int, user models.UserRequest) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.userMap[key] = models.WrapUser{User: user, TTL: time.Now().Add(c.ttl)}
}

func (c *CacheDecorator) getKeyByLogPass(login, password string) int {
	for key, wrapped := range c.userMap {
		if (wrapped.User.Login == login) && (wrapped.User.Password == password) {
			return key
		}
	}
	return 0
}

func (c *CacheDecorator) GetCacheSize() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.userMap)
}

func (c *CacheDecorator) GetIDByLogin(ctx context.Context, userRequest models.UserRequest) (models.UserResponse, error) {
	keyID := c.getKeyByLogPass(userRequest.Login, userRequest.Password)
	user, ok := c.getUserMapValue(keyID)
	if ok {
		return user.ToResponce(), nil
	}

	userResponce, err := c.userProvider.GetIDByLogin(ctx, userRequest)
	if err != nil {
		return models.UserResponse{}, errors.Wrap(err, "CacheDecorator.userProvider.GetIDByLogin:")
	}
	c.setUserMapValue(userResponce.ID, models.UserRequest{ID: userResponce.ID, Login: userRequest.Login, Password: userRequest.Password})
	return userResponce, nil
}

func (c *CacheDecorator) GetUserByID(ctx context.Context, userRequest models.UserRequest) (models.UserResponse, error) {
	user, ok := c.getUserMapValue(userRequest.ID)
	if ok {
		return user.ToResponce(), nil
	}

	userResponce, err := c.userProvider.GetUserByID(ctx, userRequest)
	if err != nil {
		return models.UserResponse{}, errors.Wrap(err, "CacheDecorator.userProvider.GetUserByID:")
	}
	c.setUserMapValue(userRequest.ID, models.UserRequest{ID: userRequest.ID, Login: userResponce.Login, Password: userRequest.Password})
	return userResponce, nil
}

func (c *CacheDecorator) AddingUser(ctx context.Context, userRequest models.UserRequest) error {
	if err := c.userProvider.AddingUser(ctx, userRequest); err != nil {
		return errors.Wrap(err, "CacheDecorator.userProvider.AddingUser:")
	}
	c.setUserMapValue(userRequest.ID, userRequest)
	return nil
}

func (c *CacheDecorator) UpdateUser(ctx context.Context, userRequest models.UserRequest) error {
	if err := c.userProvider.UpdateUser(ctx, userRequest); err != nil {
		return errors.Wrap(err, "CacheDecorator.userProvider.UpdateUser:")
	}

	user, exists := c.getUserMapValue(userRequest.ID)
	if exists {
		c.setUserMapValue(userRequest.ID, user)
	}

	return nil
}

func (c *CacheDecorator) UploadFile(ctx context.Context, file *multipart.FileHeader) func(context.Context) error {
	return c.userProvider.UploadFile(ctx, file)
}
