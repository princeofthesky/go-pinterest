package model

type Source string

const (
	Pinterest Source = Source(iota)
	Reddit
)

func (s Source) String() string {
	switch s {
	case Pinterest:
		return "Pinterest"
	case Reddit:
		return "Reddit"
	}
	return "unknown"
}

type ImageInfo struct {
	Image                  string
	Id                     uint64
	Title                  string
	Source                 Source
	SourceId               string
	SourceLink             string
	SourceOwnerName          string
	SourceOwnerLink        string
	SourceImageSizeDeTails []ImageSize
	SourceKeyWord          string
	Category               string
	CreatedTime            uint64
	CrawledTime            uint64
}

type ImageSize struct {
	Width  int
	Height int
	Url    string
}
type ImageResponse struct {
	Avatar string `json:"ava"`
	Id     int64  `json:"id"`
	Title  string `json:"title"`
}
