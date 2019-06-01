package model

// Content shows article's content those are editable
type Content struct {
	text *Text
	tags []string
}

// NewContent returns a new Content value object pointer
func NewContent(text *Text, tags []string) *Content {
	return &Content{
		text: text,
		tags: tags,
	}
}

func (c *Content) Tags() []string {
	return c.tags
}

func (c *Content) Text() Text {
	return *c.text
}
