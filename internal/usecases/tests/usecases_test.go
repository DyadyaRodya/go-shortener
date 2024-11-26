package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"

	usecases2 "github.com/DyadyaRodya/go-shortener/internal/usecases"
	usecasesMocks "github.com/DyadyaRodya/go-shortener/internal/usecases/mocks"
)

type usecasesSuite struct {
	suite.Suite

	urlStorage  *usecasesMocks.URLStorage
	idGenerator *usecasesMocks.IDGenerator

	usecases *usecases2.Usecases
}

func (u *usecasesSuite) SetupTest() {
	t := u.T()

	u.urlStorage = usecasesMocks.NewURLStorage(t)
	u.idGenerator = usecasesMocks.NewIDGenerator(t)

	u.usecases = usecases2.NewUsecases(u.urlStorage, u.idGenerator)
}

func TestRunUsecasesSuite(t *testing.T) {
	suite.Run(t, new(usecasesSuite))
}
