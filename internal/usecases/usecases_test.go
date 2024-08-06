package usecases

import (
	"errors"
	usecasesMocks "github.com/DyadyaRodya/go-shortener/internal/usecases/mocks"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type usecasesSuite struct {
	suite.Suite

	urlStorage  *usecasesMocks.URLStorage
	idGenerator *usecasesMocks.IDGenerator

	usecases *Usecases
}

func (u *usecasesSuite) SetupTest() {
	t := u.T()

	u.urlStorage = usecasesMocks.NewURLStorage(t)
	u.idGenerator = usecasesMocks.NewIDGenerator(t)

	u.usecases = NewUsecases(u.urlStorage, u.idGenerator)
}

func TestRunUsecasesSuite(t *testing.T) {
	suite.Run(t, new(usecasesSuite))
}

func (u *usecasesSuite) TestNewUsecases() {
	expected := &Usecases{u.urlStorage, u.idGenerator}
	if got := NewUsecases(u.urlStorage, u.idGenerator); !reflect.DeepEqual(got, expected) {
		u.Errorf(errors.New("NewUsecases error"), "NewUsecases() = %v, want %v", got, expected)
	}
}
