package image

type Minimal struct {
	Name string `json:"name"`
}

type ImageData struct {
	Name string
	Data []byte
}

const BaseURL = "https://image.lmm.im/"

type Image struct {
	ID        int64
	User      int64
	Name      string
	CreatedAt string
}
