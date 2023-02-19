package model

const (
	RoleTeamlead    ParticipantRole = "team_lead"
	RoleParticipant ParticipantRole = "participant"
	RoleOwner       ParticipantRole = "owner"
)

type (
	ParticipantRole string

	Participant struct {
		ID        int             `json:"id"`
		Role      ParticipantRole `json:"role"`
		ProjectID int             `json:"project_id,omitempty"`
		ShortUser `json:"User"`
	}
)

var ParticipantRoles = map[ParticipantRole]struct{}{
	RoleTeamlead:    {},
	RoleParticipant: {},
	RoleOwner:       {},
}

//teamlead participant owner
