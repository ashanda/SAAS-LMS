package handlers

import "github.com/google/uuid"

func mustUUID(value string) uuid.UUID {
	id, _ := uuid.Parse(value)
	return id
}