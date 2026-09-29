package main

import (
	"obstaclepack"

	generic "go.viam.com/rdk/components/generic"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/worldstatestore"
)

func main() {
	module.ModularMain(
		resource.APIModel{API: generic.API, Model: obstaclepack.ObstacleFromMesh},
		resource.APIModel{API: generic.API, Model: obstaclepack.Obstacle},
		resource.APIModel{API: worldstatestore.API, Model: obstaclepack.WorldStateAggregator},
	)
}
