package tag

type Minimal struct {
	Name string `json:"name"`
}

type Tag struct {
	ID   uint64 `json:"id"`
	User uint64 `json:"user"`
	Blog uint64 `json:"blog"`
	Name string `json:"name"`
}
