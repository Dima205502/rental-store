package service

type UserManager struct {
	createrRepo RepoManager
	notifier    Sender
}

type SessionManager struct {
	sessionRepo RepoManager
	finderRepo  RepoFinder
}

type EmailManager struct {
	finderRepo    RepoFinder
	userTokenRepo RepoManager
	notifier      Sender
}
