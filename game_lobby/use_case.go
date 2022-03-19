package game_lobby

import (
	"projectdeflector/matchmaking/broadcast"
	"projectdeflector/matchmaking/repositories"
	"time"

	"github.com/go-co-op/gocron"
)

type UseCase struct {
	Repo repositories.Repository
}

func (useCase UseCase) MatchPeriodically(repoFactory repositories.RepositoryFactory) {

	match := func() {
		repo, cleanup, err := repoFactory.GetRepository()
		defer cleanup()
		if err != nil {
			return
		}

		group, err := repo.SetRandomMatchMakingGroup()
		if err != nil {
			return
		}

		players, err := repo.GetMatchMakingGroup(group)
		if err != nil {
			return
		}

		start := 0
		if len(players)%2 == 1 {
			start = 1
			repo.ClearPlayerMatchMakingGroup(players[0].PlayerId)
		}

		for i := start; i < len(players); i++ {
			playerIds := []string{players[i].PlayerId, players[i+1].PlayerId}
			i += 1

			gameId, err := StartGame(playerIds)
			if err == nil {
				broadcastSuccess(playerIds, gameId)
			} else {
				broadcastError(playerIds)
			}
		}

		repo.DeleteMatchMakingGroup(group)
	}

	s := gocron.NewScheduler(time.UTC)
	s.Every(5).Seconds().Do(match)

	s.StartAsync()
}

func (useCase UseCase) FindSoloGame(playerId string) error {
	gameId, err := StartGame([]string{playerId, "system"})
	if err != nil {
		return err
	}
	broadcastSuccess([]string{playerId}, gameId)

	return nil
}

func (useCase UseCase) QueuePlayer(playerId string) error {
	err := useCase.Repo.QueuePlayer(playerId)
	return err
}

func broadcastSuccess(playerIds []string, gameId string) {
	broadcast.SocketBroadcast(playerIds, "match_found", map[string]interface{}{
		"id": gameId,
	})
}

func broadcastError(playerIds []string) {
	broadcast.SocketBroadcast(playerIds, "match_found", map[string]interface{}{
		"error": true,
	})
}
