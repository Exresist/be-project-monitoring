package model

import "github.com/google/uuid"

const (
	RoleTeamlead ParticipantRole = iota + 1
	RoleParticipant
	RoleOwner
)

type (
	ParticipantRole int

	Participant struct {
		User
		Role      ParticipantRole
		UserID    uuid.UUID
		ProjectID int
	}
)

var Roles = map[ParticipantRole]struct{}{
	RoleTeamlead:    {},
	RoleParticipant: {},
	RoleOwner:       {},
}

//teamlead participant owner
