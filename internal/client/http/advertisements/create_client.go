package advertisements

import (
	"net/http"
	"testing"
	"testing/fstest"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

// HttpPostAdvertisement отправляет объявление как multipart/form-data.
func HttpPostAdvertisement(
	t *testing.T,
	token string,
	data models.CreateRequest,
) *http.Response {
	t.Helper()

	files := fstest.MapFS{}
	fileNames := make([]string, 0, len(data.Photos))

	for _, photo := range data.Photos {
		files[photo.Name] = &fstest.MapFile{Data: photo.Data}
		fileNames = append(fileNames, photo.Name)
	}

	request := apiRunner.GetRunner().
		Auth(token).
		Create().
		UseFS(files).
		Post(path.AdvertisementPath)

	for name, value := range data.Fields {
		request = request.MultipartFormData(name, value)
	}

	if len(fileNames) > 0 {
		request = request.MultipartFile("photos", fileNames...)
	}

	return request.Expect(t).End().Response
}

// HttpDeleteAdvertisement удаляет объявление по ID.

// HttpDeleteAdvertisement удаляет объявление по ID.
func HttpDeleteAdvertisement(
	t *testing.T,
	token string,
	id string,
) *http.Response {
	t.Helper()

	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}
