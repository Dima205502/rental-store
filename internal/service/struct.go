package service

type UserManager struct {
	createrRepo RepoManager
}

type SessionManager struct {
	sessionRepo RepoManager
	finderRepo  RepoFinder
}

type EmailManager struct {
	finderRepo    RepoFinder
	userTokenRepo RepoManager
}
