package favorites

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	favoritesClient "api-tests-template/internal/client/http/favorites"
)

// checkFavoriteResponse проверяет статус и возвращает тело HTTP-ответа.
func checkFavoriteResponse(
	t *testing.T,
	response *http.Response,
	expectedStatus int,
) string {
	t.Helper()
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать ответ API избранного")

	require.Equalf(
		t,
		expectedStatus,
		response.StatusCode,
		"Ожидался HTTP %d, получен %d. Ответ: %s",
		expectedStatus,
		response.StatusCode,
		string(body),
	)

	return string(body)
}

// AddFavorite добавляет объявление в избранное и проверяет HTTP-статус.
func AddFavorite(
	t *testing.T,
	token, advertisementID string,
	expectedStatus int,
) string {
	t.Helper()

	response := favoritesClient.HttpPostFavorite(
		t, token, advertisementID,
	)

	return checkFavoriteResponse(t, response, expectedStatus)
}

// GetFavorites получает список избранного и проверяет HTTP-статус.
func GetFavorites(
	t *testing.T,
	token string,
	expectedStatus int,
) string {
	t.Helper()

	response := favoritesClient.HttpGetFavorites(t, token)

	return checkFavoriteResponse(t, response, expectedStatus)
}

// DeleteFavorite удаляет объявление из избранного и проверяет HTTP-статус.
func DeleteFavorite(
	t *testing.T,
	token, advertisementID string,
	expectedStatus int,
) string {
	t.Helper()

	response := favoritesClient.HttpDeleteFavorite(
		t, token, advertisementID,
	)

	return checkFavoriteResponse(t, response, expectedStatus)
}
