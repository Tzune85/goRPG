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
| `1` | Open buy menu |
| `2` | Open sell menu |
| `3` | Exit the shop |

After choosing buy or sell, a numbered item list appears (e.g. `[1] Health Potion = 5 Gold`).
The output will contain `[0] Back`.
Send the item number to buy/sell it, or `0` to go back to the main shop menu.

### Using Potions Outside Combat

You can use a health potion at any time during navigation by sending:

```
POST /api/action
{"action": "potion"}
```

Use it when your HP is low before entering a new room.

## Strategy to Win

To defeat the Ancient Dragon at the end of the dungeon you must be strong enough. Follow this loop:

1. **Farm enemies** — fight every enemy you find to gain XP and level up. Higher level = more damage and HP.
2. **Visit the shop** — spend your gold on better weapons and armor. Stronger gear makes all future fights easier.
3. **Heal between fights** — use potions (`potion` action) when navigating if HP is below 50%.
4. **Explore fully** — check all exits in every room before moving on. Unexplored rooms may contain the shop or shortcuts.
5. **Don't run early** — running from weak enemies wastes XP. Only run from bosses you cannot beat yet.

## Decision Logic

Parse the `output` field to decide what to do next:

- Contains `Actions: [1] attack` → you are in **combat**
- Contains `Actions: [1] Buy` → you are in a **shop**
- Contains `Exits:` → you are **navigating** — consider using a potion if HP < 50%, then choose a direction
- `done: true` → game over, stop the loop

## Spectators

Humans can watch your run live at:
```
/watch
```
