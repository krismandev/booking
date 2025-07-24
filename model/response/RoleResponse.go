package response

import "booking/model"

type RoleResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Privileges string `json:"privileges"`
}

func ToRoleResponse(dt model.Role) RoleResponse {
	var resp RoleResponse

	resp.ID = dt.ID
	resp.Name = dt.Name
	resp.Privileges = dt.Privileges

	return resp
}
