package model

// Content shows article's content those are editable
type Content struct {
	text *Text
	tags []*Tag
}

// NewContent returns a new Content value object pointer
func NewContent(title, body string, tags []string) (*Content, error) {
	text, err := NewText(title, body)
	if err != nil {
		return nil, err
	}

	tagModels := make([]*Tag, len(tags), len(tags))
	for i, name := range tags {
		tag, err := NewTag(name, uint(i+1))
		if err != nil {
			return nil, err
		}
		tagModels[i] = tag
	}

	return &Content{
		text: text,
		tags: tagModels,
	}, nil
}

func (c *Content) Tags() []*Tag {
	return c.tags
}

func (c *Content) Text() Text {
	return *c.text
}
