package socketServer

type socketServerDb interface {
	AddUser(id int) error
	RemoveUser(id int) error
}