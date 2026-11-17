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

// CreateAdvertisement создаёт объявление и проверяет ответ 201.
func CreateAdvertisement(
	t *testing.T,
	token string,
	request models.CreateRequest,
) models.CreateResponse {
	t.Helper()

	response := advertisementsClient.HttpPostAdvertisement(t, token, request)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать ответ создания объявления")

	require.Equalf(
		t,
		http.StatusCreated,
		response.StatusCode,
		"Ожидался 201 Created, получен %d. Ответ: %s",
		response.StatusCode,
		string(body),
	)

	var result models.CreateResponse
	require.NoError(t, json.Unmarshal(body, &result),
		"Не удалось разобрать ответ создания объявления")

	return result
}

// DeleteAdvertisement удаляет тестовое объявление при очистке.
func DeleteAdvertisement(t *testing.T, token, id string) {
	t.Helper()

	response := advertisementsClient.HttpDeleteAdvertisement(t, token, id)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать ответ удаления")

	require.Equalf(
		t,
		http.StatusNoContent,
		response.StatusCode,
		"Ошибка удаления объявления %s: HTTP %d, ответ: %s",
		id,
		response.StatusCode,
		string(body),
	)
}

// CreateAdvertisementExpectBadRequest проверяет отклонение некорректного объявления.
func CreateAdvertisementExpectBadRequest(
	t *testing.T,
	token string,
	request models.CreateRequest,
) models.APIError {
	t.Helper()

	response := advertisementsClient.HttpPostAdvertisement(
		t,
		token,
		request,
	)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать ответ API")

	// Если объявление неожиданно создано, удаляем его после теста.
	if response.StatusCode == http.StatusCreated {
		var created models.CreateResponse
		if json.Unmarshal(body, &created) == nil && created.ID != "" {
			t.Cleanup(func() {
				DeleteAdvertisement(t, token, created.ID)
			})
		}
	}

	require.Equalf(
		t,
		http.StatusBadRequest,
		response.StatusCode,
		"Ожидался HTTP 400, получен %d. Ответ сервера: %s",
		response.StatusCode,
		string(body),
	)

	var apiError models.APIError
	require.NoError(
		t,
		json.Unmarshal(body, &apiError),
		"Не удалось разобрать JSON ошибки",
	)

	require.NotEmpty(t, apiError.Error,
		"В ответе отсутствует поле error")

	require.NotEmpty(t, apiError.Message,
		"В ответе отсутствует поле message")

	return apiError
}
