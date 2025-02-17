package cache

import (
	"PersonalAccountAPI/internal/models"
	"PersonalAccountAPI/internal/usecase"
	"context"
	"mime/multipart"
	"sync"

	"github.com/pkg/errors"
)

type CacheDecorator struct {
	userProvider *usecase.UserUsecase

	mx      sync.RWMutex
	userMap map[int]models.UserRequest
	// userLoginMap map[string]models.UserRequest
	// ttl time.Duration
}

func New(user *usecase.UserUsecase) *CacheDecorator {
	return &CacheDecorator{
		userProvider: user,
		mx:           sync.RWMutex{},
		userMap:      make(map[int]models.UserRequest),
		// userLoginMap: make(map[string]models.UserRequest),
		// ttl
	}
}

func (c *CacheDecorator) getUserMapValue(key int) (models.UserRequest, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	user, exists := c.userMap[key]
	return user, exists
}

func (c *CacheDecorator) setUserMapValue(key int, user models.UserRequest) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.userMap[key] = user
}

// не эффективно
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

func (c *CacheDecorator) GetIDByLoginFromDB(ctx context.Context, login, password string) (int, error) {
	keyID := c.getKeyByLogPass(login, password)
	user, ok := c.getUserMapValue(keyID)
	if ok {
		return user.ID, nil
	}

	id, err := c.userProvider.GetIDByLoginFromDB(ctx, login, password)
	if err != nil {
		return id, errors.Wrap(err, "CacheDecorator.userProvider.GetIDByLoginFromDB:")
	}
	c.setUserMapValue(id, models.UserRequest{ID: id, Login: login, Password: password})
	return id, errors.Wrap(err, "CacheDecorator.GetIDByLoginFromDB:")
}

func (c *CacheDecorator) GetUserByIDFromDB(ctx context.Context, id int) (string, error) {
	user, ok := c.getUserMapValue(id)
	if ok {
		return user.Login, nil
	}

	login, err := c.userProvider.GetUserByIDFromDB(ctx, id)
	if err != nil {
		return login, errors.Wrap(err, "CacheDecorator.userProvider.GetUserByIDFromDB:")
	}
	c.setUserMapValue(id, models.UserRequest{ID: id, Login: login})
	return login, errors.Wrap(err, "CacheDecorator.userProvider.GetUserByIDFromDB:")
}

func (c *CacheDecorator) AddingUserToDB(ctx context.Context, id int, login, password string) error {
	err := c.userProvider.AddingUserToDB(ctx, id, login, password)
	if err != nil {
		return errors.Wrap(err, "CacheDecorator.userProvider.AddingUserToDB:")
	}

	c.setUserMapValue(id, models.UserRequest{ID: id, Login: login, Password: password})
	return errors.Wrap(err, "CacheDecorator.userProvider.AddingUserToDB:")
}

func (c *CacheDecorator) UpdateUserInDB(ctx context.Context, id int, login, password string) error {
	err := c.userProvider.UpdateUserInDB(ctx, id, login, password)
	if err != nil {
		return errors.Wrap(err, "CacheDecorator.userProvider.UpdateUserInDB:")
	}

	user, exists := c.getUserMapValue(id)
	if exists {
		c.setUserMapValue(id, user)
	}
	return errors.Wrap(err, "CacheDecorator.userProvider.UpdateUserInDB:")
}

func (c *CacheDecorator) SetFile(file *multipart.FileHeader) {
	c.userProvider.SetFile(file)
}
func (c *CacheDecorator) UploadFile(ctx context.Context) error {
	return errors.Wrap(c.userProvider.UploadFile(ctx), "CacheDecorator.userProvider.UploadFile:")
}
