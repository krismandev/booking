package request

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=20,password"`
	RoleID   string `json:"roleId" validate:"required"`
}

type UserListRequest struct {
	GlobalListDataRequest
	Filter string
}

type UpdateUserRequest struct {
	UserId   string `json:"userId" validate:"required"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password" validate:"omitempty,min=8,max=20,password"`
}

type DeactivateUserRequest struct {
	UserID string `json:"userId" validate:"required"`
}
