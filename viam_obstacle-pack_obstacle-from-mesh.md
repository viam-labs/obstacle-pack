# Model viam:obstacle-pack:obstacle-from-mesh

A generic component that loads a mesh from a [PLY](https://en.wikipedia.org/wiki/PLY_(file_format)) file at startup and exposes it via `Geometries`.

> **Note:** This model is registered under the `rdk:component:generic` API. The generic API gRPC server in RDK does **not** implement `GetGeometries`, so the mesh will not be picked up by the motion planner or frame system automatically. If you want the mesh to act as a collision obstacle for motion planning, use [`viam:obstacle-pack:obstacle-from-mesh-gripper`](viam_obstacle-pack_obstacle-from-mesh-gripper.md) instead. Use this generic variant when you only need to read the geometry from application code (e.g. via the SDK's `GetGeometries` client call against the gripper API, or for introspection purposes).

## Configuration

```json
{
  "mesh_path": <string>
}
```

### Attributes

| Name        | Type   | Inclusion | Description                                                                                                                       |
|-------------|--------|-----------|-----------------------------------------------------------------------------------------------------------------------------------|
| `mesh_path` | string | Required  | Filesystem path to a `.ply` mesh file. Loaded once at construction time. Must end in `.ply` (case-insensitive). |

### Example Configuration

```json
{
  "name": "my-obstacle",
  "api": "rdk:component:generic",
  "model": "viam:obstacle-pack:obstacle-from-mesh",
  "attributes": {
    "mesh_path": "/home/viam/obstacles/table.ply"
  }
}
```

## DoCommand

This model does not implement `DoCommand`; calling it returns `ErrDoUnimplemented`.
