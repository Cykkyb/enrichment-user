package service

import (
	"database/sql"
	"encoding/json"
	"enrichment-user-info/internal/entity"
	"enrichment-user-info/internal/repository/postgres"
	"log/slog"
	"net/http"
)

type EnrichmentService struct {
	repository repository.Enrichment
	log        *slog.Logger
}

func NewEnrichmentService(repo repository.Enrichment, log *slog.Logger) *EnrichmentService {
	return &EnrichmentService{
		repository: repo,
		log:        log,
	}
}

func (s *EnrichmentService) CreateUser(user *entity.User) error {
	if err := s.EnrichUser(user); err != nil {
		return err
	}

	if err := user.Validate(); err != nil {
		return err
	}

	return s.repository.CreateUser(user)
}

func (s *EnrichmentService) GetUser(id int) (*entity.User, error) {
	return s.repository.GetUser(id)
}

func (s *EnrichmentService) GetUsersWithFilter(filter *entity.UserFilter) ([]entity.User, error) {
	return s.repository.GetUsersWithFilter(filter)
}

func (s *EnrichmentService) UpdateUser(user *entity.User) error {
	existingUser, err := s.GetUser(user.Id)
	if err != nil {
		return err
	}

	existingUser.Name = user.Name
	existingUser.Surname = user.Surname
	existingUser.Patronymic = user.Patronymic
	existingUser.Age = user.Age
	existingUser.Gender = user.Gender
	existingUser.Nationality = user.Nationality

	if err = existingUser.Validate(); err != nil {
		s.log.Error(err.Error())
		return err
	}

	if err = s.repository.UpdateUser(existingUser); err != nil {
		return err
	}

	return nil
}

func (s *EnrichmentService) DeleteUser(id int) error {
	_, err := s.GetUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.ErrUserNotFound
		}

		return err
	}

	return s.repository.DeleteUser(id)
}

func (s *EnrichmentService) EnrichUser(user *entity.User) error {
	data, err := getEnrichmentData(user.Name)
	if err != nil {
		return err
	}

	user.Age = data.Age
	user.Gender = data.Gender
	user.Nationality = data.Nationality

	return nil
}

type EnrichmentData struct {
	Age         int           `json:"age"`
	Gender      string        `json:"gender"`
	Nationality string        `json:"nationality"`
	CountryInfo []CountryData `json:"country"`
}

type CountryData struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

func getEnrichmentData(name string) (EnrichmentData, error) {
	ageData, err := fetchEnrichmentData("https://api.agify.io/?name=" + name)
	if err != nil {
		return EnrichmentData{}, err
	}

	genderData, err := fetchEnrichmentData("https://api.genderize.io/?name=" + name)
	if err != nil {
		return EnrichmentData{}, err
	}

	nationalityData, err := fetchEnrichmentData("https://api.nationalize.io/?name=" + name)
	if err != nil {
		return EnrichmentData{}, err
	}

	return EnrichmentData{
		Age:         ageData.Age,
		Gender:      genderData.Gender,
		Nationality: nationalityData.Nationality,
	}, nil
}

func fetchEnrichmentData(url string) (EnrichmentData, error) {
	response, err := http.Get(url)
	if err != nil {
		return EnrichmentData{}, err
	}

	defer response.Body.Close()

	var data EnrichmentData
	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return EnrichmentData{}, err
	}

	if len(data.CountryInfo) > 0 {
		data.Nationality = data.CountryInfo[0].CountryID
	}

	return data, nil
}
