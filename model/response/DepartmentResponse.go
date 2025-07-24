package response

import "booking/model"

type DepartmentResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ToDepartmentResponse(data model.Department) DepartmentResponse {
	var resp DepartmentResponse

	resp.ID = data.ID
	resp.Name = data.Name
	return resp
}
