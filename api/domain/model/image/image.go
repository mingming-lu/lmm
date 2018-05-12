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
	ID        uint64
	User      uint64
	Name      string
	CreatedAt string
}
