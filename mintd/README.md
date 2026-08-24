# Parafa mintd

**mintd** is the server institutions run. Wallets talk to it to get notes issued and redeemed.

## How it works

**mintd** signs blinded messages and stores spent notes.

It doesn't manage accounts/identities, or funds. All of that is done by the institution.

## Status

Early development. It runs but it can't issue or redeem anything yet.

Working:

- 2 HTTP servers, public and admin
- Configuration via environment variables and flags
- Warning if the admin API is not on a local address
- Graceful shutdown

Not built yet:

- Seed generation, encryption and loading
- Key derivation
- Signing
- Every endpoint except `/ping`

## Servers

**Public API.** Wallets talk to this. In production it's open to the internet. (it is on a loopback address by default, so you will need a reverse proxy, for example using Caddy/Nginx servers)

**Admin API.** For the institution's own systems, for payment confirmations/withdrawals. (it is on a loopback address by default, if you change the host, you will get a warning)

## Configuration

Flags overwrite environment variables, which overwrite the defaults. You don't need to rebuild if you have your own settings.

| Setting | Environment variable | Flag | Default |
|---|---|---|---|
| Seed file | `PARAFA_SEED_PATH` | `--seed-path` | `/var/lib/parafa/seed` |
| Public API | `PARAFA_API_ADDRESS` | `--api-addr` | `127.0.0.1:8080` |
| Admin API | `PARAFA_ADMIN_ADDRESS` | `--admin-addr` | `127.0.0.1:8081` |

The seed path must include the filename.

Run `mintd --help` for the full list.

## Run it

clone repo, then:

```sh
go build -o bin/mintd ./mintd
./bin/mintd
```
OR using make:

```sh
make build mintd
./bin/mintd
```