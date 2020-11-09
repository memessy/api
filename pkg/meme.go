package pkg

import (
	"net/url"
	"time"
)

type Meme struct {
	Id          string
	FileUrl     url.URL
	Description string
	ParsedText  string
	Categories  []string
	ProcessedAt time.Time
	CreatedAt   time.Time
}
