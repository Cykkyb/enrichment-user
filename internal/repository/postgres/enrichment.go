package repository

import (
	"enrichment-user-info/internal/entity"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strings"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type EnrichmentPostgres struct {
	db *sqlx.DB
}

func NewEnrichmentPostgres(db *sqlx.DB) *EnrichmentPostgres {
	return &EnrichmentPostgres{
		db: db,
	}
}

func (r *EnrichmentPostgres) CreateUser(user *entity.User) error {
	query := fmt.Sprintf(`INSERT INTO %s (name, surname, patronymic, age, gender, nationality) VALUES ($1, $2, $3, $4, $5, $6) returning id`, usersTable)

	var id int
	err := r.db.QueryRow(query, user.Name, user.Surname, user.Patronymic, user.Age, user.Gender, user.Nationality).Scan(&id)
	if err != nil {
		return err
	}

	user.Id = id

	return nil
}

func (r *EnrichmentPostgres) GetUser(id int) (*entity.User, error) {
	query := fmt.Sprintf(`SELECT id, name, surname, patronymic, age, gender, nationality FROM %s WHERE id = $1`, usersTable)

	var user entity.User

	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *EnrichmentPostgres) UpdateUser(user *entity.User) error {
	query := fmt.Sprintf(`UPDATE %s SET name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nationality = $6 WHERE id = $7`, usersTable)

	_, err := r.db.Exec(query, user.Name, user.Surname, user.Patronymic, user.Age, user.Gender, user.Nationality, user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *EnrichmentPostgres) DeleteUser(id int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, usersTable)

	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

type filterField struct {
	value     string
	condition string
}

var filterFields = []filterField{
	{"Id", "id = $%d"},
	{"Name", "name = $%d"},
	{"Surname", "surname = $%d"},
	{"Age", "age = $%d"},
	{"MinAge", "age >= $%d"},
	{"MaxAge", "age <= $%d"},
	{"Nationality", "nationality = $%d"},
}

func (r *EnrichmentPostgres) constructFilterConditions(filter *entity.UserFilter) ([]string, []interface{}) {
	var conditions []string
	var args []interface{}

	counter := 1

	for _, field := range filterFields {
		fieldValue := reflect.ValueOf(filter).Elem().FieldByName(field.value)

		if fieldValue.IsValid() {
			value := fieldValue.Interface()

			switch v := value.(type) {
			case string:
				if v != "" {
					conditions = append(conditions, fmt.Sprintf(field.condition, counter))
					args = append(args, value)
					counter++
				}
			case int:
				if v > 0 {
					conditions = append(conditions, fmt.Sprintf(field.condition, counter))
					args = append(args, value)
					counter++
				}
			}
		}
	}

	return conditions, args
}

func (r *EnrichmentPostgres) GetUsersWithFilter(filter *entity.UserFilter) ([]entity.User, error) {
	var users []entity.User

	conditions, args := r.constructFilterConditions(filter)

	query := "SELECT * FROM users"

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" OFFSET %d", filter.Offset)

	if filter.Limit > 0 {
		if filter.Limit > 100 {
			filter.Limit = 100
		}
		query += fmt.Sprintf(" LIMIT %d", filter.Limit)
	}

	err := r.db.Select(&users, query, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}
