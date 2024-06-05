package api

// this is used to create a new user in the server database
type CreateUserRequest struct {
	Identity string
	Password string
}
