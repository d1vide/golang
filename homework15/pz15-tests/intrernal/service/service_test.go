package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRepo struct {
	users map[string]User
	calls []string
}

func (r *testRepo) ByEmail(email string) (User, error) {
	r.calls = append(r.calls, email)

	u, ok := r.users[email]
	if !ok {
		return User{}, ErrNotFound
	}
	return u, nil
}

func TestService_FindIDByEmail(t *testing.T) {
	users := map[string]User{
		"user@example.com":  {ID: 1, Email: "user@example.com"},
		"admin@example.com": {ID: 2, Email: "admin@example.com"},
		"test@example.com":  {ID: 3, Email: "test@example.com"},
	}

	repo := &testRepo{
		users: users,
		calls: make([]string, 0),
	}
	service := New(repo)

	t.Run("пользователь найден", func(t *testing.T) {
		email := "user@example.com"
		expectedID := int64(1)

		actualID, err := service.FindIDByEmail(email)

		require.NoError(t, err, "Не должно быть ошибки при поиске существующего пользователя")
		assert.Equal(t, expectedID, actualID, "ID должен соответствовать ожидаемому")

		assert.Contains(t, repo.calls, email, "Репозиторий должен был быть вызван с этим email")
	})

	t.Run("другой пользователь найден", func(t *testing.T) {
		email := "admin@example.com"
		expectedID := int64(2)

		actualID, err := service.FindIDByEmail(email)

		require.NoError(t, err)
		assert.Equal(t, expectedID, actualID)
	})

	t.Run("пользователь не найден", func(t *testing.T) {
		email := "nonexistent@example.com"

		actualID, err := service.FindIDByEmail(email)

		require.Error(t, err, "Должна быть ошибка при поиске несуществующего пользователя")
		assert.ErrorIs(t, err, ErrNotFound, "Ошибка должна быть ErrNotFound")
		assert.Equal(t, int64(0), actualID, "При ошибке должен возвращаться 0")

		assert.Contains(t, repo.calls, email, "Репозиторий должен был быть вызван даже для несуществующего пользователя")
	})

	t.Run("пустой email", func(t *testing.T) {
		email := ""

		actualID, err := service.FindIDByEmail(email)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Equal(t, int64(0), actualID)
	})

	t.Run("регистр email не важен для репозитория", func(t *testing.T) {
		email := "User@Example.com"

		actualID, err := service.FindIDByEmail(email)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Equal(t, int64(0), actualID)
	})
}
