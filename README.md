# Module obstacle-pack

A collection of Viam components that publish static obstacles to a machine. Each model builds its geometry from configuration and exposes it via `Geometries`.

## Models

| Model                                                                               | API                     | Use when                                                                                   |
| ----------------------------------------------------------------------------------- | ----------------------- | ------------------------------------------------------------------------------------------ |
| [`viam:obstacle-pack:obstacle-from-mesh`](viam_obstacle-pack_obstacle-from-mesh.md) | `rdk:component:generic` | A generic component that reads a mesh and renders a geometry from it through GetGeometries |
| [`viam:obstacle-pack:obstacle`](viam_obstacle-pack_obstacle.md)                     | `rdk:component:generic` | You want to declare obstacle primitives (box/sphere/capsule) inline in the config          |
