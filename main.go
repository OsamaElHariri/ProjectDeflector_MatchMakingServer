package main

import (
	"log"
	"projectdeflector/matchmaking/broadcast"
	"projectdeflector/matchmaking/game_board"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	playersAwaitingMatch := make([]string, 0)

	app.Post("/solo", func(c *fiber.Ctx) error {
		payload := struct {
			PlayerId string `json:"playerId"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			return c.SendStatus(400)
		}

		gameId, err := game_board.StartGame([]string{payload.PlayerId, "system"})
		if err != nil {
			return c.SendStatus(400)
		}
		broadcast.SocketBroadcast([]string{payload.PlayerId}, "match_found", map[string]interface{}{
			"id": gameId,
		})
		result := fiber.Map{
			"status": "ok",
		}
		return c.JSON(result)
	})
	app.Post("/find", func(c *fiber.Ctx) error {
		payload := struct {
			PlayerId string `json:"playerId"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			return c.SendStatus(400)
		}

		playersAwaitingMatch = append(playersAwaitingMatch, payload.PlayerId)

		if len(playersAwaitingMatch) == 2 {
			gameId, err := game_board.StartGame(playersAwaitingMatch)
			if err != nil {
				return c.SendStatus(400)
			}
			broadcast.SocketBroadcast(playersAwaitingMatch, "match_found", map[string]interface{}{
				"id": gameId,
			})
			playersAwaitingMatch = make([]string, 0)
		}
		result := fiber.Map{
			"status": "ok",
		}
		return c.JSON(result)
	})

	log.Fatal(app.Listen(":3004"))
}
