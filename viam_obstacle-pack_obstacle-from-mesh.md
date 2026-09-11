# Model viam:obstacle-pack:obstacle-from-mesh

A generic component that loads a mesh from a [PLY](https://en.wikipedia.org/wiki/PLY_(file_format)) file at startup and exposes it via `Geometries`.

> **Note:** This model is registered under the `rdk:component:generic` API. The generic API
> serves `GetGeometries`, and the frame system will pick the mesh up for motion planning as
> long as the component is configured with a `frame` that has a `parent` and does not itself
> declare a `geometry`.

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
