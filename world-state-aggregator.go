package obstaclepack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	commonpb "go.viam.com/api/common/v1"
	pb "go.viam.com/api/service/worldstatestore/v1"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/worldstatestore"
	"google.golang.org/protobuf/encoding/protojson"
)

var WorldStateAggregator = resource.NewModel("viam", "obstacle-pack", "world-state-aggregator")

const defaultSubscriberBufferSize = 64

func init() {
	resource.RegisterService(worldstatestore.API, WorldStateAggregator,
		resource.Registration[worldstatestore.Service, *WorldStateAggregatorConfig]{
			Constructor: newWorldStateAggregator,
		},
	)
}

// WorldStateAggregatorConfig is the configuration for the world-state-aggregator service.
type WorldStateAggregatorConfig struct {
	// SubscriberBufferSize bounds each subscriber's channel. Default 64 when zero
	// or unset. When a subscriber's buffer fills, changes are dropped and logged
	// for that subscriber only; producers are never blocked.
	SubscriberBufferSize int `json:"subscriber_buffer_size,omitempty"`
}

// Validate ensures the config is well-formed. The aggregator has no upstream
// dependencies, so this never returns required/optional dep lists.
func (c *WorldStateAggregatorConfig) Validate(path string) ([]string, []string, error) {
	return nil, nil, nil
}

func newWorldStateAggregator(
	ctx context.Context,
	deps resource.Dependencies,
	conf resource.Config,
	logger logging.Logger,
) (worldstatestore.Service, error) {
	cfg, err := resource.NativeConfig[*WorldStateAggregatorConfig](conf)
	if err != nil {
		return nil, err
	}
	bufSize := cfg.SubscriberBufferSize
	if bufSize <= 0 {
		bufSize = defaultSubscriberBufferSize
	}
	return &worldStateAggregator{
		Named:       conf.ResourceName().AsNamed(),
		logger:      logger,
		store:       newTransformStore(),
		subscribers: make(map[uint64]*subscriber),
		bufferSize:  bufSize,
	}, nil
}

type worldStateAggregator struct {
	resource.Named
	resource.TriviallyReconfigurable
	resource.TriviallyCloseable

	logger logging.Logger
	store  *transformStore

	subMu       sync.Mutex
	subscribers map[uint64]*subscriber
	nextSubID   uint64
	bufferSize  int
}

type subscriber struct {
	id uint64
	ch chan worldstatestore.TransformChange
}

func (a *worldStateAggregator) ListUUIDs(ctx context.Context, extra map[string]any) ([][]byte, error) {
	return a.store.list(), nil
}

func (a *worldStateAggregator) GetTransform(ctx context.Context, uuid []byte, extra map[string]any) (*commonpb.Transform, error) {
	if t := a.store.get(string(uuid)); t != nil {
		return t, nil
	}
	return nil, worldstatestore.ErrNilResponse
}

func (a *worldStateAggregator) StreamTransformChanges(ctx context.Context, extra map[string]any) (*worldstatestore.TransformChangeStream, error) {
	sub := a.register()
	go func() {
		<-ctx.Done()
		a.unregister(sub)
	}()
	return worldstatestore.NewTransformChangeStreamFromChannel(ctx, sub.ch), nil
}

func (a *worldStateAggregator) register() *subscriber {
	a.subMu.Lock()
	defer a.subMu.Unlock()
	a.nextSubID++
	s := &subscriber{
		id: a.nextSubID,
		ch: make(chan worldstatestore.TransformChange, a.bufferSize),
	}
	a.subscribers[s.id] = s
	return s
}

func (a *worldStateAggregator) unregister(s *subscriber) {
	a.subMu.Lock()
	defer a.subMu.Unlock()
	if _, ok := a.subscribers[s.id]; !ok {
		return
	}
	delete(a.subscribers, s.id)
	close(s.ch)
}

// publish sends change to every live subscriber. Non-blocking per subscriber:
// a full buffer drops the event for that subscriber only.
func (a *worldStateAggregator) publish(change worldstatestore.TransformChange) {
	a.subMu.Lock()
	defer a.subMu.Unlock()
	for _, s := range a.subscribers {
		select {
		case s.ch <- change:
		default:
			a.logger.Warnw("dropped change for subscriber; buffer full",
				"subscriber_id", s.id,
				"uuid", string(change.Transform.Uuid),
				"change_type", change.ChangeType.String(),
			)
		}
	}
}

func (a *worldStateAggregator) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	if v, ok := cmd["set_transform"]; ok {
		return a.doSetTransform(v)
	}
	if v, ok := cmd["remove_transform"]; ok {
		return a.doRemoveTransform(v)
	}
	if _, ok := cmd["list_transforms"]; ok {
		return a.doListTransforms()
	}
	return nil, fmt.Errorf("unknown command, expected 'set_transform', 'remove_transform', or 'list_transforms'")
}

func (a *worldStateAggregator) doSetTransform(v interface{}) (map[string]interface{}, error) {
	args, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("set_transform: value must be an object, got %T", v)
	}
	uuid, err := extractUUID(args, "set_transform")
	if err != nil {
		return nil, err
	}
	t, err := parseTransform(args)
	if err != nil {
		return nil, fmt.Errorf("set_transform: %w", err)
	}
	t.Uuid = []byte(uuid)
	existed := a.store.set(uuid, t)

	changeType := pb.TransformChangeType_TRANSFORM_CHANGE_TYPE_ADDED
	if existed {
		changeType = pb.TransformChangeType_TRANSFORM_CHANGE_TYPE_UPDATED
	}
	a.publish(worldstatestore.TransformChange{ChangeType: changeType, Transform: t})

	return map[string]interface{}{"success": true, "uuid": uuid}, nil
}

func (a *worldStateAggregator) doRemoveTransform(v interface{}) (map[string]interface{}, error) {
	args, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("remove_transform: value must be an object, got %T", v)
	}
	uuid, err := extractUUID(args, "remove_transform")
	if err != nil {
		return nil, err
	}
	if prev := a.store.remove(uuid); prev != nil {
		a.publish(worldstatestore.TransformChange{
			ChangeType: pb.TransformChangeType_TRANSFORM_CHANGE_TYPE_REMOVED,
			Transform:  prev,
		})
	}
	return map[string]interface{}{"success": true, "uuid": uuid}, nil
}

func (a *worldStateAggregator) doListTransforms() (map[string]interface{}, error) {
	snap := a.store.snapshot()
	out := make([]map[string]interface{}, 0, len(snap))
	for _, t := range snap {
		raw, err := protojson.Marshal(t)
		if err != nil {
			return nil, fmt.Errorf("list_transforms: %w", err)
		}
		var m map[string]interface{}
		if err := json.Unmarshal(raw, &m); err != nil {
			return nil, fmt.Errorf("list_transforms: %w", err)
		}
		out = append(out, m)
	}
	return map[string]interface{}{"transforms": out}, nil
}

func extractUUID(args map[string]interface{}, cmd string) (string, error) {
	raw, ok := args["uuid"]
	if !ok {
		return "", fmt.Errorf("%s: missing required field \"uuid\"", cmd)
	}
	uuid, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("%s: \"uuid\" must be a string, got %T", cmd, raw)
	}
	if uuid == "" {
		return "", errors.New(cmd + ": \"uuid\" must be non-empty")
	}
	return uuid, nil
}

func parseTransform(args map[string]interface{}) (*commonpb.Transform, error) {
	payload := make(map[string]interface{}, len(args))
	for k, v := range args {
		if k == "uuid" {
			continue
		}
		payload[k] = v
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	var t commonpb.Transform
	if err := protojson.Unmarshal(raw, &t); err != nil {
		return nil, fmt.Errorf("invalid transform fields: %w", err)
	}
	return &t, nil
}

type transformStore struct {
	mu         sync.RWMutex
	transforms map[string]*commonpb.Transform
}

func newTransformStore() *transformStore {
	return &transformStore{transforms: make(map[string]*commonpb.Transform)}
}

func (s *transformStore) list() [][]byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([][]byte, 0, len(s.transforms))
	for k := range s.transforms {
		out = append(out, []byte(k))
	}
	return out
}

func (s *transformStore) get(uuid string) *commonpb.Transform {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.transforms[uuid]
}

func (s *transformStore) set(uuid string, t *commonpb.Transform) (existed bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, existed = s.transforms[uuid]
	s.transforms[uuid] = t
	return existed
}

// remove deletes uuid and returns the previous value, or nil if it wasn't set.
func (s *transformStore) remove(uuid string) *commonpb.Transform {
	s.mu.Lock()
	defer s.mu.Unlock()
	prev, ok := s.transforms[uuid]
	if !ok {
		return nil
	}
	delete(s.transforms, uuid)
	return prev
}

func (s *transformStore) snapshot() []*commonpb.Transform {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*commonpb.Transform, 0, len(s.transforms))
	for _, t := range s.transforms {
		out = append(out, t)
	}
	return out
}
