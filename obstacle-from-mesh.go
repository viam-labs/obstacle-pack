package obstaclepack

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	generic "go.viam.com/rdk/components/generic"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
)

var ObstacleFromMesh = resource.NewModel("viam", "obstacle-pack", "obstacle-from-mesh")

func init() {
	resource.RegisterComponent(generic.API, ObstacleFromMesh,
		resource.Registration[resource.Resource, *Config]{
			Constructor: newObstaclePackObstacleFromMesh,
		},
	)
}

// Config is shared by the generic and gripper variants of obstacle-from-mesh.
type Config struct {
	// MeshPath is the filesystem path to a PLY mesh file that will be loaded
	// at construction time and returned as a geometry from Geometries.
	MeshPath         string   `json:"mesh_path"`
	DecimationFactor *float64 `json:"decimation_factor,omitempty"`
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns three values:
//  1. Required dependencies: other resources that must exist for this resource to work.
//  2. Optional dependencies: other resources that may exist but are not required.
//  3. An error if any Config fields are missing or invalid.
//
// The `path` parameter indicates
// where this resource appears in the machine's JSON configuration
// (for example, "components.0"). You can use it in error messages
// to indicate which resource has a problem.
func (cfg *Config) Validate(path string) ([]string, []string, error) {
	if cfg.MeshPath == "" {
		return nil, nil, resource.NewConfigValidationFieldRequiredError(path, "mesh_path")
	}
	if !strings.EqualFold(filepath.Ext(cfg.MeshPath), ".ply") {
		return nil, nil, resource.NewConfigValidationError(path,
			fmt.Errorf("mesh_path must point to a .ply file, got %q", cfg.MeshPath))
	}
	if cfg.DecimationFactor != nil && *cfg.DecimationFactor <= 0 {
		return nil, nil, resource.NewConfigValidationError(path,
			fmt.Errorf("decimation_factor must be > 0, got %v", *cfg.DecimationFactor))
	}
	return nil, nil, nil
}

type obstaclePackObstacleFromMesh struct {
	resource.AlwaysRebuild
	resource.Named

	name resource.Name

	logger logging.Logger
	cfg    *Config

	mesh *spatialmath.Mesh

	cancelCtx  context.Context
	cancelFunc func()
}

func newObstaclePackObstacleFromMesh(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (resource.Resource, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}
	return NewObstacleFromMesh(ctx, deps, rawConf.ResourceName(), conf, logger)
}

func NewObstacleFromMesh(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *Config, logger logging.Logger) (resource.Resource, error) {
	s, err := newObstacleMeshResource(name, conf, logger)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// newObstacleMeshResource is the shared constructor used by both the generic
// and gripper variants. It loads the PLY file and initialises the struct.
func newObstacleMeshResource(name resource.Name, conf *Config, logger logging.Logger) (*obstaclePackObstacleFromMesh, error) {
	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	mesh, err := spatialmath.NewMeshFromPLYFile(conf.MeshPath)
	if err != nil {
		cancelFunc()
		return nil, fmt.Errorf("failed to load mesh from %q: %w", conf.MeshPath, err)
	}

	if conf.DecimationFactor != nil {
		mesh, err = mesh.ConservativeDecimate(len(mesh.Triangles()) / int(*conf.DecimationFactor))
		if err != nil {
			cancelFunc()
			return nil, fmt.Errorf("failed to decimate mesh: %w", err)
		}
	}

	return &obstaclePackObstacleFromMesh{
		name:       name,
		logger:     logger,
		cfg:        conf,
		mesh:       mesh,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
	}, nil
}

func (s *obstaclePackObstacleFromMesh) Name() resource.Name {
	return s.name
}

func (s *obstaclePackObstacleFromMesh) Geometries(ctx context.Context, extra map[string]interface{}) ([]spatialmath.Geometry, error) {
	if s.mesh == nil {
		return nil, fmt.Errorf("mesh not loaded")
	}
	return []spatialmath.Geometry{s.mesh}, nil
}

func (s *obstaclePackObstacleFromMesh) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, resource.ErrDoUnimplemented
}

func (s *obstaclePackObstacleFromMesh) Status(ctx context.Context) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *obstaclePackObstacleFromMesh) Close(context.Context) error {
	s.cancelFunc()
	return nil
}
