package aas

type RoleInfo struct {
	Service string `json:"service,omitempty"`
	// Name: UpdateHost
	Name string `json:"name" gorm:"not null"`
	// 1234-88769876-28768
	Context string `json:"context,omitempty"`
}

type RoleCreate struct {
	Name    string `json:"name"`
	Service string `json:"service"`
}

type RoleCreateResponse struct {
	Service string `json:"service"`
	Name    string `json:"name"`
	ID      string `json:"role_id"`
}

type RoleIDs struct {
	RoleUUIDs []string `json:"role_ids"`
}

type RoleSlice struct {
	Roles []RoleInfo `json:"roles"`
}

type UserCreate struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type UserCreateResponse struct {
	ID   string `json:"user_id"`
	Name string `json:"username"`
}

type UserRoleCreate struct {
	ID      string `json:"user_id"`
	RoleIds RoleIDs
}

type UserCred struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
