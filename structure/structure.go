package structure

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Program struct {
	Name      string `json:"name"`
	CreatedAt primitive.DateTime
	UpdatedAt primitive.DateTime
	Swag      bool   `json:"offers_swag"`
	Bounty    bool   `json:"offers_bounties"`
	Target    Target `json:"targets"`
}

type Target struct {
	InScope  []InScope `json:"in_scope"`
	OutScope []InScope `json:"out_of_scope"`
}

type InScope struct {
	AssetIdentifier string `json:"asset_identifier"`
	AssetType       string `json:"asset_type"` //CIDR,URL
}

type OutScope struct {
	AssetIdentifier string `json:"asset_identifier"`
	AssetType       string `json:"asset_type"` //CIDR,URL
}

type Domain struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Final     []string           `bson:"subs,omitempty"`
	Hidden    []string           `bson:"hidden,omitempty"`
	ProgramID primitive.ObjectID `bson:"program_id,omitempty"`
}
