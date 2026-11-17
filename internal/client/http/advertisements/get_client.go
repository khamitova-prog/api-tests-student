package advertisements

import (
	"fmt"
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpGetAdvertisementByID получает объявление по ID через query-параметр.
func HttpGetAdvertisementByID(
	t *testing.T,
	token, id string,
) *http.Response {
	t.Helper()

	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementPhotos получает фотографии объявления по ID.
func HttpGetAdvertisementPhotos(
	t *testing.T,
	token, id string,
) *http.Response {
	t.Helper()

	photosPath := fmt.Sprintf(
		"%s/%s/photos",
		path.AdvertisementsPath,
		id,
	)

	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Get(photosPath).
		Expect(t).
		End().Response
}
