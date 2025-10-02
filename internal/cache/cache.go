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
	ttl          time.Duration

	mu      sync.RWMutex
	userMap map[int]models.UserRequest
	ttls    map[int]time.Time
	// userLoginMap map[string]models.UserRequest
}

func New(user *usecase.UserUsecase, ttl time.Duration) *CacheDecorator {
	cache := &CacheDecorator{
		userProvider: user,
		ttl:          ttl,
		mu:           sync.RWMutex{},
		userMap:      make(map[int]models.UserRequest),
		ttls:         make(map[int]time.Time),
		// userLoginMap: make(map[string]models.UserRequest),
	}

	return cache
}

func (c *CacheDecorator) delete() {
	for key, ttl := range c.ttls {
		if time.Now().After(ttl) {
			c.mu.Lock()
			delete(c.userMap, key)
			delete(c.ttls, key)
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
			c.mu.RLock()
			defer c.mu.RUnlock()
			c.delete()
		case <-context.Background().Done():
			return
		}
	}
}

func (c *CacheDecorator) getUserMapValue(key int) (models.UserRequest, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	outOfTTL := c.ttls[key]
	if time.Now().After(outOfTTL) {
		return models.UserRequest{}, false
	}

	user, exists := c.userMap[key]
	return user, exists
}

func (c *CacheDecorator) setUserMapValue(key int, user models.UserRequest) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.userMap[key] = user
	c.ttls[key] = time.Now().Add(c.ttl)
}

// todo: не эффективно, добавить карту пользователей по логину
func (c *CacheDecorator) getKeyByLogPass(login, password string) int {
	for key, mapValue := range c.userMap {
		if (mapValue.Login == login) && (mapValue.Password == password) {
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

func (c *CacheDecorator) GetIDByLogin(ctx context.Context, userRequest models.UserRequest) (int, error) {
	keyID := c.getKeyByLogPass(userRequest.Login, userRequest.Password)
	user, ok := c.getUserMapValue(keyID)
	if ok {
		return user.ID, nil
	}

	id, err := c.userProvider.GetIDByLogin(ctx, userRequest)
	if err != nil {
		return 0, errors.Wrap(err, "CacheDecorator.userProvider.GetIDByLogin:")
	}
	c.setUserMapValue(id, models.UserRequest{ID: id, Login: userRequest.Login, Password: userRequest.Password})
	return id, nil
}

func (c *CacheDecorator) GetUserByID(ctx context.Context, userRequest models.UserRequest) (string, error) {
	user, ok := c.getUserMapValue(userRequest.ID)
	if ok {
		return user.Login, nil
	}

	login, err := c.userProvider.GetUserByID(ctx, userRequest)
	if err != nil {
		return "", errors.Wrap(err, "CacheDecorator.userProvider.GetUserByID:")
	}
	c.setUserMapValue(userRequest.ID, models.UserRequest{ID: userRequest.ID, Login: login, Password: userRequest.Password})
	return login, nil
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
