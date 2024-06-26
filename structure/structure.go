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
	ProgramID primitive.ObjectID `bson:"programID,omitempty"`
	Sub       string             `bson:"subdomain,omitempty"`
	SC        int                `bson:"sc"` // status code
	Locatoin  string             `bson:"location,omitempty"`
	CDN       bool               `bson:"cdn"`
	Icon      string             `bson:"icon,omitempty"`
	Detail    Detail             `bson:"detail,omitempty"`
}

type Detail struct {
	A       []string               // dns A record
	Cname   []string               // dns CNAME record
	Tech    []string               // technologies
	Headers map[string]interface{} // response headers
}
