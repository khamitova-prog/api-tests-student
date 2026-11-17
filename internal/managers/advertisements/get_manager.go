package advertisements

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	advertisementsClient "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
)

// decodeOKResponse проверяет статус 200 и преобразует JSON-ответ в модель.
func decodeOKResponse(
	t *testing.T,
	response *http.Response,
	result any,
) {
	t.Helper()
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать ответ API")

	require.Equalf(
		t,
		http.StatusOK,
		response.StatusCode,
		"Ожидался HTTP 200, получен %d. Ответ: %s",
		response.StatusCode,
		string(body),
	)

	require.NoError(
		t,
		json.Unmarshal(body, result),
		"Не удалось разобрать JSON-ответ",
	)
}

// GetAdvertisementByID получает созданное объявление и проверяет статус ответа.
func GetAdvertisementByID(
	t *testing.T,
	token, id string,
) models.CreateResponse {
	t.Helper()

	response := advertisementsClient.HttpGetAdvertisementByID(
		t, token, id,
	)

	var result models.CreateResponse
	decodeOKResponse(t, response, &result)

	return result
}

// GetAdvertisementPhotos получает список фотографий объявления.
func GetAdvertisementPhotos(
	t *testing.T,
	token, id string,
) []models.CreatedPhoto {
	t.Helper()

	response := advertisementsClient.HttpGetAdvertisementPhotos(
		t, token, id,
	)

	var photos []models.CreatedPhoto
	decodeOKResponse(t, response, &photos)

	return photos
}
