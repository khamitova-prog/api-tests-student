
package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

func ContainsAdvertisementByTitle(
	t *testing.T,
	items []models.Advertisement,
	title string,
) {
	t.Helper()

	found := false

	for _, item := range items {
		if item.Title == title {
			found = true
			break
		}
	}

	require.True(t, found, "Объявление с названием %q не найдено", title)
}
