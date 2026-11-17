package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

// HttpPostFavorite добавляет объявление в избранное пользователя.
func HttpPostFavorite(
	t *testing.T,
	token, advertisementID string,
) *http.Response {
	t.Helper()

	request := models.AddFavoriteRequest{
		AdvertisementID: advertisementID,
	}

	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(request).
		Expect(t).
		End().Response
}

// HttpGetFavorites возвращает первую страницу избранного.
func HttpGetFavorites(t *testing.T, token string) *http.Response {
	t.Helper()

	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Get(path.FavoritesPath).
		Query("page", "1").
		Query("page_size", "50").
		Expect(t).
		End().Response
}

// HttpDeleteFavorite удаляет объявление из избранного.
func HttpDeleteFavorite(
	t *testing.T,
	token, advertisementID string,
) *http.Response {
	t.Helper()

	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Delete(path.FavoritesPath + "/" + advertisementID).
		Expect(t).
		End().Response
}
