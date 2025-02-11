package cache

import (
	"PersonalAccountAPI/internal/models"
	"PersonalAccountAPI/internal/usecase"
	"context"
	"sync"
)

// редко будет использоваться (или вообще не булет),
// так что нееффективно
func getKeyByLogPass(m *map[int]models.UserRequest, login, password string) (int, bool) {
	for key, mapValue := range *m {
		if (mapValue.Login == login) && (mapValue.Password == password) {
			return key, true
		}
	}
	return 0, false
}

type CacheDecorator struct {
	mx           sync.RWMutex
	userProvider *usecase.UserUsecase
	userMap      map[int]models.UserRequest
	userLoginMap map[string]models.UserRequest
	// второстепенно ttl time.Duration
}

func New(user *usecase.UserUsecase) *CacheDecorator {
	return &CacheDecorator{
		mx:           sync.RWMutex{},
		userProvider: user,
		userMap:      make(map[int]models.UserRequest),
		userLoginMap: make(map[string]models.UserRequest),
	}
}

func (c *CacheDecorator) getUserMapValue(key int) (models.UserRequest, bool) {
	c.mx.RLock()
	user, exists := c.userMap[key]
	c.mx.RUnlock()
	return user, exists
}

func (c *CacheDecorator) setUserMapValue(key int, user models.UserRequest) {
	c.mx.Lock()
	c.userMap[key] = user
	c.mx.Unlock()
}

func (c *CacheDecorator) GetIDByLoginFromDB(ctx context.Context, login, password string) (int, error) {
	keyId, _ := getKeyByLogPass(&c.userMap, login, password)
	user, exists := c.getUserMapValue(keyId)
	if exists {
		return user.Id, nil
	} else {
		id, err := c.userProvider.GetIDByLoginFromDB(ctx, login, password)
		if err != nil {
			return id, err
		}
		c.setUserMapValue(id, models.UserRequest{Id: id, Login: login, Password: password})
		return id, err
	}
}

func (c *CacheDecorator) GetUserByIDFromDB(ctx context.Context, id int) (string, error) {
	keyId := id
	user, exists := c.getUserMapValue(keyId)
	if exists {
		return user.Login, nil
	} else {
		login, err := c.userProvider.GetUserByIDFromDB(ctx, id)
		if err != nil {
			return login, err
		}
		c.setUserMapValue(keyId, models.UserRequest{Id: id, Login: login})
		return login, err
	}
}

// Нужны ли эти методы? Нет, так как они изменяют базу
func (c *CacheDecorator) AddingUserToDB(ctx context.Context, id int, login, password string) error {
	c.setUserMapValue(id, models.UserRequest{Id: id, Login: login, Password: password})
	return c.userProvider.AddingUserToDB(ctx, id, login, password)
}

func (c *CacheDecorator) UpdateUserInDB(ctx context.Context, id int, login, password string) error {
	user, exists := c.getUserMapValue(id)
	if exists {
		c.setUserMapValue(id, user)
	}
	return c.userProvider.UpdateUserInDB(ctx, id, login, password)
}
