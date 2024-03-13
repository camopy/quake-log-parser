# quake-log-parser

This is a command line tool that parses a Quake 3 Arena log file and outputs a report of the games kills and players.

> Note: &lt;world&gt; kills or suicides are counted as a negative kill score for the involved player.

##  1. Usage

```bash
$ quake-log-parser <log-file>
```

### 1.1. Example with output file flag

```bash
$ quake-log-parser -o <output-file> <log-file> 
```

#### Example output

```json
"game-10": {
  "total_kills": 60,
  "players": [
    "Oootsimo",
    "Dono da Bola",
    "Zeh",
    "Chessus",
    "Mal",
    "Assasinu Credi",
    "Isgalamido"
  ],
  "kills": {
    "Assasinu Credi": 3,
    "Chessus": 5,
    "Dono da Bola": 3,
    "Isgalamido": 6,
    "Mal": 1,
    "Oootsimo": -1,
    "Zeh": 7
  },
  "kills_by_means": {
    "MOD_BFG": 2,
    "MOD_BFG_SPLASH": 2,
    "MOD_CRUSH": 1,
    "MOD_MACHINEGUN": 1,
    "MOD_RAILGUN": 7,
    "MOD_ROCKET": 4,
    "MOD_ROCKET_SPLASH": 1,
    "MOD_TELEFRAG": 25,
    "MOD_TRIGGER_HURT": 17
  }
}
```