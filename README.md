# EXP

by [orpheus](github.com/orpheus)

A leveling and reward system for real life skills to incentivize productive behavior.
Based on the Old School RuneScape leveling system, one of the most iconic and proven
game mechanics, users will log or add their time to a skill and level up, unlocking
rewards, titles, and achievements. A questing system will be implemented to gameify 
goal setting and achievement. And blockchain integration will allow for NFT rewards.

### Start

Make sure `postgres` is installed on your local system and running with a database
called `exp`. Swap the environment variables with ones for your system.

With `go` installed, run:

```
DB_USER=postgres \
DB_PASS=password \
DB_NAME=exp \
DB_HOST=localhost \
DB_PORT=5432 \
go run main.go
```

Service will be listening on port `8080`