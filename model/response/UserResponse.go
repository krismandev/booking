package response

import "booking/model"

type CreateUserResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	CreatedTime string `json:"createdTime"`
}

type UserResponse struct {
	ID           string              `json:"id"`
	CreatedAt    *string             `json:"createdAt"`
	Email        string              `json:"email"`
	Name         string              `json:"name"`
	Roles        *[]string           `json:"roles,omitempty"`
	Role         *RoleResponse       `json:"role,omitempty"`
	IsActive     bool                `json:"isActive"`
	DepartmentID string              `json:"departmentid"`
	Department   *DepartmentResponse `json:"deparment"`
}

func ToUserResponse(dt model.User, role *model.Role) UserResponse {
	var resp UserResponse

	resp.ID = dt.ID
	resp.Email = dt.Email
	resp.Name = dt.Name
	resp.CreatedAt = dt.CreatedAt
	resp.IsActive = dt.IsActive
	resp.DepartmentID = dt.DepartmentID

	if role != nil {
		roleResp := ToRoleResponse(*role)
		resp.Role = &roleResp
	}

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

type UserListResponse struct {
	GlobalListDataResponse
	Data []UserResponse
}
