package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES ($1, $2,
            CURRENT_TIMESTAMP AT TIME ZONE 'UTC',
            (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') + ($3 * INTERVAL '1 day'))
    RETURNING id`

	var id int
	err := m.DB.QueryRow(stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	// Write the SQL statement we want to execute. Again, I've split it over two
	// lines for readability.
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') AND id = $1`
	//инициализируем структуру snippet
	var s Snippet
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		// If the query returns no rows, then row.Scan() will return a
		// sql.ErrNoRows error. We use the errors.Is() function check for that
		// error specifically, and return our own ErrNoRecord error
		// instead (we'll create this in a moment).
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	// If everything went OK, then return the filled Snippet struct.
	return s, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
    ORDER BY id DESC
    LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Expires, &s.Created)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)

	}

	if err = rows.Err(); err != nil {
		return nil, err // Возвращаем ошибку, а не битый результат
	}
	return snippets, err

}

func (m *SnippetModel) Delete(id int) error {
	stmt := `DELETE FROM snippets WHERE id = $1`

	result, err := m.DB.Exec(stmt, id)

	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("entries not found")
	}
	return nil
}
func (m *SnippetModel) Update(id int, title, content string, expires int) error {

	stmt := `
		UPDATE snippets
		SET title = $1,
		    content = $2,
		    expires = (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') + ($3 * INTERVAL '1 day')
		WHERE id = $4`

	result, err := m.DB.Exec(stmt, title, content, expires, id)

	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("entries not update")
	}
	return nil
}
