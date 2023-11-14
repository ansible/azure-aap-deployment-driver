package handler

type LoginManager interface {
	GetLoginHandler() HandleFuncWithDB
}
