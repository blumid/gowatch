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
type Subdomain struct {
	// ID        primitive.ObjectID `bson:"_id"`
	ProgramID primitive.ObjectID `bson:"programID,omitempty"`
	Sub       string             `bson:"subdomain,omitempty"`
	SC        int                `bson:"sc"` //status code
	// CL        int                `bson:"cl"` //content length
	Locatoin string `bson:"location,omitempty"`
	Detail   Detail `bson:"detail,omitempty"`
}

type Detail struct {
	CDN     bool   `bson:"cdn"`
	Icon    string `bson:"icon,omitempty"`
	A       []string
	Cname   []string
	Tech    []string               // technology
	Headers map[string]interface{} //response headers
}
