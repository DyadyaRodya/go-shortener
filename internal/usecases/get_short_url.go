package usecases

import "github.com/DyadyaRodya/go-shortener/internal/domain/entity"

func (u *Usecases) GetShortURL(ID string) (*entity.ShortURL, error) {
	shortURL, err := u.urlStorage.GetURLByID(ID)
	if err != nil {
		return nil, err
	}
	return shortURL, nil
}
