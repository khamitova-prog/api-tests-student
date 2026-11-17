package advertisementCreate

import (
	"os"
	"strings"
	"testing"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"

	"github.com/stretchr/testify/require"
)

const (
	negativePrice = "-100"
	invalidPrice  = "not-a-number"
	longTitleSize = 2048
)

// validCreateRequest создаёт корректный запрос для негативных тестов.
func validCreateRequest(t *testing.T) models.CreateRequest {
	t.Helper()

	return models.CreateRequest{
		Fields: map[string]string{
			"title":       "QA Negative " + utils.RandomString(12),
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
}

func TestCreateAdvertisementInvalidParameters(t *testing.T) {
	base.SetupSuite()
	defer base.TearDownSuite()

	t.Log("Подготовка: авторизация пользователя")

	login := auth.Login(
		t,
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)

	testCases := []struct {
		name            string
		modify          func(*models.CreateRequest)
		expectedMessage string
	}{
		{
			name: "Отсутствует title",
			modify: func(r *models.CreateRequest) {
				delete(r.Fields, "title")
			},
			expectedMessage: "title is required",
		},
		{
			name: "Пустой title",
			modify: func(r *models.CreateRequest) {
				r.Fields["title"] = ""
			},
			expectedMessage: "title is required",
		},
		{
			name: "Отсутствует description",
			modify: func(r *models.CreateRequest) {
				delete(r.Fields, "description")
			},
			expectedMessage: "description is required",
		},
		{
			name: "Пустой description",
			modify: func(r *models.CreateRequest) {
				r.Fields["description"] = ""
			},
			expectedMessage: "description is required",
		},
		{
			name: "Отрицательная цена",
			modify: func(r *models.CreateRequest) {
				r.Fields["price"] = negativePrice
			},
			expectedMessage: "price must be between 0 and 1000000 inclusive",
		},
		{
			name: "Некорректный формат цены",
			modify: func(r *models.CreateRequest) {
				r.Fields["price"] = invalidPrice
			},
			expectedMessage: "cannot convert string value not-a-number to float",
		},
		{
			name: "Слишком длинное название",
			modify: func(r *models.CreateRequest) {
				r.Fields["title"] = strings.Repeat("A", longTitleSize)
			},
			expectedMessage: "title must not exceed 50 characters provided 2048",
		},
		{
			name: "Отсутствуют фотографии",
			modify: func(r *models.CreateRequest) {
				r.Photos = nil
			},
			expectedMessage: "at least 1 photo is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Подготовка: создание тестовых данных")

			request := validCreateRequest(t)
			tc.modify(&request)

			t.Log("Шаг 1: отправка некорректного объявления")

			apiError := advertisements.CreateAdvertisementExpectBadRequest(
				t,
				login.Token,
				request,
			)
			require.Equal(t, "invalid_request", apiError.Error,
				"Некорректный тип ошибки")

			require.Equal(t, tc.expectedMessage, apiError.Message,
				"Некорректное сообщение об ошибке")

			t.Logf(
				"Получена ошибка: error=%q, message=%q, details=%q",
				apiError.Error,
				apiError.Message,
				apiError.Details,
			)

			t.Log("Проверка HTTP 400 и структуры ошибки пройдена")
		})
	}
}
