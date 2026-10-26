# Parafa

**Parafa** is a collection of programs and packages written in Go that make up a **lightweight**, privacy-preserving blind signing system, using BDHKE (Blind Diffie-Hellman Key Exchange) on [secp256k1](https://pkg.go.dev/github.com/decred/dcrd/dcrec/secp256k1/v4).

Parafa is built for a low attack surface, small binaries and open use.

### What is this?

An operator checks something about you, then gives you a voucher for it. When you redeem the voucher, the operator can see that it's valid, but not which person it gave it to.

Systems that already do this are each built for one use (e.g. digital cash or age verification). Parafa isn't. **mintd** signs vouchers, checks them, and records which ones were used. What a voucher means is up to the operator.

Right now a voucher only works with the operator that issued it. The goal is to let operators accept each other's vouchers.

#### Example use
- You want to withdraw cash from your bank that runs a Parafa mintd service.
- Using a wallet, you request a withdrawal from the bank.
- Your bank debits your account and provides a stamp that your wallet turns into a valid voucher (serial + stamp).
- You can also give the voucher to your friends/family if you want to (through encrypted channels).
- They can swap the voucher for a new one so the previous owner's voucher becomes invalid.
- To redeem it, all they have to do is send a redemption request using their wallet. The bank then credits their account accordingly.

A voucher is just a piece of data that by design isn't linked to your identity. It works like physical cash. If someone can access your voucher(s), they can spend them.

### Components

[mintd](./mintd/) - The server operators run. Wallets talk to it to get vouchers issued and redeemed.

[Demo](./demo/) - A straightforward demonstration of serial generation, blinding, signing, verification & spent check.

### Status

Early development. NOT ready for production. See [mintd](./mintd/) for what works and what doesn't.

### AI use

No AI-generated code is in this repository. AI was used for learning and research about design choices. It contributed to README files.

### License

AGPL-3.0. See [LICENSE](./LICENSE).