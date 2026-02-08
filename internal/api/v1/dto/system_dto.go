package dto

type SystemStatusResponse struct {
	Installed bool    `json:"installed"`
	SiteName  *string `json:"site_name,omitempty"`
}

type SystemSetupRequest struct {
	SiteName      string `json:"site_name" binding:"required,min=1,max=100"`
	AdminUsername string `json:"admin_username" binding:"required,min=1,max=50"`
	AdminEmail    string `json:"admin_email" binding:"required,email,max=255"`
	AdminPassword string `json:"admin_password" binding:"required,min=8,max=72"`
}
