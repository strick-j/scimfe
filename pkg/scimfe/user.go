package scimfe

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UsersResponse struct {
	Users []User `json:"users"`
}

func (c Client) Users(t Token) ([]User, error) {
	rsp := new(UsersResponse)
	return rsp.Users, c.get("/users", rsp, t)
}

func (c Client) UserByID(uid string, t Token) (*User, error) {
	rsp := new(User)
	return rsp, c.get("/users/"+uid, rsp, t)
}

func (c Client) CurrentUser(t Token) (*User, error) {
	rsp := new(User)
	return rsp, c.get("/users/self", rsp, t)
}
