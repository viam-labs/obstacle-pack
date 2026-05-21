# Module obstacle-pack

A collection of Viam components that publish static mesh obstacles to a machine. Each model loads a PLY mesh from a configured file path and exposes it via `Geometries`.

## Models

| Model                                                                               | API                     | Use when                                                                                   |
| ----------------------------------------------------------------------------------- | ----------------------- | ------------------------------------------------------------------------------------------ |
| [`viam:obstacle-pack:obstacle-from-mesh`](viam_obstacle-pack_obstacle-from-mesh.md) | `rdk:component:generic` | A generic component that reads a mesh and renders a geometry from it through GetGeometries |
