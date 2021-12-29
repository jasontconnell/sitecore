package data

import "github.com/google/uuid"

type Blob interface {
	GetId() uuid.UUID
	GetData() []byte
}

type blob struct {
	id   uuid.UUID
	data []byte
}

func (b *blob) GetId() uuid.UUID {
	return b.id
}

func (b *blob) GetData() []byte {
	return b.data
}

func NewBlob(id uuid.UUID, data []byte) Blob {
	b := &blob{id: id, data: data}
	return b
}
