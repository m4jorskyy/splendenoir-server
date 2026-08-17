package inpost

type SendParcelRequest struct {
	Receiver struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	} `json:"receiver"`
	CustomAttributes struct {
		TargetPoint string `json:"target_point"`
	} `json:"custom_attributes"`
}
