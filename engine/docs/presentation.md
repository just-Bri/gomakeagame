# Presentation helpers (UI, maps, camera, font)

Last updated: 2026-09-20

Gomag keeps **sim** (ecs/grid/pathing/arena) separate from **presentation**.
These packages are the shared draw/HUD layer for TD and RPG titles.

| Package | Role |
|---------|------|
| `font` | Default Pure-Go TTF (`gofont`) + `Apply` for gg |
| `ui` | Immediate-mode `Panel` / `Bar` / `Label` / `Button` / `Hit` |
| `tilemap` | Integer layers → ortho rectangles or iso diamonds; `FromWalkable` |
| `camera` | `Follow2D` / `FollowXZ` (feeds `project.Iso`) |
| `debug` | FPS tracker + line overlay |
| `project` | Soft-3D iso math + `DrawBox` / `DrawSorted` |
| `move` | `AABB` / `AABBXZ` + `ClampPoint*` |

## Engine vs game

| Concern | Engine | Game |
|---------|--------|------|
| Panel/bar chrome | `ui` | Layout, icons, quest text |
| Font face loading | `font` | Copy / dialogue content |
| Paint walkability | `tilemap.FromWalkable` | Tile art atlases later |
| Camera lerp | `camera` | Input that picks the target |
| Soft-3D boxes | `project.DrawBox` | Character silhouettes |
| Zone bounds | `move.AABBXZ` | Concrete min/max constants |

## Smoke

```bash
mise run hello_ui
```
