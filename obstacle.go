package obstaclepack

import (
	"context"
	"fmt"

	generic "go.viam.com/rdk/components/generic"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
)

var Obstacle = resource.NewModel("viam", "obstacle-pack", "obstacle")

func init() {
	resource.RegisterComponent(generic.API, Obstacle,
		resource.Registration[resource.Resource, *ObstacleConfig]{
			Constructor: newObstaclePackObstacle,
		},
	)
}

// ObstacleConfig holds the geometries that are returned from Geometries. They are
// parsed once at construction time.
type ObstacleConfig struct {
	Geometries []spatialmath.GeometryConfig `json:"geometries"`
}

func (cfg *ObstacleConfig) ParseGeometries() ([]spatialmath.Geometry, error) {
	gs := []spatialmath.Geometry{}

	for _, gc := range cfg.Geometries {
		g, err := gc.ParseConfig()
		if err != nil {
			return nil, err
		}
		gs = append(gs, g)
	}

	return gs, nil
}

// Validate ensures the configured geometries can be parsed.
// Returns three values:
//  1. Required dependencies: other resources that must exist for this resource to work.
//  2. Optional dependencies: other resources that may exist but are not required.
//  3. An error if any Config fields are missing or invalid.
func (cfg *ObstacleConfig) Validate(path string) ([]string, []string, error) {
	_, err := cfg.ParseGeometries()
	return nil, nil, err
}

type obstaclePackObstacle struct {
	resource.AlwaysRebuild
	resource.TriviallyCloseable

	name resource.Name

	logger logging.Logger
	cfg    *ObstacleConfig

	obstacles []spatialmath.Geometry
}

func newObstaclePackObstacle(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (resource.Resource, error) {
	conf, err := resource.NativeConfig[*ObstacleConfig](rawConf)
	if err != nil {
		return nil, err
	}
	return NewObstacle(ctx, deps, rawConf.ResourceName(), conf, logger)
}

func NewObstacle(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *ObstacleConfig, logger logging.Logger) (resource.Resource, error) {
	gs, err := conf.ParseGeometries()
	if err != nil {
		return nil, err
	}

	return &obstaclePackObstacle{
		name:      name,
		logger:    logger,
		cfg:       conf,
		obstacles: gs,
	}, nil
}

func (o *obstaclePackObstacle) Name() resource.Name {
	return o.name
}

func (o *obstaclePackObstacle) Geometries(ctx context.Context, extra map[string]interface{}) ([]spatialmath.Geometry, error) {
	return o.obstacles, nil
}

func (o *obstaclePackObstacle) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, resource.ErrDoUnimplemented
}

func (o *obstaclePackObstacle) Status(ctx context.Context) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}
