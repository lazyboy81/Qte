# Query Terminal (QTe)

Query Terminal is a Go terminal application for connecting to databases, writing SQL, and viewing query results in the terminal. It uses Bubble Tea v2, Bubbles, Huh, and Lip Gloss for its interface.

> **Development is paused.** This project will not be actively developed for now. I plan to return to it later, but there is no timeline for resuming development. The repository is left in an unfinished state.

## Current state

The project includes database connection forms, saved connection profiles, a SQL editor, and a result table. The current source tree is partway through a UI and connection-tab refactor: the root model renders the home and information panels, but screen switching and event forwarding to the connection and query screens are disabled. The full connection-to-query workflow is therefore not currently accessible from the application.

The sections below describe the existing components and their intended workflow. Multiple connection tabs are still unfinished.

## Features in the codebase

- PostgreSQL, SQLite, and Oracle database drivers.
- Connection forms with a database type, connection string, and profile name.
- Saved connection profiles with encrypted connection strings and an encryption key stored in the operating system's keyring.
- A multiline SQL editor with line numbers, displayed above a query result table.
- Tabular results for statements beginning with `SELECT` or `WITH`.
- Execution of other statements, including data changes and DDL, with affected-row and last-insert-ID values where the driver provides them.
- A 60-second timeout for query execution, loading indicators, and query error messages.

## Requirements

- Go 1.25.12 or newer, as specified in `go.mod`.
- CGO enabled and a working C compiler for the bundled SQLite and Oracle drivers.
- An interactive terminal.
- An accessible operating system keyring for saving and loading connection profiles.
- A reachable database server, or a local SQLite database file. Oracle connections also require the Oracle client libraries used by `godror`.

Docker is optional; the repository includes a PostgreSQL development service.

## Run and build

From the repository root:

```sh
go mod download
go run ./cmd/qte
```

Build a Windows executable:

```sh
go build -o build/qte.exe ./cmd/qte
```

For Linux or macOS, use an output name without `.exe`:

```sh
go build -o build/qte ./cmd/qte
```

The Makefile also provides `make run`, `make build`, and `make test`. Its default build target produces `build/qte.exe`. The `build-all` target lists Linux, Windows, and macOS builds, but the CGO drivers require an appropriate compiler and build environment for each target.

The application writes JSON logs to `qte.log` in the working directory.

## Connection workflow

Once screen routing is restored, the connection components are designed for this workflow:

1. Choose **New Connection** from the home menu.
2. Select the database type and enter its connection string in **DB URL**.
3. Enter a profile **Name**. All three fields are required.
4. After the connection succeeds, the profile is saved and the query screen opens.
5. Write SQL and press `Ctrl+E` to execute it.
6. Use **Existing Connection** from the home menu to reconnect with a saved profile.

Credentials belong in the driver's connection string; there are no separate username or password fields. Profile names identify saved connections, so saving different connection details under an existing name replaces that profile.

| Database | Driver | DB URL value |
| --- | --- | --- |
| PostgreSQL | `pgx` | `postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable` for the local development service below |
| SQLite | `sqlite3` | A file path such as `./example.db`, or `:memory:` for a temporary in-memory database |
| Oracle | `godror` | An Oracle connection string accepted by `godror` |

### Local PostgreSQL service

Start the service defined in `docker-compose.yml`:

```sh
docker compose up -d postgres
```

It runs PostgreSQL 16 Alpine on port `5432`, with database `mydb`, user `myuser`, and password `mypassword`. These are local development credentials. Data persists in the `postgres_data` Docker volume.

Stop the service with:

```sh
docker compose down
```

### Keyboard controls

These bindings are defined in the screen components; the current root-screen refactor prevents access to most of them.

| Key | Action |
| --- | --- |
| `Ctrl+C` | Quit the application |
| `Enter` | Select a home-menu item or saved connection |
| `?` | Request the information screen |
| `Ctrl+R` | Return from the information screen to the home screen |
| `Ctrl+E` | Submit SQL from the editor |
| `Up` / `Down` | Navigate result rows |
| `Ctrl+H` / `Ctrl+L` | Cycle the query model's existing tabs |

## Saved connections

Connection profiles are stored in `QTe/connections.json` under the directory returned by Go's `os.UserConfigDir()`:

| Platform | Default location |
| --- | --- |
| Windows | `%APPDATA%\QTe\connections.json` |
| macOS | `~/Library/Application Support/QTe/connections.json` |
| Linux | `$XDG_CONFIG_HOME/QTe/connections.json`, or `~/.config/QTe/connections.json` when unset |

Connection strings are encrypted with AES-256-GCM. The key is stored in the OS keyring under service `QTe`, entry `encryption-key`. Profile names and driver names remain readable in the JSON file. The saved-connection picker displays decrypted connection strings.

## Tests

```sh
go test ./...
```

The existing database tests use in-memory SQLite databases to check column and row extraction, table creation, inserts, and selection. They do not require the Docker service and do not cover the complete interactive workflow or live PostgreSQL and Oracle connections.

## Project layout

```text
cmd/qte/                 Application entry point and logging
internal/common/         Shared messages and screen identifiers
internal/connection/     Connection form, picker, establishment, and modal
internal/db/             Database drivers, connections, query execution, and tests
internal/home/           Home menu
internal/info/           Information screen
internal/query/          Query screen, editor, and result table
internal/root/           Root UI model and navigation
internal/store/          Connection storage, encryption, and keyring integration
```

## Development history and remaining work

Remaining work includes restoring the full UI workflow, completing multiple connection tabs, adding an error overlay and SQL formatting, and supporting semicolon-separated queries with separate result tabs. See [todo.md](todo.md) for the existing task list. These items are deferred while development is paused.

Query dispatch currently checks only whether trimmed SQL starts with `SELECT` or `WITH`; it is not a SQL parser. Other statements use `ExecContext`, and multiple statements are not split into individual queries or result sets.
