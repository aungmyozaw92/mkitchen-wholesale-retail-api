package admin

type CreateUsersRequest struct {
	Name string `validate:"required,min=3,max=200"`
	Username string `validate:"required,min=4,max=20"`
    Email    string `validate:"required,email"`
	Password string `validate:"required,min=6,max=20"`
}