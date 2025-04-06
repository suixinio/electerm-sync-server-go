# Electerm sync server go

A simple electerm data sync server with go.

## Use

Requires go 1.16+

```bash
git clone git@github.com:electerm/electerm-sync-server-go.git
cd electerm-sync-server-go

# Install dependencies
go mod download

# create env file, then edit .env
cp sample.env .env

# Run in development mode
go run src/main.go

# would show something like
# server running at http://127.0.0.1:7837

# in electerm sync settings, set custom sync server with:
# server url: http://127.0.0.1:7837
# Then you can use http://127.0.0.1:7837/api/sync as API Url in electerm custom sync

# JWT_SECRET: your JWT_SECRET in .env
# JWT_USER_NAME: one JWT_USER in .env
```

## Build and Run in production

For Unix-like systems (Linux/macOS):

```bash
# Run the build script
./bin/build.sh

# Run the server (after configuring .env)
# For macOS:
GIN_MODE=release ./bin/electerm-sync-server-mac

# For Linux:
GIN_MODE=release ./bin/electerm-sync-server-linux
```

## Test

```bash
bin/test.sh
```

## Write your own data store

Just take [src/store/filestore.go](src/store/filestore.go) as an example, write your own read/write method

## Sync server in other languages

- [electerm-sync-server-kotlin](https://github.com/electerm/electerm-sync-server-kotlin)
- [electerm-sync-server-vercel](https://github.com/electerm/electerm-sync-server-vercel)
- [electerm-sync-server-rust](https://github.com/electerm/electerm-sync-server-rust)
- [electerm-sync-server-cpp](https://github.com/electerm/electerm-sync-server-cpp)
- [electerm-sync-server-java](https://github.com/electerm/electerm-sync-server-java)
- [electerm-sync-server-node](https://github.com/electerm/electerm-sync-server-node)
- [electerm-sync-server-python](https://github.com/electerm/electerm-sync-server-python)
- [electerm-sync-server-deno](https://github.com/electerm/electerm-sync-server-deno)
- [electerm-sync-server-go](https://github.com/electerm/electerm-sync-server-go)

## License

MIT
