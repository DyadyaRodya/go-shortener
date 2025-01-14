package dto

import usecasesdto "github.com/DyadyaRodya/go-shortener/internal/usecases/dto"

// StatsData Structure of response with total numbers of shortened urls and existing users
type StatsData struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}

// NewStatsDataResponse Creates new response *StatsData with given stats
func NewStatsDataResponse(data *usecasesdto.StatsResponse) *StatsData {
	return &StatsData{
		URLs:  data.URLs,
		Users: data.Users,
	}
}
