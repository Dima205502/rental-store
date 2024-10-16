package controller

type AuthManager struct {
	creater CreaterService
	session SessionService
	checker CheckerService
}
