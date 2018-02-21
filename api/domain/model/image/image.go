package image

type Minimal struct {
	URL string `json:"url"`
}

type ImageType int

const (
	TypeAvatar ImageType = iota
	TypeBlog
	TypePhoto
)

const BaseURL = "https://image.lmm.im/"

type Image struct {
	ID        int64
	User      int64
	Type      ImageType
	URL       string
	CreatedAt string
}
