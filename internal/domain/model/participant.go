package model

const (
	RoleTeamlead    ParticipantRole = "TEAM_LEAD"
	RoleParticipant ParticipantRole = "PARTICIPANT"
	RoleOwner       ParticipantRole = "OWNER"
)

type (
	ParticipantRole string

	Participant struct {
		ID        int             `json:"id"`
		Role      ParticipantRole `json:"role"`
		ProjectID int             `json:"projectId,omitempty"`
		ShortUser `json:"user"`
	}
)

var ParticipantRoles = map[ParticipantRole]struct{}{
	RoleTeamlead:    {},
	RoleParticipant: {},
	RoleOwner:       {},
}
