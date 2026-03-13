package dto

// CheckDBRequest defines setup database connectivity validation payload.
type CheckDBRequest struct {
	Host string `json:"host" binding:"required"`
	Port int    `json:"port" binding:"required"`
	User string `json:"user" binding:"required"`
	Pass string `json:"pass"`
	Name string `json:"name" binding:"required"`
}
