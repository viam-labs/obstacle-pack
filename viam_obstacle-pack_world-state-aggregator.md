# Model viam:obstacle-pack:world-state-aggregator

An in-memory aggregator implementing `rdk:service:world_state_store`. It collects runtime-attached
geometries from many producers into a single world state store. Producers push transforms via
`DoCommand`; consumers either poll (`ListUUIDs` + `GetTransform`) or subscribe
(`StreamTransformChanges`).

State is in-memory only and not persisted across restarts. Producers are the source of truth — they
republish on their own startup.

## Configuration

```json
{
  "subscriber_buffer_size": <int>
}
```

### Attributes

| Name                     | Type | Inclusion | Description                                                                                                                                              |
|--------------------------|------|-----------|----------------------------------------------------------------------------------------------------------------------------------------------------------|
| `subscriber_buffer_size` | int  | Optional  | Per-subscriber channel buffer. Default `64`. When a subscriber's buffer fills, changes are dropped and logged for that subscriber only; producers are never blocked. |

### Example Configuration

```json
{
  "name": "world-state",
  "api": "rdk:service:world_state_store",
  "model": "viam:obstacle-pack:world-state-aggregator",
  "attributes": {
    "subscriber_buffer_size": 64
  }
}
```

## Read API

Standard `rdk:service:world_state_store` methods:

- `ListUUIDs(ctx, extra)` — snapshot of all current transform UUIDs.
- `GetTransform(ctx, uuid, extra)` — current value for one UUID.
- `StreamTransformChanges(ctx, extra)` — subscribe to add/update/remove events. No replay of existing
  state; consumers call `ListUUIDs` + `GetTransform` first to bootstrap.

## DoCommand (Write API)

| Command            | Args                                                                           | Effect                                                                                     |
|--------------------|--------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------|
| `set_transform`    | `{uuid, reference_frame, pose_in_observer_frame, physical_object?, metadata?}` | Upsert. Emits `ADDED` if new, `UPDATED` if existing.                                       |
| `remove_transform` | `{uuid}`                                                                       | Delete. Emits `REMOVED` if the UUID existed. Unknown UUID is a silent no-op (no event).    |
| `list_transforms`  | `{}`                                                                           | Snapshot of all transforms. For CLI debug.                                                 |

UUIDs are operator-namespaced strings on the wire (e.g. `tool-changer-1/attached`), stored as bytes
internally. Transform fields use the protobuf JSON names of `common.v1.Transform`.

### Example

```json
{
  "set_transform": {
    "uuid": "tool-changer-1/attached",
    "reference_frame": "gripper-1",
    "pose_in_observer_frame": {
      "reference_frame": "gripper-1",
      "pose": { "x": 0, "y": 0, "z": 50, "o_z": 1 }
    },
    "physical_object": {
      "box": { "dims_mm": { "x": 40, "y": 40, "z": 100 } },
      "label": "tool"
    }
  }
}
```
