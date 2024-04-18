package structure

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ------- hackerOne ----------
type Program struct {
	Name      string `json:"name"`
	CreatedAt primitive.DateTime
	UpdatedAt primitive.DateTime
	Url       string `json:"url"`
	Target    Target `json:"targets"`
}

type Target struct {
	InScope  []InScope  `json:"in_scope"`
	OutScope []OutScope `json:"out_of_scope"`
}

type InScope struct {
	Asset string `json:"asset_identifier"`
	Type  string `json:"asset_type"` //CIDR,URL,WILD
}

type OutScope struct {
	Asset     string `json:"asset_identifier"`
	AssetType string `json:"asset_type"` //CIDR,URL,WILD
}

// ------------------------------

// -------- Intigriti -----------

type Program_I struct {
	Name      string `json:"name"`
	CreatedAt primitive.DateTime
	UpdatedAt primitive.DateTime
	Url       string `json:"url"`
	Target    struct {
		InScope  []I
		OutScope []I
	} `json:"targets"`
}

type I struct {
	Asset string `json:"endpoint"`
	Type  string `json:"type"` //CIDR,URL,WILD
}

// ------------------------------

// -------- BugCrowd ------------

type Program_B struct {
	Name      string `json:"name"`
	CreatedAt primitive.DateTime
	UpdatedAt primitive.DateTime
	Url       string `json:"url"`
	Target    struct {
		InScope  []B
		OutScope []B
	} `json:"targets"`
}

type B struct {
	Asset string `json:"target"`
	Type  string `json:"type"` //CIDR,URL,WILD
}

// ------------------------------

type Result_1 struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Target Target             `bson:"target"`
}

type Domain struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Subs      []Sub              `bson:"subs,omitempty"`
	ProgramID primitive.ObjectID `bson:"program_id,omitempty"`
}

type Sub struct {
	Name   string
	Hidden bool
}

type Message struct {
	Username  string `json:"username,omitempty"`
	AvatarUrl string `json:"avatar_url,omitempty"`
	Content   string `json:"content,omitempty"`
	// Embeds    []Embed `json:"embeds,omitempty"`
}
