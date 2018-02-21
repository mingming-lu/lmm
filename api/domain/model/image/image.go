package image

type Minimal struct {
	URL string `json:"url"`
}

type ImageType int

const (
	Blog ImageType = iota
	Avatar
	Photo
)

type Image struct {
	ID        int64
	User      int64
	Type      ImageType
	URL       string
	CreatedAt string
}
