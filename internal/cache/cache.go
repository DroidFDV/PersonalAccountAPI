package cache

import (
	"PersonalAccountAPI/internal/models"
	"PersonalAccountAPI/internal/repository"
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type CacheDecorator struct {
	repoProvider repository.RepoProvider
	mu           sync.RWMutex
	userMap      map[int]models.WrapUser
	ttl          time.Duration
}

func New(repository repository.RepoProvider, ttl time.Duration) *CacheDecorator {
	cache := &CacheDecorator{
		repoProvider: repository,
		mu:           sync.RWMutex{},
		userMap:      make(map[int]models.WrapUser),
		ttl:          ttl,
	}

	return cache
}

func (c *CacheDecorator) delete() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, wrapped := range c.userMap {
		if time.Now().After(wrapped.TTL) {
			delete(c.userMap, key)
		}
	}
}

func (c *CacheDecorator) RunCleaner(ctx context.Context, checkInterval time.Duration) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.delete()
		//WARNING: not work
		case <-ctx.Done():
			return
		}
	}
}

func (c *CacheDecorator) getUserMapValue(key int) (models.UserDTO, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	wrapped := c.userMap[key]
	if time.Now().After(wrapped.TTL) {
		return models.UserDTO{}, false
	}

	wrapped, exists := c.userMap[key]
	return wrapped.User, exists
}

func (c *CacheDecorator) setUserMapValue(key int, user models.UserDTO) {
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

func (c *CacheDecorator) GetIDByLogin(ctx context.Context, userRequest models.UserDTO) (models.UserDTO, error) {
	keyID := c.getKeyByLogPass(userRequest.Login, userRequest.Password)
	user, ok := c.getUserMapValue(keyID)
	if ok {
		return user, nil
	}

	userResponce, err := c.repoProvider.GetIDByLogin(ctx, userRequest)
	if err != nil {
		return models.UserDTO{}, errors.Wrap(err, "CacheDecorator.userProvider.GetIDByLogin:")
	}
	c.setUserMapValue(userResponce.ID, models.UserDTO{ID: userResponce.ID, Login: userRequest.Login, Password: userRequest.Password})
	return userResponce, nil
}

func (c *CacheDecorator) GetUserByID(ctx context.Context, userRequest models.UserDTO) (models.UserDTO, error) {
	user, ok := c.getUserMapValue(userRequest.ID)
	if ok {
		return user, nil
	}

	userResponce, err := c.repoProvider.GetUserByID(ctx, userRequest)
	if err != nil {
		return models.UserDTO{}, errors.Wrap(err, "CacheDecorator.userProvider.GetUserByID:")
	}
	c.setUserMapValue(userRequest.ID, models.UserDTO{ID: userRequest.ID, Login: userResponce.Login, Password: userRequest.Password})
	return userResponce, nil
}

func (c *CacheDecorator) AddingUser(ctx context.Context, userRequest models.UserDTO) error {
	if err := c.repoProvider.AddingUser(ctx, userRequest); err != nil {
		return errors.Wrap(err, "CacheDecorator.userProvider.AddingUser:")
	}
	c.setUserMapValue(userRequest.ID, userRequest)
	return nil
}

func (c *CacheDecorator) UpdateUser(ctx context.Context, userRequest models.UserDTO) error {
	if err := c.repoProvider.UpdateUser(ctx, userRequest); err != nil {
		return errors.Wrap(err, "CacheDecorator.userProvider.UpdateUser:")
	}

	user, exists := c.getUserMapValue(userRequest.ID)
	if exists {
		c.setUserMapValue(userRequest.ID, user)
	}

	return nil
}
