package mock

import (
	"testing"

	"githib.com/VladimirStepanov/snippetbox/pkg/models"
)

type SnippetData struct {
	Title    string
	Content  string
	Expire   int
	IsPublic bool
}

func getPreparedSnippetStore(t *testing.T) (*SnippetStore, int64) {
	us := &UsersStore{DB: map[int64]*models.User{}}

	userID, err := us.Insert("test", "test", "test", "test")

	if err != nil {
		t.Fatal(err)
	}

	return &SnippetStore{DB: []*models.Snippet{}, UsersMap: us.DB}, userID
}

func TestInsertSnippet(t *testing.T) {
	tests := map[string]struct {
		WantError  error
		Data       *SnippetData
		GetOwnerID func(id int64) int64
	}{
		"Success insert": {
			WantError: nil,
			Data: &SnippetData{
				Title:   "Title",
				Content: "Content",
				Expire:  1,
			},
			GetOwnerID: func(id int64) int64 {
				return id
			},
		},
		"Unknown owner ID": {
			WantError: models.ErrUnknownOwnerID,
			Data: &SnippetData{
				Title:   "Title",
				Content: "Content",
				Expire:  1,
			},
			GetOwnerID: func(id int64) int64 {
				return id + 10
			},
		},
	}
	for name, value := range tests {
		t.Run(name, func(t *testing.T) {
			ss, userID := getPreparedSnippetStore(t)

			_, err := ss.Insert(value.Data.Title, value.Data.Content, value.Data.Expire, value.Data.IsPublic, value.GetOwnerID(userID))

			if value.WantError != nil && value.WantError != err {
				t.Fatalf("Want: %v, Get: %v\n", value.WantError, err)
			}

			if value.WantError == nil && err != nil {
				t.Fatal(err)
			}
		})
	}

}

func TestGetSnippet(t *testing.T) {
	tests := map[string]struct {
		WantError error
		Data      *SnippetData
		GetID     func(id int64) int64
	}{
		"Get ErrNoRecord": {
			WantError: models.ErrNoRecord,
			Data:      nil,
		},
		"Success get": {
			WantError: nil,
			Data: &SnippetData{
				Title:   "Title",
				Content: "Content",
				Expire:  1,
			},
		},
		"Snippet expire": {
			WantError: models.ErrNoRecord,
			Data: &SnippetData{
				Title:   "Title",
				Content: "Content",
				Expire:  -1,
			},
		},
	}

	for name, value := range tests {
		t.Run(name, func(t *testing.T) {
			ss, ownerID := getPreparedSnippetStore(t)

			var snippetID int64
			var err error

			if value.Data != nil {
				snippetID, err = ss.Insert(value.Data.Title, value.Data.Content, value.Data.Expire, value.Data.IsPublic, ownerID)
				if err != nil {
					t.Fatal(err)
				}
			}

			snippet, err := ss.Get(snippetID)

			if value.WantError != nil && value.WantError != err {
				t.Fatalf("Want: %v, Get: %v\n", value.WantError, err)
			}

			if value.WantError == nil && err != nil {
				t.Fatal(err)
			}

			if snippet != nil && value.Data != nil && value.WantError == nil {
				if snippet.Content != value.Data.Content || snippet.IsPublic != value.Data.IsPublic || snippet.Title != value.Data.Title || snippet.OwnerID != ownerID {
					t.Fatalf("Want: %v, Get: %v", snippet, value.Data)
				}
			}
		})
	}
}

func TestDeleteSnippet(t *testing.T) {
	tests := map[string]struct {
		WantError  error
		Data       *SnippetData
		GetID      func(id int64) int64
		GetOwnerID func(id int64) int64
	}{
		"Delete row not found": {
			WantError: models.ErrNoRecord,
			Data:      nil,
			GetID: func(id int64) int64 {
				return id + 10
			},
			GetOwnerID: func(id int64) int64 {
				return id
			},
		},
		"Delete row not found because bad user": {
			WantError: models.ErrNoRecord,
			Data: &SnippetData{
				Title:   "Title",
				Content: "Content",
				Expire:  1,
			},
			GetID: func(id int64) int64 {
				return id
			},
			GetOwnerID: func(id int64) int64 {
				return id + 1
			},
		},
		"Success delete": {
			WantError: nil,
			Data: &SnippetData{
				Title:   "Title",
				Content: "Content",
				Expire:  1,
			},
			GetID: func(id int64) int64 {
				return id
			},
			GetOwnerID: func(id int64) int64 {
				return id
			},
		},
	}

	for name, value := range tests {
		t.Run(name, func(t *testing.T) {
			ss, ownerID := getPreparedSnippetStore(t)

			var snippetID int64
			var err error

			if value.Data != nil {
				snippetID, err = ss.Insert(value.Data.Title, value.Data.Content, value.Data.Expire, value.Data.IsPublic, ownerID)
				if err != nil {
					t.Fatal(err)
				}
			}

			err = ss.Delete(value.GetID(snippetID), value.GetOwnerID(ownerID))

			if value.WantError != nil && value.WantError != err {
				t.Fatalf("Want: %v, Get: %v\n", value.WantError, err)
			}

			if value.WantError == nil && err != nil {
				t.Fatal(err)
			}
		})
	}

}

func TestUpdate(t *testing.T) {
	tests := map[string]struct {
		WantError             error
		Data                  *SnippetData
		GetSnippetAfterUpdate func(s *SnippetData) *SnippetData
	}{
		"Update ErrNoRecord": {
			WantError: models.ErrNoRecord,
			Data:      nil,
			GetSnippetAfterUpdate: func(s *SnippetData) *SnippetData {
				return &SnippetData{}
			},
		},
		"Update success": {
			WantError: nil,
			Data: &SnippetData{
				Title:   "Title",
				Content: "Content",
				Expire:  1,
			},
			GetSnippetAfterUpdate: func(s *SnippetData) *SnippetData {
				s.Title += "Hello world"
				return s
			},
		},
	}

	for name, value := range tests {
		t.Run(name, func(t *testing.T) {
			ss, ownerID := getPreparedSnippetStore(t)
			var snippetID int64
			var err error

			if value.Data != nil {
				snippetID, err = ss.Insert(
					value.Data.Title, value.Data.Content, value.Data.Expire, value.Data.IsPublic, ownerID,
				)
				if err != nil {
					t.Fatal(err)
				}
			}

			updatedSnippet := value.GetSnippetAfterUpdate(value.Data)

			err = ss.Update(
				&models.Snippet{ID: snippetID,
					Title:    updatedSnippet.Title,
					Content:  updatedSnippet.Content,
					IsPublic: updatedSnippet.IsPublic,
				},
				ownerID,
			)

			if value.WantError != nil && value.WantError != err {
				t.Fatalf("Want: %v, Get: %v\n", value.WantError, err)
			}

			if value.WantError == nil && err != nil {
				t.Fatal(err)
			}

			if value.Data != nil && snippetID != 0 {
				snippet, err := ss.Get(snippetID)

				if err != nil {
					t.Fatalf("Error while get: %v %d %d", err, snippetID, ownerID)
				}

				if snippet.Content != updatedSnippet.Content || snippet.IsPublic != updatedSnippet.IsPublic || snippet.Title != updatedSnippet.Title {
					t.Fatalf("Want: %v, Get: %v", snippet, updatedSnippet)
				}
			}

		})
	}
}

func TestLatestAll(t *testing.T) {
	snippets := []*SnippetData{
		{"1", "2", 1, true},
		{"11", "22", 1, false},
		{"111", "222", 1, true},
		{"exp11", "exp22", -1, false},
	}

	tests := map[string]struct {
		WantError  error
		WantResult []*SnippetData
		WantLatest int
		WantPage   int
		GetOwnerID func(int64) int64
	}{
		"All public": {
			WantError:  nil,
			WantResult: []*SnippetData{snippets[0], snippets[2]},
			WantLatest: 5,
			WantPage:   1,
			GetOwnerID: func(id int64) int64 {
				return -1
			},
		},
		"All for owner ID": {
			WantError:  nil,
			WantResult: snippets[:3],
			WantLatest: 5,
			WantPage:   1,
			GetOwnerID: func(id int64) int64 {
				return id
			},
		},
		"Empty latest": {
			WantError:  nil,
			WantResult: []*SnippetData{},
			WantLatest: 5,
			WantPage:   1,
			GetOwnerID: func(id int64) int64 {
				return id + 3
			},
		},
		"Test first page": {
			WantError:  nil,
			WantResult: snippets[:2],
			WantLatest: 2,
			WantPage:   1,
			GetOwnerID: func(id int64) int64 {
				return id
			},
		},
		"Test last page": {
			WantError:  nil,
			WantResult: []*SnippetData{snippets[2]},
			WantLatest: 2,
			WantPage:   2,
			GetOwnerID: func(id int64) int64 {
				return id
			},
		},
		"Test third page": {
			WantError:  nil,
			WantResult: []*SnippetData{snippets[2]},
			WantLatest: 1,
			WantPage:   3,
			GetOwnerID: func(id int64) int64 {
				return id
			},
		},
	}

	for name, value := range tests {
		t.Run(name, func(t *testing.T) {
			ss, ownerID := getPreparedSnippetStore(t)

			for _, snippet := range snippets {
				_, err := ss.Insert(snippet.Title, snippet.Content, snippet.Expire, snippet.IsPublic, ownerID)
				if err != nil {
					t.Fatal(err)
				}
			}

			latestSnippets, err := ss.LatestAll(value.GetOwnerID(ownerID), value.WantLatest, value.WantPage)
			if value.WantError != nil && value.WantError != err {
				t.Fatalf("Want: %v, Get: %v\n", value.WantError, err)
			}

			if value.WantError == nil && err != nil {
				t.Fatal(err)
			}

			if len(value.WantResult) > 0 {
				if len(value.WantResult) != len(latestSnippets) {
					t.Fatalf("Want len %d, Got len: %d", len(value.WantResult), len(latestSnippets))
				}

				for i := range latestSnippets {
					if latestSnippets[i].Content != value.WantResult[i].Content || latestSnippets[i].IsPublic != value.WantResult[i].IsPublic || latestSnippets[i].Title != value.WantResult[i].Title {
						t.Fatalf("Want: %v, Get: %v", latestSnippets[i], value.WantResult[i])
					}
				}
			} else if len(value.WantResult) == 0 && len(latestSnippets) != 0 {
				t.Fatalf("Want len: 0, got len: %d", len(latestSnippets))
			}

		})
	}
}
