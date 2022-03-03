package models

type Skill struct {
	Id           string `json:"id"`
	SkillId      string `json:"name"`
	UserId       string `json:"description"`
	Exp          string `json:"exp"`
	Txp          string `json:"txp"`
	Level        string `json:"level"`
	DateCreated  string `json:"dateCreated"`
	DateModified string `json:"dateModified"`
}
