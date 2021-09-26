package redisStore

type IUserManager interface {
	UserIsActive(userId string) (bool, error)
	ActiveUser(userId string) error
	DeactiveUser(userId string) error
}
