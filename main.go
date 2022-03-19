package main

import (
	"log"
	"math/rand"
	"projectdeflector/matchmaking/game_lobby"
	"projectdeflector/matchmaking/repositories"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	rand.Seed(time.Now().Unix())

	repoFactory := repositories.GetRepositoryFactory()
	game_lobby.UseCase{}.MatchPeriodically(repoFactory)

	app.Use("/", func(c *fiber.Ctx) error {
		repo, cleanup, err := repoFactory.GetRepository()
		if err != nil {
			return c.SendStatus(400)
		}

		defer cleanup()
		c.Locals("repo", repo)

		return c.Next()
	})

	app.Post("/solo", func(c *fiber.Ctx) error {
		payload := struct {
			PlayerId string `json:"playerId"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			return c.SendStatus(400)
		}

		repo := c.Locals("repo").(repositories.Repository)
		useCase := game_lobby.UseCase{
			Repo: repo,
		}

		err := useCase.FindSoloGame(payload.PlayerId)
		if err != nil {
			return c.SendStatus(400)
		}

		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})
	app.Post("/find", func(c *fiber.Ctx) error {
		payload := struct {
			PlayerId string `json:"playerId"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			return c.SendStatus(400)
		}

		repo := c.Locals("repo").(repositories.Repository)
		useCase := game_lobby.UseCase{
			Repo: repo,
		}

		err := useCase.QueuePlayer(payload.PlayerId)
		if err != nil {
			return c.SendStatus(400)
		}

		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	log.Fatal(app.Listen(":3004"))
}
