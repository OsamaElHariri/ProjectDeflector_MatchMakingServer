package repositories

type RepositoryFactory interface {
	GetRepository() (Repository, func(), error)
}

type Repository interface {
	QueuePlayer(playerId string) error
	UnqueuePlayer(playerId string) error
	GetMatchMakingGroup(matchMakingGroup string) ([]MatchMakingPlayer, error)
	SetRandomMatchMakingGroup() (string, error)
	DeleteMatchMakingGroup(matchMakingGroup string) error
	ClearPlayerMatchMakingGroup(playerId string) error
}

func GetRepositoryFactory() RepositoryFactory {
	return getMongoRepositoryFactory()
}
