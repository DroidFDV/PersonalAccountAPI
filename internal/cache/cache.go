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

	mx      sync.RWMutex
	userMap map[int]models.UserRequest
	ttls    map[int]time.Time
	// userLoginMap map[string]models.UserRequest
}

func New(user *usecase.UserUsecase, ttl time.Duration, checkInterval time.Duration) *CacheDecorator {
	cache := &CacheDecorator{
		userProvider: user,
		ttl:          ttl,
		mx:           sync.RWMutex{},
		userMap:      make(map[int]models.UserRequest),
		ttls:         make(map[int]time.Time),
		// userLoginMap: make(map[string]models.UserRequest),
	}

	go cache.Cleaner(checkInterval)

	return cache
}

func (c *CacheDecorator) Cleaner(checkInterval time.Duration) {
	for {
		time.Sleep(checkInterval)
		c.mx.RLock()

		for key, ttl := range c.ttls {
			if time.Now().After(ttl) {
				c.mx.Lock()
				delete(c.userMap, key)
				delete(c.ttls, key)
				c.mx.Unlock()
			}
		}

		c.mx.RUnlock()
	}
}

func (c *CacheDecorator) getUserMapValue(key int) (models.UserRequest, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	outOfTTL := c.ttls[key]
	if time.Now().After(outOfTTL) {
		return models.UserRequest{}, false
	}

	user, exists := c.userMap[key]
	return user, exists
}

func (c *CacheDecorator) setUserMapValue(key int, user models.UserRequest) {
	c.mx.Lock()
	defer c.mx.Unlock()

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
	c.mx.RLock()
	defer c.mx.RUnlock()

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
		return id, errors.Wrap(err, "CacheDecorator.userProvider.GetIDByLogin:")
	}
	c.setUserMapValue(id, models.UserRequest{ID: id, Login: userRequest.Login, Password: userRequest.Password})
	return id, errors.Wrap(err, "CacheDecorator.GetIDByLogin:")
}

func (c *CacheDecorator) GetUserByID(ctx context.Context, userRequest models.UserRequest) (string, error) {
	user, ok := c.getUserMapValue(userRequest.ID)
	if ok {
		return user.Login, nil
	}

	login, err := c.userProvider.GetUserByID(ctx, userRequest)
	if err != nil {
		return login, errors.Wrap(err, "CacheDecorator.userProvider.GetUserByID:")
	}
	c.setUserMapValue(userRequest.ID, models.UserRequest{ID: userRequest.ID, Login: login, Password: userRequest.Password})
	return login, errors.Wrap(err, "CacheDecorator.userProvider.GetUserByID:")
}

func (c *CacheDecorator) AddingUser(ctx context.Context, userRequest models.UserRequest) error {
	err := c.userProvider.AddingUser(ctx, userRequest)
	if err != nil {
		return errors.Wrap(err, "CacheDecorator.userProvider.AddingUser:")
	}

	c.setUserMapValue(userRequest.ID, userRequest)
	return errors.Wrap(err, "CacheDecorator.userProvider.AddingUser:")
}

func (c *CacheDecorator) UpdateUser(ctx context.Context, userRequest models.UserRequest) error {
	err := c.userProvider.UpdateUser(ctx, userRequest)
	if err != nil {
		return errors.Wrap(err, "CacheDecorator.userProvider.UpdateUser:")
	}

	user, exists := c.getUserMapValue(userRequest.ID)
	if exists {
		c.setUserMapValue(userRequest.ID, user)
	}

	return errors.Wrap(err, "CacheDecorator.userProvider.UpdateUser:")
}

func (c *CacheDecorator) UploadFile(ctx context.Context, file *multipart.FileHeader) func(context.Context) error {
	return c.userProvider.UploadFile(ctx, file)
}
