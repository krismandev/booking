package response

import "booking/model"

type RoomResponse struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	LocationID string            `json:"locationId"`
	Floor      string            `json:"floor"`
	Capacity   string            `json:"capacity"`
	CreatedAt  string            `json:"createdAt"`
	UpdatedAt  string            `json:"updatedAt"`
	Location   *LocationResponse `json:"location,omitempty"`
}

func ToRoomResponse(data model.Room, location model.Location) RoomResponse {
	var resp RoomResponse
	resp.ID = data.ID
	resp.Name = data.Name
	resp.LocationID = data.LocationID
	resp.Floor = data.Floor
	resp.Capacity = data.Capacity
	resp.CreatedAt = data.CreatedAt
	resp.UpdatedAt = data.UpdatedAt

	var loc LocationResponse
	loc.ID = location.ID
	loc.Name = location.Name

	resp.Location = &loc
	return resp
}
