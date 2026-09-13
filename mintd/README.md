# Parafa mintd

**mintd** is the server operators run. Wallets talk to it to get notes issued and redeemed.

It's essentially a universal system for banks, exchanges or other licensed operators to convert your money into "notes" that you own and control just like physical cash.
Example use:
- You want to withdraw cash from your bank that runs a parafa mintd service.
- Using a wallet, you request a withdrawal from the bank.
- Your bank debits your account and provides a signature that your wallet turns into a valid note.
- You can also give the note to your friends/family if you want to (through encrypted channels).
- They can swap the note to a new one with the same value so the previous owner's note becomes invalid.
- To redeem it, all they have to do is send a redemption request using their wallet. Bank then credits their account accordingly.

A note is just a piece of data that by design isn't linked to your identity. Works exactly like physical cash. If someone can access your note(s), they can spend it.

The bank knows you withdrew, and knows someone redeemed, but cannot connect the two.

## How it works

**mintd** signs blinded serials and stores spent notes.

It doesn't manage accounts/identities, or funds. All of that is done by the operator!

**mintd** has a secret ***seed*** which is stored in a file (by default at `/var/lib/parafa/seed`), every key derives from this seed, it MUST be backed up and secured by the operator!

The program asks for a passphrase, either to encrypt a new seed file or to decrypt an existing one. You can also feed it in through a pipe, from any source (e.g. `pass parafa/seed-passphrase | ./bin/mintd`).

Keep the passphrase somewhere safe and NOT ANYWHERE NEAR the encrypted seed file (use a vault, or a secrets mount for instance)

## Status

Early development. It runs but it can't issue or redeem anything over the API yet.

**DO NOT run mintd with real funds in its current state!**

Working:

- 2 HTTP servers, public and admin
- Configuration via environment variables and flags
- Warning if the admin API is not on a local address
- Graceful shutdown
- Seed generation & loading
- Seed file and directory permission checks, refusing to start if too open
- Seed encryption
- Key derivation
- Sign and Verify
- Basic demo

Not built yet:
- DLEQ proofs: this ensures mintd cannot deanonymize users by signing their notes with unique keys without the user knowing. with DLEQ, users/wallets can do this verification themselves.
- API Endpoints (public & admin)
- Side-channel hardening: point multiplication on curve uses NonConst operations meaning it takes more time for one operation to finish than another; this could potentially be exploited with our setup.
- Wallet library
- CLI Wallet
- Mock operator running mintd with fake money.


## Servers

**Public API.** Wallets talk to this. It is on a loopback address by default, you need a reverse proxy in front of it to make it accessible.

**Admin API.** For the operator's own systems, for payment confirmations/withdrawals. (it is on a loopback address by default, if you change the host, you will get a warning)

## Configuration

Flags overwrite environment variables, which overwrite the defaults. You don't need to rebuild if you have your own settings.

| Setting | Environment variable | Flag | Default |
|---|---|---|---|
| Seed file | `PARAFA_SEED_PATH` | `--seed-path` | `/var/lib/parafa/seed` |
| Public API | `PARAFA_API_ADDRESS` | `--api-addr` | `127.0.0.1:8080` |
| Admin API | `PARAFA_ADMIN_ADDRESS` | `--admin-addr` | `127.0.0.1:8081` |

The seed path must include the filename.

Run `mintd --help` for the full list.

### **Notice**

**mintd** checks the permissions of the seed file and its parent directory, but securing the path above it is the operator's job!

## Run it (Linux)
Go 1.26.5

### Mintd:
clone repo, then:

```sh
go build -o bin/mintd ./mintd
./bin/mintd
```
OR using make:

```sh
make build-mintd
./bin/mintd
```

### Basic demo:
```sh
go build -o bin/demo ./demo
./bin/demo
```
OR using make:

```sh
make build-demo
./bin/demo
```

If you don't pipe a passphrase in, mintd will ask for one. In production, pipe it from wherever you keep it.

## AI use
No AI-generated code is in this repository. An LLM was used for learning and research about design choices. It contributed substantially to the wording of this README. Earlier AI-assisted test files (*_test.go) were removed and are not part of the codebase.

## License

AGPL-3.0. See [LICENSE](../LICENSE).