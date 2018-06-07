package model

type Entity struct {
	id uint64
}

func (e *Entity) ID() uint64 {
	return e.id
}

func NewEntity(id uint64) *Entity {
	return &Entity{id: id}
}
