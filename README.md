# ProjectDeflector MatchMakingServer

This repo holds the code to handle queueing players and matching them together for the `Hit Bounce` mobile game.


## Mobile App

The mobile app, and a high level introduction to this project can be found in the [JsClientGame](https://github.com/OsamaElHariri/ProjectDeflector_JsClientGame) repo, which is intended to be the entry point to understanding this project.


## Overview of This Project

This is a Go server that uses the [Fiber](https://gofiber.io/) web framework. The routes are just in `main.go`, and the function of the routes is to only validate the inputs, then call a use case (which is what runs the business logic). The use cases can be found in the `use_cases.go` file, and these should incapsulate all the functions that this server can do.

Note that this project has a `.devcontainer` and is meant to be run inside a dev container.


## Outputting a Binary

To output the binary of this Go code, run the VSCode task using `CTRL+SHIFT+B`. This should be done while inside the dev container.


Once you have the binary, you need to build the docker image _outside_ the dev container. I use this command and just overwrite the image everytime. This keeps the [Infra](https://github.com/OsamaElHariri/ProjectDeflector_Infra) repo simpler.

```
docker build -t project_deflector/matchmaking_server:1.0 .
```

## Matchmaking Overview

I kept it real simple with this one. The matchmaking strategy is simply a function that runs every 5 seconds. This function loops over all the players in the queue, and just matches the player at index `i` with the player at index `i+1`.

