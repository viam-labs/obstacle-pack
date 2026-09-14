# Model viam:obstacle-pack:obstacle

A generic component that turns a list of geometries declared in its configuration into a
static obstacle, exposed via `Geometries`. Geometries are parsed once at construction time.

Use this when you want to place simple primitives (boxes, spheres, capsules) into a scene.
To load an obstacle from a mesh file instead, see
[`viam:obstacle-pack:obstacle-from-mesh`](viam_obstacle-pack_obstacle-from-mesh.md).

> **Note on motion planning:** when this component is given a `frame` with a `parent`, the
> frame system asks it for its geometry and includes it in motion planning. A frame link
> holds only **one** geometry, so if `geometries` contains more than one entry only the
> first is picked up by the frame system (`viam-server` logs a warning). All configured
> geometries are always returned over `GetGeometries`. To place several primitives into a
> plan, configure one component per geometry.
>
> If the component's `frame` declares its own `geometry`, that takes precedence and the
> configured `geometries` are not consulted by the frame system.

## Configuration

```json
{
  "geometries": [
    {
      "type": <string>,
      "translation": { "x": <number>, "y": <number>, "z": <number> },
      "orientation": { ... },
      "label": <string>
    }
  ]
}
```

### Attributes

| Name         | Type  | Inclusion | Description                                                                                                                                               |
|--------------|-------|-----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------|
| `geometries` | array | Optional  | Geometry definitions in the standard Viam geometry format (`box`, `sphere`, `capsule`). Poses are relative to the component's own frame origin. Omitting it, or giving an empty list, is accepted and yields an obstacle with no geometries. |

Each entry uses the same schema as a `frame` geometry: `type` plus its dimensions
(`x`/`y`/`z` for a box, `r` for a sphere, `r`/`l` for a capsule), an optional
`translation` and `orientation`, and an optional `label`.

### Example Configuration

```json
{
  "name": "my-obstacle",
  "api": "rdk:component:generic",
  "model": "viam:obstacle-pack:obstacle",
  "attributes": {
    "geometries": [
      {
        "type": "box",
        "x": 1000,
        "y": 600,
        "z": 40,
        "translation": { "x": 0, "y": 0, "z": -20 },
        "label": "table-top"
      }
    ]
  },
  "frame": {
    "parent": "world",
    "translation": { "x": 400, "y": 0, "z": 0 }
  }
}
```

## DoCommand

This model does not implement `DoCommand`; calling it returns `ErrDoUnimplemented`.
