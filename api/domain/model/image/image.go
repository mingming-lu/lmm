package image

type Minimal struct {
	URL string `json:"url"`
}

type ImageType int

const (
	TypeBlog ImageType = iota
	TypeAvatar
	TypePhoto
)

type Image struct {
	ID        int64
	User      int64
	Type      ImageType
	URL       string
	CreatedAt string
}
