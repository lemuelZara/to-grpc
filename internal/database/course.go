package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type Course struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
	CategoryID  string
}

func NewCourse(db *sql.DB) Course {
	return Course{db: db}
}

func (c Course) Save(name string, description string, categoryID string) (Course, error) {
	id := uuid.New().String()

	_, err := c.db.Exec("INSERT INTO course(id, name, description, category_id) VALUES($1, $2, $3, $4)", id, name, description, categoryID)
	if err != nil {
		return Course{}, err
	}

	return Course{ID: id, Name: name, Description: description, CategoryID: categoryID}, nil
}

func (c Course) FindAll() ([]Course, error) {
	rows, err := c.db.Query("SELECT id, name, description, category_id FROM course")
	if err != nil {
		return []Course{}, err
	}
	defer rows.Close()

	courses := []Course{}
	for rows.Next() {
		var c Course
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CategoryID); err != nil {
			return []Course{}, err
		}

		courses = append(courses, c)
	}

	return courses, nil
}

func (c Course) FindByCategoryID(categoryID string) ([]Course, error) {
	rows, err := c.db.Query("SELECT id, name, description, category_id FROM course WHERE category_id = $1", categoryID)
	if err != nil {
		return []Course{}, err
	}
	defer rows.Close()

	courses := []Course{}
	for rows.Next() {
		var c Course
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CategoryID); err != nil {
			return []Course{}, err
		}

		courses = append(courses, c)
	}

	return courses, nil
}
