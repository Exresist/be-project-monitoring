package model

const (
	RoleTeamlead    ParticipantRole = "team_lead"
	RoleParticipant ParticipantRole = "participant"
	RoleOwner       ParticipantRole = "owner"
)

type (
	ParticipantRole string

	Participant struct {
		ShortUser
		Role      ParticipantRole `json:"role"`
		ID        int             `json:"id"`
		ProjectID int             `json:"project_id"`
	}
)

var ParticipantRoles = map[ParticipantRole]struct{}{
	RoleTeamlead:    {},
	RoleParticipant: {},
	RoleOwner:       {},
}

//teamlead participant owner
