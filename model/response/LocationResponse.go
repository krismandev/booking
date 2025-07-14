package response

import (
	"booking/model"
)

type LocationResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ToLocationResponse(data model.Location) LocationResponse {
	var resp LocationResponse

	resp.ID = data.ID
	resp.Name = data.Name
	return resp
}
