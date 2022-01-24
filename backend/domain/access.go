package domain

type Access struct {
	Token string
	User  User
	Poll  Poll
}

type AccessRepository interface {
	GetByAccessToken(token string) (Access, error)
	GetAccessAsAdmin(admin User) ([]Access, error)
	GetAllAccess(poll Poll) ([]Access, error)
	GrantAccess(user User, poll Poll) (Access, error)
	RevokeAccess(user User, poll Poll) error
}
