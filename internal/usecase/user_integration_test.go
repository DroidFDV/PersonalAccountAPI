//go:build integration
// +build integration

package usecase

import (
	"PersonalAccountAPI/internal/models"
	"PersonalAccountAPI/internal/repository"
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib" // for sql.Open with pgx
	"github.com/stretchr/testify/require"
)

const (
	testDBHost     = "localhost"
	testDBPort     = "5432"
	testDBUser     = "postgres"
	testDBPassword = "postgres"
	testDBName     = "test_personal_account"
)

func getTestDBDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		testDBUser,
		testDBPassword,
		testDBHost,
		testDBPort,
		testDBName,
	)
}

func setupTestDB(t *testing.T) *pgx.Conn {
	dsn := getTestDBDSN()
	conn, err := pgx.Connect(context.Background(), dsn)
	require.NoError(t, err)

	_, err = conn.Exec(context.Background(), "TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	t.Cleanup(func() {
		conn.Close(context.Background())
	})

	return conn
}

func TestUserUsecase_Integration(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New(db)
	usecase := New(repo)

	ctx := context.Background()

	// === Тест: AddingUser ===
	userReq := models.UserRequest{
		ID:       1,
		Login:    "testuser",
		Password: "securepass",
	}
	err := usecase.AddUser(ctx, userReq)
	require.NoError(t, err)

	// === Тест: GetIDByLogin ===
	got, err := usecase.GetIDByLogin(ctx, userReq)
	require.NoError(t, err)
	require.Equal(t, 1, got.ID)
	require.Equal(t, "testuser", got.Login)
	require.Equal(t, "securepass", got.Password)

	// === Тест: GetUserByID ===
	gotByID, err := usecase.GetUserByID(ctx, models.UserRequest{ID: got.ID})
	require.NoError(t, err)
	require.Equal(t, got.ID, gotByID.ID)
	require.Equal(t, "testuser", gotByID.Login)

	// === Тест: UpdateUser ===
	updateReq := models.UserRequest{
		ID:       got.ID,
		Login:    "updateduser",
		Password: "newpass",
	}
	err = usecase.UpdateUser(ctx, updateReq)
	require.NoError(t, err)

	updated, err := usecase.GetUserByID(ctx, models.UserRequest{ID: got.ID})
	require.NoError(t, err)
	require.Equal(t, "updateduser", updated.Login)
}
