# goRPG

**Dungeon of Shadows** — a small text adventure RPG written in Go.

Pick a class, explore rooms, fight enemies, gather loot and try to beat the boss.

## Requirements

- Go 1.26+

## Run

```bash
go run .
```

Or build the binary:

```bash
make build
./goRpg
```

## Gameplay

At startup you choose a name and one of three classes:

| Class   | HP  | ATK |
|---------|-----|-----|
| Warrior | 100 | 15  |
| Mage    | 70  | 25  |
| Thief   | 85  | 18  |

### Commands

| Command                   | Alias           | Description                |
| ------------------------- | --------------- | -------------------------- |
| `go <dir>` / `move <dir>` | `n` `s` `e` `w` | Move in a direction        |
| `look`                    | `l`             | Describe the current room  |
| `inventory`               | `i`, `inv`      | Show items                 |
| `status`                  |                 | Show character sheet       |
| `help`                    | `?`             | Show available commands    |
| `quit`                    | `q`             | Exit the game              |

## Project structure

```text
.
├── game/         # Game logic (player, enemies, world, combat)
├── main.go       # Entry point
├── Makefile      # build / test / cover targets
└── go.mod
```

## Testing

```bash
make test     # run all tests with -v
make cover    # generate coverage.out + coverage.html
```
