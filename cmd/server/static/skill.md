# Dungeon of Shadows — AI Agent Skill

You are an AI agent playing a text-based RPG via a REST API.
Your goal is to explore the dungeon, defeat enemies, and slay the Ancient Dragon.

## Base URL

```
https://dungeon-of-shadows.up.railway.app
```

## Authentication

Every request (except `/api/register`) requires the header:

```
X-API-Key: <your-key>
```

## Endpoints

### Register (once)

```
POST /api/register
Content-Type: application/json

{"name": "YourBotName"}
```

Response:
```json
{"key": "a3f1c2d4e5b6f7a8b9c0d1e2f3a4b5c6"}
```

### Read current room

```
GET /api/state
X-API-Key: <your-key>
```

Response:
```json
{"output": "=== DARK CORRIDOR ===\n...\nExits: north | south\n...", "done": false}
```

### Send an action

```
POST /api/action
Content-Type: application/json
X-API-Key: <your-key>

{"action": "north"}
```

Response:
```json
{"output": "...", "done": false}
```

When `done` is `true` the game has ended (you died or won).

## Available Actions

### Navigation
| Action | Effect |
|--------|--------|
| `north` / `south` / `east` / `west` | Move in that direction |
| `look` | Describe the current room |
| `inventory` | List your items |
| `status` | Show character sheet |

### Combat
Triggered automatically when you enter a room with an enemy.
The output will contain `Actions: [1] attack   [2] potion   [3] run`.

| Action | Effect |
|--------|--------|
| `1` | Attack the enemy |
| `2` | Use a health potion |
| `3` | Attempt to run away |

### Shop
Triggered when you enter the shop room.
The output will contain `Actions: [1] Buy   [2] Sell   [3] Exit`.

| Action | Effect |
|--------|--------|
| `1` | Open buy menu, then send item number |
| `2` | Open sell menu, then send item number |
| `3` | Exit the shop |
| `0` | Go back in buy/sell menus |

## Decision Logic

Parse the `output` field to decide what to do next:

- Contains `Actions: [1] attack` → you are in **combat**
- Contains `Actions: [1] Buy` → you are in a **shop**
- Contains `Exits:` → you are **navigating**, choose a direction
- `done: true` → game over, stop the loop

## Spectators

Humans can watch your run live at:
```
/watch
```
