package model

type Source string

type ImageInfo struct {
	Image            string      `json:"image"`
	Id               int64       `json:"id"`
	Title            string      `json:"title"`
	Description      string      `json:"description"`
	SourceId         string      `json:"-"`
	Link             string      `json:"link"`
	OwnerName        string      `json:"owner_name"`
	OwnerUrl         string      `json:"owner_url"`
	BoardName        string      `json:"board_name"`
	BoardUrl         string      `json:"board_url"`
	BoardDescription string      `json:"board_description"`
	Images           []ImageSize `json:"images"`
	KeyWords         []string    `json:"keywords"`
	Annotations      []string    `json:"annotations"`
	Hashtags         []string    `json:"hashtags"`
	CreatedTime      int64       `json:"created_time"`
	CrawledTime      int64       `json:"crawled_time"`
	Category         string      `json:"category"`
}

type ListImageInfo struct {
	Images     []ImageInfo `json:"images"`
	NextOffset int64       `json:"next_offset"`
}
type ImageSize struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Url    string `json:"url"`
}
