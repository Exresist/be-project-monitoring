package model

const (
	RoleTeamlead    ParticipantRole = "Teamlead"
	RoleParticipant ParticipantRole = "Participant"
	RoleOwner       ParticipantRole = "Owner"
)

type (
	ParticipantRole string

	Participant struct {
		User
		Role      ParticipantRole
		ID        int
		ProjectID int
	}
)

var ParticipantRoles = map[ParticipantRole]struct{}{
	RoleTeamlead:    {},
	RoleParticipant: {},
	RoleOwner:       {},
}

//teamlead participant owner
