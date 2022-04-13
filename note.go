package main

import (
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Note struct {
	ID          int       `json:"id,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreateAt    time.Time `json:"create_at,omitempty"`
	UpdateAt    time.Time `json:"update_at,omitempty"`
}

func (n *Note) Create() error {
	db := GetConnection()

	q := `INSERT INTO notes (title, description, update_at) VALUES (?, ?, ?)`

	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	r, err := stmt.Exec(n.Title, n.Description, time.Now())
	if err != nil {
		return err
	}
	if i, err := r.RowsAffected(); err != nil || i != 1 {
		return errors.New("ERROR: Se esperaba una fila afectada")
	}

	return nil
}

func (n *Note) GetAll() ([]Note, error) {
	db := GetConnection()
	q := `SELECT id, title, description, create_at, update_at FROM notes`

	rows, err := db.Query(q)

	if err != nil {
		return []Note{}, err
	}
	defer rows.Close()

	notes := []Note{}

	for rows.Next() {
		rows.Scan(&n.ID, &n.Title, &n.Description, &n.CreateAt, &n.UpdateAt)
		notes = append(notes, *n)
	}
	return notes, nil
}

func (n *Note) Update() error {
	db := GetConnection()

	q := `UPDATE notes SET title = ?, description = ?, update_at = ? WHERE id = ?`

	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	r, err := stmt.Exec(n.Title, n.Description, time.Now(), n.ID)
	if err != nil {
		return err
	}
	if i, err := r.RowsAffected(); err != nil || i != 1 {
		return errors.New("ERROR: Se esperaba una fila afectada")
	}
	return nil
}

func (n *Note) Delete(id int) error {
	db := GetConnection()
	q := `DELETE FROM notes WHERE id = ?`

	stmt, err := db.Prepare(q)

	if err != nil {
		return err
	}
	defer stmt.Close()

	r, err := stmt.Exec(id)

	if err != nil {
		return err
	}

	if i, err := r.RowsAffected(); err != nil || i != 1 {
		return errors.New("ERROR: Se esperaba una fila afectada")
	}
	return nil
}
