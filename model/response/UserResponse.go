package response

import "booking/model"

type CreateUserResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	CreatedTime string `json:"createdTime"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	CreatedAt *string   `json:"createdAt"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Roles     *[]string `json:"roles,omitempty"`
	Role      string    `json:"role"`
}

func ToUserResponse(dt model.User) UserResponse {
	var resp UserResponse

	resp.ID = dt.ID
	resp.Email = dt.Email
	resp.Name = dt.Name
	resp.CreatedAt = dt.CreatedAt

	return resp
}

type UpdateUserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type AuthUserDetailResponse struct {
	User UserResponse `json:"user"`
	Role RoleResponse `json:"role"`
}

type RoleResponse struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Privileges []string `roles`
}

type UserListResponse struct {
	GlobalListDataResponse
	Data []UserResponse
}
