package r_models

type PetRequest struct {
	PetID  uint `json:"pet_id"`  // ID of the pet being requested
	UserID uint `json:"user_id"` // ID of the user requesting
}

type PetFosterHomeContact struct {
	PetID         uint   `json:"pet_id"`          // ID of the pet for which contact is being made
	UserID        uint   `json:"user_id"`         // ID of the user making the contact
	ContactUserID uint   `json:"contact_user_id"` // ID of the foster home contact user
	Reason        string `json:"reason"`          // Reason for contacting the foster family
	Message       string `json:"message"`         // Message to be sent to the foster family
}

type Base64Image struct {
	Base64 string `json:"base64"` // Base64 encoded image string
	Name   string `json:"name"`   // Name of the image file
}
