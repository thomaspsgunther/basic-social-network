package users

type WrongUsernameOrPasswordError struct{}
type UserAlreadyExistsError struct{}

func (m *WrongUsernameOrPasswordError) Error() string {
	return "wrong username or password"
}

func (m *UserAlreadyExistsError) Error() string {
	return "user already exists"
}
