package main

import (
	"log"
	"math/rand"
	"os"
	"projectdeflector/matchmaking/game_lobby"
	"projectdeflector/matchmaking/repositories"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}
	err := godotenv.Load("env/." + env + ".env")
	if err != nil {
		log.Fatalf("could not load env vars")
	}

	app := fiber.New()
	app.Use(recover.New())

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

	app.Use("/", func(c *fiber.Ctx) error {
		userId := c.Get("x-user-id")
		if userId != "" {
			c.Locals("userId", userId)
		}
		return c.Next()
	})

	app.Post("/solo", func(c *fiber.Ctx) error {
		playerId := c.Locals("userId").(string)
		repo := c.Locals("repo").(repositories.Repository)
		useCase := game_lobby.UseCase{
			Repo: repo,
		}

		err := useCase.FindSoloGame(playerId)
		if err != nil {
			return c.SendStatus(400)
		}

		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	app.Post("/find", func(c *fiber.Ctx) error {
		playerId := c.Locals("userId").(string)

		repo := c.Locals("repo").(repositories.Repository)
		useCase := game_lobby.UseCase{
			Repo: repo,
		}

		err := useCase.QueuePlayer(playerId)
		if err != nil {
			return c.SendStatus(400)
		}

		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	app.Post("/cancel", func(c *fiber.Ctx) error {
		playerId := c.Locals("userId").(string)

		repo := c.Locals("repo").(repositories.Repository)
		useCase := game_lobby.UseCase{
			Repo: repo,
		}

		err := useCase.UnqueuePlayer(playerId)
		if err != nil {
			return c.SendStatus(400)
		}

		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	log.Fatal(app.Listen(":3004"))
}
