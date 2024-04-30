package structure

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ------- hackerOne ----------
type Program struct {
	Name      string `json:"name"`
	Url       string `json:"url"`
	Target    Target `json:"targets"`
	Owner     string
	Bounty    string `json:"bounty"`
	CreatedAt primitive.DateTime
	UpdatedAt primitive.DateTime
}

type Target struct {
	InScope  []InScope  `json:"in_scope"`
	OutScope []OutScope `json:"out_of_scope"`
}

type InScope struct {
	Asset string `json:"asset"`
	Type  string `json:"type"` //CIDR,URL,WILDCARD,API
}

type OutScope struct {
	Asset string `json:"asset"`
	Type  string `json:"type"` //CIDR,URL,WILDCARD,API
}

type Result_1 struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Target Target             `bson:"target"`
}

type Message struct {
	Username  string `json:"username,omitempty"`
	AvatarUrl string `json:"avatar_url,omitempty"`
	Content   string `json:"content,omitempty"`
	// Embeds    []Embed `json:"embeds,omitempty"`
}

// --------------------- asset collection: --------------------
type Asset struct {
	ID        primitive.ObjectID `bson:"_id"`
	ProgramID primitive.ObjectID `bson:"program_id,omitempty"`
	WILDCARD  string             `bson:"name"`
	Subs      []Sub              `bson:"subs,omitempty"`
}

type Sub struct {
	DName string // domain name
	SC    int    //status code
	CT    int    //content length
	// Dns     Dns
	Tech    []string // technology
	Headers []string //response headers
}

type Dns struct {
	A     string
	Cname string
	Mail  string
}
