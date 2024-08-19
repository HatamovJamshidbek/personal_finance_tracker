package models

type PrimaryKey struct {
	ID string `json:"id"`
}

type GetListRequest struct {
	Page   int32  `json:"page"`
	Limit  int32  `json:"limit"`
	Search string `json:"search"`
}

type Response struct {
	StatusCode  int
	Description string
	Data        interface{}
}

type UpdatePasswordRequest struct {
	Login       string `json:"login"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
