
package advertisementSearch

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	advertisementsManager "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"

	base "api-tests-template/tests"
)

const (
	existingAdTitle = "Внимание"
	searchBeginning = "Вни"
	searchMiddle    = "има"
	searchEnd       = "ние"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	require.NotEmpty(s.T(), existingAdTitle)
	require.NotEmpty(s.T(), searchBeginning)
	require.NotEmpty(s.T(), searchMiddle)
	require.NotEmpty(s.T(), searchEnd)

	s.loginData = auth.Login(
		s.T(),
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) searchAndCheck(searchText string) {
	result := advertisementsManager.GetAdvertisementsList(
		s.T(),
		s.loginData.Token,
		searchText,
	)

	assertions.ContainsAdvertisementByTitle(
		s.T(),
		result.Items,
		existingAdTitle,
	)
}

func (s *TestSuite) TestSearchByTitleParts() {
	testCases := []struct {
		name  string
		query string
	}{
		{"Поиск по началу названия", searchBeginning},
		{"Поиск по середине названия", searchMiddle},
		{"Поиск по концу названия", searchEnd},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.searchAndCheck(tc.query)
		})
	}
}

func (s *TestSuite) TestSearchByDifferentCases() {
	testCases := []struct {
		name  string
		query string
	}{
		{"Строчные буквы", strings.ToLower(existingAdTitle)},
		{"Заглавные буквы", strings.ToUpper(existingAdTitle)},
		{"Смешанный регистр", "ВнИмАнИе"},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.searchAndCheck(tc.query)
		})
	}
}
