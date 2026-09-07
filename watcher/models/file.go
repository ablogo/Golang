package models

type Files struct {
	Type    string `json:"type"`
	Path    string `json:"path"`
	Change  bool   `json:"change"`
	Content string `json:"content"`
}
