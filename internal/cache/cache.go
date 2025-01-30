package cache

import (
	"PersonalAccountAPI/internal/usecase"
	"context"

	"github.com/jackc/pgx/v5"
)

type CacheDecorator struct {
	userProvider usecase.UserUsecase
	cMap         map[usecase.User]usecase.User
	// второстепенно ttl time.Duration
}

func NewCache(conn *pgx.Conn) *CacheDecorator {
	return &CacheDecorator{
		userProvider: *usecase.NewUser(conn),
		// Возможно для cache разумно использовать ограниченный размер
		cMap: make(map[usecase.User]usecase.User),
	}
}

// Возможно стоит использовать xxhash

func (cache *CacheDecorator) GetIDByLoginFromDB(ctx context.Context, login, password string) (int, error) {
	key := usecase.User{Login: login, Password: password}
	user, exists := cache.cMap[key]
	if exists {
		return user.Id, nil
	} else {
		id, err := cache.userProvider.GetIDByLoginFromDB(ctx, login, password)
		cache.cMap[key] = usecase.User{Id: id, Login: login, Password: password}
		return id, err
	}
}

func (cache *CacheDecorator) GetUserByIDFromDB(ctx context.Context, id int) (string, error) {
	key := usecase.User{Id: id}
	user, exists := cache.cMap[key]
	if exists {
		return user.Login, nil
	} else {
		login, err := cache.userProvider.GetUserByIDFromDB(ctx, id)
		cache.cMap[key] = usecase.User{Id: id, Login: login}
		return login, err
	}
}

// Нужны ли эти методы? Нет, так как они изменяют базу
// func (user *CacheDecorator) AddingUserToDB(ctx context.Context, id int, login, password string) error {

// }

// func (user *CacheDecorator) UpdateUserInDB(ctx context.Context, id int, login, password string) error {

// }
