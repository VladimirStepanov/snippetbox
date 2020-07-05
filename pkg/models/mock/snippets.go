package mock

import (
	"math/rand"
	"sort"
	"time"

	"githib.com/VladimirStepanov/snippetbox/pkg/models"
)

func remove(slice []*models.Snippet, s int) []*models.Snippet {
	return append(slice[:s], slice[s+1:]...)
}

func getRandSnippetID(slice []*models.Snippet) int64 {
	var found bool
	for {
		id := rand.Int63()
		for _, v := range slice {
			if v.ID == id {
				found = true
				break
			}
		}

		if !found {
			return id
		}
	}
}

//SnippetStore mock for snippets
type SnippetStore struct {
	DB       []*models.Snippet
	UsersMap map[int64]*models.User
}

//Insert snippet to map
func (s *SnippetStore) Insert(title, content string, expire int, isPublic bool, ownerID int64) (int64, error) {
	if _, ok := s.UsersMap[ownerID]; !ok {
		return 0, models.ErrUnknownOwnerID
	}

	id := getRandSnippetID(s.DB)

	s.DB = append(s.DB, &models.Snippet{
		ID:       id,
		Title:    title,
		Content:  content,
		Created:  time.Now(),
		Expires:  time.Now().AddDate(0, 0, expire),
		OwnerID:  ownerID,
		IsPublic: isPublic,
	})

	return id, nil
}

//Get specific snippet
func (s *SnippetStore) Get(snippetID, userID int64) (*models.Snippet, error) {
	for _, value := range s.DB {
		if value.ID == snippetID && value.OwnerID == userID && value.Expires.After(time.Now()) {
			return value, nil
		}
	}

	return nil, models.ErrNoRecord
}

//Delete from snippets
func (s *SnippetStore) Delete(snippetID, userID int64) error {
	for i, value := range s.DB {
		if value.ID == snippetID && value.OwnerID == userID && value.Expires.After(time.Now()) {
			remove(s.DB, i)
			return nil
		}
	}

	return models.ErrNoRecord
}

//Update from snippets
func (s *SnippetStore) Update(snippet *models.Snippet, ownerID int64) error {
	for _, value := range s.DB {
		if value.ID == snippet.ID && value.OwnerID == ownerID && value.Expires.After(time.Now()) {
			value.Title = snippet.Title
			value.Content = snippet.Content
			value.IsPublic = snippet.IsPublic
			return nil
		}
	}

	return models.ErrNoRecord
}

//LatestAll return latest snippets sorted by create_date
func (s *SnippetStore) LatestAll(ownerID int64, count, page int) ([]*models.Snippet, error) {
	start := page*count - count + 1
	if page == 1 {
		start = 0
	}
	sort.SliceStable(s.DB, func(i, j int) bool {
		return s.DB[i].Expires.Before(s.DB[j].Expires)
	})
	res := []*models.Snippet{}

	for _, val := range s.DB[start:] {
		if val.Expires.After(time.Now()) {
			if ownerID == -1 && val.IsPublic {
				res = append(res, val)
			} else if ownerID != -1 && ownerID == val.OwnerID {
				res = append(res, val)
			}

			if len(res) == count {
				break
			}
		}
	}

	return res, nil

}
