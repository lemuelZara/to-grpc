package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type Category struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
}

func NewCategory(db *sql.DB) Category {
	return Category{db: db}
}

func (c Category) Save(name string, description string) (Category, error) {
	id := uuid.New().String()

	_, err := c.db.Exec("INSERT INTO category(id, name, description) VALUES($1, $2, $3)", id, name, description)
	if err != nil {
		return Category{}, err
	}

	return Category{ID: id, Name: name, Description: description}, nil
}

func (c Category) FindAll() ([]Category, error) {
	rows, err := c.db.Query("SELECT id, name, description FROM category")
	if err != nil {
		return []Category{}, err
	}
	defer rows.Close()

	categories := []Category{}
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return []Category{}, err
		}

		categories = append(categories, c)
	}

	return categories, nil
}

func (c Category) FindByCourseID(courseID string) (Category, error) {
	row := c.db.QueryRow(
		"SELECT ca.id, ca.name, ca.description FROM category ca INNER JOIN course co ON ca.id = co.category_id WHERE co.id = $1",
		courseID,
	)

	var ca Category
	if err := row.Scan(&ca.ID, &ca.Name, &ca.Description); err != nil {
		return Category{}, err
	}

	return ca, nil
}
