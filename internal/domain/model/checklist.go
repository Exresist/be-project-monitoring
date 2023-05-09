package model

type Checklist struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ProjectID int    `json:"projectId"`
	Checked   bool   `json:"checked"`
}
