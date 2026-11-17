package favorites

import (
	"bytes"
	"image"
	"image/png"
	"net/http"
	"os"
	"testing"

	advertisements "api-tests-template/internal/managers/advertisements"
	advertisementModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	favoritesManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

const missingFavoriteAdID = "00000000-0000-4000-8000-000000000000"

// makeFavoriteTestPhoto создаёт PNG для тестового объявления.
func makeFavoriteTestPhoto(t *testing.T) []byte {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, 64, 64))

	var buffer bytes.Buffer
	require.NoError(t, png.Encode(&buffer, img))

	return buffer.Bytes()
}

// checkFavoriteAPIError проверяет структуру ответа API с ошибкой.
func checkFavoriteAPIError(t *testing.T, body string) {
	t.Helper()

	require.True(t, gjson.Valid(body),
		"Ответ сервера должен содержать JSON")

	require.NotEmpty(t, gjson.Get(body, "error").String(),
		"Поле error не должно быть пустым")

	require.NotEmpty(t, gjson.Get(body, "message").String(),
		"Поле message не должно быть пустым")
}

func TestAddNonexistentAdvertisementToFavorites(t *testing.T) {
	base.SetupSuite()
	defer base.TearDownSuite()

	t.Log("Подготовка: авторизация пользователя")

	login := auth.Login(
		t,
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)

	t.Log("Шаг 1: добавление несуществующего объявления")

	body := favoritesManager.AddFavorite(
		t,
		login.Token,
		missingFavoriteAdID,
		http.StatusNotFound,
	)

	checkFavoriteAPIError(t, body)

	t.Logf("Ответ API: %s", body)
	t.Log("Проверка 404 Not Found пройдена")
}

func TestAddOwnAdvertisementToFavorites(t *testing.T) {
	base.SetupSuite()
	defer base.TearDownSuite()

	t.Log("Подготовка: авторизация пользователя")

	login := auth.Login(
		t,
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)

	t.Log("Шаг 1: создание собственного объявления")

	request := advertisementModels.CreateRequest{
		Fields: map[string]string{
			"title":       "QA Own Favorite " + utils.RandomString(12),
			"description": "Проверка запрета добавления собственного объявления",
			"price":       "4999",
			"quantity":    "2",
		},
		Photos: []advertisementModels.PhotoFile{
			{
				Name: "photo.png",
				Data: makeFavoriteTestPhoto(t),
			},
		},
	}

	created := advertisements.CreateAdvertisement(
		t,
		login.Token,
		request,
	)

	require.NotEmpty(t, created.ID)

	// Удаляем созданное объявление после завершения теста.
	t.Cleanup(func() {
		t.Log("Очистка: удаление собственного объявления")

		advertisements.DeleteAdvertisement(
			t,
			login.Token,
			created.ID,
		)
	})

	t.Logf("Создано объявление ID=%s", created.ID)

	t.Log("Шаг 2: попытка добавить своё объявление в избранное")

	body := favoritesManager.AddFavorite(
		t,
		login.Token,
		created.ID,
		http.StatusForbidden,
	)

	checkFavoriteAPIError(t, body)

	t.Logf("Ответ API: %s", body)
	t.Log("Проверка 403 Forbidden пройдена")
}
