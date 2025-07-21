package response

import "booking/model"

type CreateUserResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	CreatedTime string `json:"createdTime"`
}

type UserResponse struct {
	ID        string   `json:"id"`
	CreatedAt *string  `json:"createdAt"`
	Email     string   `json:"email"`
	Name      string   `json:"name"`
	Roles     []string `json:"roles"`
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
}
