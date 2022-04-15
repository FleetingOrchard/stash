package models

type Bookmark struct {
	ID       int    `db:"id" json:"id"`
	URL      string `db:"url" json:"url"`
	Name     string `db:"name" json:"name"`
	Position int    `db:"position" json:"position"`
}

func NewBookmark(url string, name string) *Bookmark {
	return &Bookmark{
		URL:  url,
		Name: name,
	}
}

type Bookmarks []*Bookmark

func (t *Bookmarks) Append(o interface{}) {
	*t = append(*t, o.(*Bookmark))
}

func (t *Bookmarks) New() interface{} {
	return &Bookmark{}
}
