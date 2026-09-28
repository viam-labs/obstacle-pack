# Module obstacle-pack

A collection of Viam resources for publishing obstacles to a machine. The component models build static geometry from configuration and expose it via `Geometries`; the world state store service aggregates geometries pushed at runtime by other modules.

## Models

| Model                                                                               | API                     | Use when                                                                                   |
| ----------------------------------------------------------------------------------- | ----------------------- | ------------------------------------------------------------------------------------------ |
| [`viam:obstacle-pack:obstacle-from-mesh`](viam_obstacle-pack_obstacle-from-mesh.md) | `rdk:component:generic` | A generic component that reads a mesh and renders a geometry from it through GetGeometries |
| [`viam:obstacle-pack:obstacle`](viam_obstacle-pack_obstacle.md)                     | `rdk:component:generic` | You want to declare obstacle primitives (box/sphere/capsule) inline in the config          |
| [`viam:obstacle-pack:world-state-aggregator`](viam_obstacle-pack_world-state-aggregator.md) | `rdk:service:world_state_store` | Other modules need to push runtime-attached geometries into a single world state store that consumers can poll or subscribe to |
