package models

type PrimaryKey struct {
	ID string `bson:"id"`
}

type GetListRequest struct {
	Page   int32  `bson:"page"`
	Limit  int32  `bson:"limit"`
	Search string `bson:"search"`
}

type Response struct {
	StatusCode  int
	Description string
	Data        interface{}
}

type UpdatePasswordRequest struct {
	Login       string `bson:"login"`
	OldPassword string `bson:"old_password"`
	NewPassword string `bson:"new_password"`
}
