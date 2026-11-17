package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	favoritesManager "api-tests-template/internal/managers/favorites"
	base "api-tests-template/tests"
)

const nonexistentAdvertisementID = "00000000-0000-4000-8000-000000000000"

func TestFavoritesUnauthorized(t *testing.T) {
	base.SetupSuite()
	defer base.TearDownSuite()

	testCases := []struct {
		name            string
		token           string
		expectedMessage string
	}{
		{
			name:            "Без токена",
			token:           "",
			expectedMessage: "Authorization header required",
		},
		{
			name:            "Неверный токен",
			token:           "incorrect_token",
			expectedMessage: "Invalid or expired token",
		},
	}

	operations := []struct {
		name string
		call func(*testing.T, string) string
	}{
		{
			name: "POST favorites",
			call: func(t *testing.T, token string) string {
				return favoritesManager.AddFavorite(
					t, token, nonexistentAdvertisementID,
					http.StatusUnauthorized,
				)
			},
		},
		{
			name: "GET favorites",
			call: func(t *testing.T, token string) string {
				return favoritesManager.GetFavorites(
					t, token, http.StatusUnauthorized,
				)
			},
		},
		{
			name: "DELETE favorites",
			call: func(t *testing.T, token string) string {
				return favoritesManager.DeleteFavorite(
					t, token, nonexistentAdvertisementID,
					http.StatusUnauthorized,
				)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, operation := range operations {
				t.Run(operation.name, func(t *testing.T) {
					t.Log("Шаг: отправка запроса без корректной авторизации")

					body := operation.call(t, tc.token)

					require.True(t, gjson.Valid(body),
						"Ответ должен содержать корректный JSON")

					require.Equal(t, "unauthorized",
						gjson.Get(body, "error").String(),
						"Некорректный тип ошибки")

					require.Equal(t, tc.expectedMessage,
						gjson.Get(body, "message").String(),
						"Некорректное сообщение об ошибке")

					t.Log("Проверка 401 Unauthorized пройдена")
				})
			}
		})
	}
}
