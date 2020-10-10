package accounts

type PermissionGroup string

// A list of common expected errors.
const (
	GeneralUser       PermissionGroup = "general"
	OperationTeam PermissionGroup = "operation"
)

type Account struct {
	ID       string `json:"id"    bson:"_id"`
	Name     string `json:"name"    bson:"name"`
	Password string `json:"password"    bson:"password"`
	Balance int `json:"balance"    bson:"balance"`
	UserGroup PermissionGroup `json:"userGroup"    bson:"userGroup"`
}

type AccountReadModel struct {
	ID       string `json:"id"    bson:"_id"`
	Name     string `json:"name"    bson:"name"`
	Balance int `json:"balance"    bson:"balance"`
}

type AccountLoginParam struct {
	Email    string `json:"email"    bson:"email"`
	Password string `json:"password"    bson:"password"`
}

type LoginResponse struct {
	Status bool   `json:"status"    bson:"status"`
	Token  string `json:"token"    bson:"token"`
}
