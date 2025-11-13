package domain

// Preferences represents trip preferences and rules
type Preferences struct {
	PetsAllowed    bool `json:"pets_allowed" bson:"pets_allowed"`
	SmokingAllowed bool `json:"smoking_allowed" bson:"smoking_allowed"`
	MusicAllowed   bool `json:"music_allowed" bson:"music_allowed"`
}
