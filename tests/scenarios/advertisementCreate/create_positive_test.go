package advertisementCreate

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
)

const (
	testDescription = "Тестовое объявление для проверки API создания"
	testPrice       = "4999"
	testQuantity    = "2"
)

// makeTestPhoto создаёт корректное PNG-изображение в памяти.
func makeTestPhoto(t *testing.T) []byte {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, 64, 64))

	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			img.SetRGBA(x, y, color.RGBA{
				R: uint8(x * 4),
				G: uint8(y * 4),
				B: 120,
				A: 255,
			})
		}
	}

	var buffer bytes.Buffer
	require.NoError(t, png.Encode(&buffer, img))

	return buffer.Bytes()
}

func TestCreateAdvertisementPositive(t *testing.T) {
	base.SetupSuite()
	defer base.TearDownSuite()

	t.Log("Подготовка: авторизация пользователя")
	login := auth.Login(
		t,
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)

	title := "QA Test " + utils.RandomString(12)

	request := models.CreateRequest{
		Fields: map[string]string{
			"title":       title,
			"description": testDescription,
			"price":       testPrice,
			"quantity":    testQuantity,
		},
		Photos: []models.PhotoFile{
			{Name: "photo.png", Data: makeTestPhoto(t)},
		},
	}

	t.Log("Шаг 1: создание объявления со всеми параметрами")

	created := advertisements.CreateAdvertisement(
		t,
		login.Token,
		request,
	)

	require.NotEmpty(t, created.ID, "ID объявления пуст")

	t.Cleanup(func() {
		t.Log("Очистка: удаление созданного объявления")
		advertisements.DeleteAdvertisement(t, login.Token, created.ID)
	})

	t.Log("Шаг 2: проверка параметров созданного объявления")

	require.Equal(t, title, created.Title, "Название не совпадает")
	require.Equal(t, testDescription, created.Description,
		"Описание не совпадает")
	require.Equal(t, float64(4999), created.Price,
		"Цена не совпадает")
	require.Equal(t, 2, created.Quantity,
		"Количество не совпадает")
	require.Equal(t, login.User.Id, created.UserID,
		"ID владельца не совпадает")

	require.Len(t, created.Photos, 1,
		"Ожидалась одна загруженная фотография")
	require.NotEmpty(t, created.Photos[0].ID,
		"ID фотографии пуст")
	require.NotEmpty(t, created.Photos[0].URL,
		"URL фотографии пуст")

	t.Log("Шаг 3: получение объявления через GET /advertisement")

	actual := advertisements.GetAdvertisementByID(
		t,
		login.Token,
		created.ID,
	)

	require.Equal(t, created.ID, actual.ID,
		"ID полученного объявления не совпадает")

	require.Equal(t, title, actual.Title,
		"Название объявления не совпадает")

	require.Equal(t, testDescription, actual.Description,
		"Описание объявления не совпадает")

	require.Equal(t, float64(4999), actual.Price,
		"Цена объявления не совпадает")

	require.Equal(t, 2, actual.Quantity,
		"Количество товара не совпадает")

	require.Equal(t, login.User.Id, actual.UserID,
		"Владелец объявления не совпадает")

	require.Len(t, actual.Photos, 1,
		"У объявления должна быть одна фотография")

	require.Equal(t, created.Photos[0].ID, actual.Photos[0].ID,
		"ID фотографии не совпадает")

	t.Log("GET по ID: все параметры совпадают")

	t.Log("Шаг 4: поиск созданного объявления по полному названию")

	searchResult := advertisements.GetAdvertisementsList(
		t,
		login.Token,
		title,
	)

	found := false

	for _, item := range searchResult.Items {
		if item.ID == created.ID {
			require.Equal(t, title, item.Title,
				"Название найденного объявления не совпадает")
			found = true
			break
		}
	}

	require.Truef(
		t,
		found,
		"Созданное объявление ID=%s с названием %q не найдено в поиске",
		created.ID,
		title,
	)

	t.Log("Поиск: созданное объявление найдено")

	t.Log("Шаг 5: проверка фотографии через GET /advertisements/{id}/photos")

	photos := advertisements.GetAdvertisementPhotos(
		t,
		login.Token,
		created.ID,
	)

	require.Len(t, photos, 1,
		"Сервер должен вернуть одну загруженную фотографию")

	require.Equal(t, created.Photos[0].ID, photos[0].ID,
		"ID фотографии не совпадает с ответом создания")

	require.Equal(t, created.Photos[0].URL, photos[0].URL,
		"URL фотографии не совпадает с ответом создания")

	require.NotEmpty(t, photos[0].URL,
		"URL загруженной фотографии пустой")

	t.Log("Фотография успешно найдена на сервере")

	t.Logf("Объявление %s успешно создано и проверено", created.ID)
}
