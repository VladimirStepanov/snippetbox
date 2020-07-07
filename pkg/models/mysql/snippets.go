package mysql

import (
	"database/sql"
	"fmt"

	"githib.com/VladimirStepanov/snippetbox/pkg/models"
	"github.com/go-sql-driver/mysql"
)

//SnippetStore struct for working with snippets table
type SnippetStore struct {
	DB *sql.DB
}

// Insert snippet into database
func (s *SnippetStore) Insert(title, content string, expire int, isPublic bool, ownerID int64) (int64, error) {
	res, err := s.DB.Exec(
		`INSERT into snippets (title, content, create_date, expiration_date, is_public, owner_id) 
		VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), ?, ?)`,
		title,
		content,
		expire,
		isPublic,
		ownerID,
	)

	if err != nil {
		if me, ok := err.(*mysql.MySQLError); ok {
			if me.Number == 1452 {
				return 0, models.ErrUnknownOwnerID
			}
		}
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return id, nil
}

//Delete from snippets
func (s *SnippetStore) Delete(snippetID, userID int64) error {
	res, err := s.DB.Exec("DELETE from snippets WHERE id=? and owner_id=?", snippetID, userID)

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if count == 0 {
		return models.ErrNoRecord
	}

	return nil

}

//Get specific snippet
func (s *SnippetStore) Get(snippetID int64) (*models.Snippet, error) {
	res := &models.Snippet{}
	row := s.DB.QueryRow(
		`SELECT id, title, content, create_date, expiration_date, is_public, owner_id from snippets 
		WHERE id=? AND expiration_date > CURDATE()`,
		snippetID,
	)

	err := row.Scan(&res.ID, &res.Title, &res.Content, &res.Created, &res.Expires, &res.IsPublic, &res.OwnerID)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

//Update snippet
func (s *SnippetStore) Update(snippet *models.Snippet, ownerID int64) error {
	res, err := s.DB.Exec(
		"update snippets set title = ?, content = ?, is_public = ? where id = ? and owner_id = ?",
		snippet.Title,
		snippet.Content,
		snippet.IsPublic,
		snippet.ID,
		ownerID,
	)

	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ra == 0 {
		return models.ErrNoRecord
	}

	return nil
}

func (s *SnippetStore) getSnippets(rows *sql.Rows) ([]*models.Snippet, error) {

	snippets := []*models.Snippet{}

	for rows.Next() {
		res := &models.Snippet{}

		err := rows.Scan(&res.ID, &res.Title, &res.Content, &res.Created, &res.Expires, &res.IsPublic, &res.OwnerID)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, res)
	}

	return snippets, nil
}

//LatestAll return latest snippets sorted by create_date
func (s *SnippetStore) LatestAll(ownerID int64, count, page int) ([]*models.Snippet, error) {
	var rows *sql.Rows
	var err error
	var limit string

	if page == 1 {
		limit = fmt.Sprintf("LIMIT %d", count)
	} else {
		limit = fmt.Sprintf("LIMIT %d, %d", count*page-count, count)
	}
	if ownerID == -1 {
		rows, err = s.DB.Query(
			`SELECT id, title, content, create_date, expiration_date, is_public, owner_id from snippets 
			WHERE expiration_date > CURDATE() AND is_public = 1 ORDER BY create_date DESC ` + limit,
		)
	} else {
		rows, err = s.DB.Query(
			`SELECT id, title, content, create_date, expiration_date, is_public, owner_id from snippets
			WHERE expiration_date > CURDATE() AND owner_id = ? ORDER BY create_date DESC `+limit, ownerID,
		)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res, err := s.getSnippets(rows)

	return res, err
}
