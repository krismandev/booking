package response

type CreateUserResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	CreatedTime string `json:"createdTime"`
}
