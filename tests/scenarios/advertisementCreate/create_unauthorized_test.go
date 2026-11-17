package advertisementCreate

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	advertisementsClient "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
)

func TestCreateAdvertisementUnauthorized(t *testing.T) {
	base.SetupSuite()
	defer base.TearDownSuite()

	testCases := []struct {
		name  string
		token string
	}{
		{
			name:  "Отсутствует токен",
			token: "",
		},
		{
			name:  "Неверный токен",
			token: "incorrect_token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Подготовка: формирование объявления")

			request := models.CreateRequest{
				Fields: map[string]string{
					"title":       "QA Test " + utils.RandomString(12),
					"description": testDescription,
					"price":       testPrice,
					"quantity":    testQuantity,
				},
				Photos: []models.PhotoFile{
					{
						Name: "photo.png",
						Data: makeTestPhoto(t),
					},
				},
			}

			t.Log("Шаг 1: создание объявления без корректной авторизации")

			response := advertisementsClient.HttpPostAdvertisement(
				t,
				tc.token,
				request,
			)
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			require.NoError(t, err, "Ошибка чтения ответа API")

			t.Logf("HTTP статус: %d", response.StatusCode)
			t.Logf("Ответ сервера: %s", string(body))

			require.Equal(t, http.StatusUnauthorized,
				response.StatusCode,
				"Ожидался HTTP 401, ответ: %s", string(body))

			require.True(t, gjson.ValidBytes(body),
				"Ответ с ошибкой должен содержать корректный JSON")

			require.NotEmpty(t, gjson.GetBytes(body, "error").String(),
				"В ответе отсутствует поле error")

			require.NotEmpty(t, gjson.GetBytes(body, "message").String(),
				"В ответе отсутствует сообщение об ошибке")

			t.Log("Проверка 401 Unauthorized успешно пройдена")
		})
	}
}
