# Parafa mintd

**mintd** is the server operators run. Wallets talk to it to get notes issued and redeemed.

See the [main README](../README.md) for what Parafa is and how a note works.

## How it works

**mintd** signs blinded serials and stores spent notes.

It doesn't manage accounts/identities, or funds. All of that is done by the operator!

**mintd** has a secret ***seed*** which is stored in a file (by default at `/var/lib/parafa/seed`), every key (one key per denomination + epoch) derives from this seed, it MUST be backed up and secured by the operator!

The program asks for a passphrase, either to encrypt a new seed file or to decrypt an existing one. You can also feed it in through a pipe, from any source (e.g. `pass parafa/seed-passphrase | ./bin/mintd`).

Keep the passphrase somewhere safe and NOT ANYWHERE NEAR the encrypted seed file (use a vault, or a secrets mount for instance)

## Status

Early development. It runs but it can't issue or redeem anything over the API yet.

**DO NOT run mintd with real funds in its current state!**

Working:

- 2 HTTP servers, public and admin
- Configuration via environment variables, flags and config file
- Warning if the admin API is not on a local address
- Graceful shutdown
- Seed generation & loading
- Seed file and directory permission checks, refusing to start if too open
- Seed encryption
- Key derivation
- Sign and Verify

Not built yet:

- Tests
- DLEQ proofs: this ensures mintd cannot deanonymize users by signing their notes with unique keys without the user knowing. With DLEQ, users/wallets can do this verification themselves.
- API endpoints (public & admin)
- Side-channel hardening: point multiplication on curve uses NonConst operations meaning it takes more time for one operation to finish than another; this could potentially be exploited with our setup.
- Wallet library
- CLI wallet
- Mock operator running mintd with fake money

## Servers

**Public API.** Wallets talk to this. It is on a loopback address by default, you need a reverse proxy in front of it to make it accessible.

**Admin API.** For the operator's own systems, for payment confirmations/withdrawals. (it is on a loopback address by default, if you change the host, you will get a warning)

## Configuration

Flags overwrite environment variables, which overwrite the config file, which overwrites the defaults. You don't need to rebuild if you have your own settings.

| Setting | Environment variable | Flag | Default |
|---|---|---|---|
| Seed file | `PARAFA_SEED_PATH` | `--seed-path` | `/var/lib/parafa/seed` |
| Public API | `PARAFA_API_ADDRESS` | `--api-addr` | `127.0.0.1:8080` |
| Admin API | `PARAFA_ADMIN_ADDRESS` | `--admin-addr` | `127.0.0.1:8081` |
| Config file | `PARAFA_CONFIG` | `--config` | `/etc/parafa/parafa.conf` |

The seed path must include the filename.

### Config file
**It is not necessary to have a configuration file**, but if you don't want to always provide flags or environment variables for settings, you can write a config file and save it. It is recommended to create a file at `/etc/parafa/parafa.conf` as that is the path the program looks for by default without having to provide any paths to find it.

This config file uses `key = value` format, use `#` at the start of the line for comments.

Config keys use underscores (`seed_path`), while flags use hyphens (`--seed-path`).

Example config:
```ini
# /etc/parafa/parafa.conf
seed_path = /var/lib/parafa/seed
api_addr = 127.0.0.1:8080
admin_addr = 127.0.0.1:8081
```

Run `mintd --help` for the full list.

### **Notice**

**mintd** checks the permissions of the seed file and its parent directory, but securing the path above it is the operator's job!

## Run it (Linux)

Go 1.26.5. Clone the repo, then:

```sh
make build-mintd
./bin/mintd
```

Or without make:

```sh
go build -o bin/mintd ./mintd
./bin/mintd
```

If you don't pipe a passphrase in, mintd will ask for one. In production, pipe it from wherever you keep it.
