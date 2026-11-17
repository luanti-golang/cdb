package cdb

type VCSPackage struct {
	Aliases          []string `json:"aliases,omitempty"`
	Author           string   `json:"author"`
	Featured         bool     `json:"featured,omitempty"`
	Name             string   `json:"name"`
	Release          int      `json:"release"`
	Repo             *string  `json:"repo"`
	ShortDescription string   `json:"short_description"`
	Thumbnail        *string  `json:"thumbnail"`
	Title            string   `json:"title"`
	Type             string   `json:"type"` // "game", "mod" or "txp"
}
