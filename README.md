# Parafa

**Parafa** is a collection of programs and packages written in Go that make up a **lightweight**, privacy-preserving blind signing system, using BDHKE (Blind Diffie-Hellman Key Exchange) on [secp256k1](https://pkg.go.dev/github.com/decred/dcrd/dcrec/secp256k1/v4).

Parafa is built for a low attack surface, scalability, small binaries and open use.

### What is this?

An operator (e.g. a company, organization or institute) using this system can give out vouchers for a verified fact, e.g. "You have a valid ticket for the concert" or "You just withdrew 100 euros". They give vouchers once the fact has been verified.

When you redeem this voucher, e.g. at the gate for ticket inspection or when depositing cash, they cannot tell who they gave it to. They can only tell that this voucher is valid (because it has a valid signature).

It works because you give a blinded "message" to the operator. After they verify the fact (e.g. your account was charged), they sign it and give the signature back. You remove the blinding from the signature, and then you have your voucher. A voucher is just a pair of data: the "message" (serial) and the signature (stamp). You can use it as proof (only once per voucher) without revealing your identity, e.g. to get into the concert.

Read more about how blind signatures are used [here](https://en.wikipedia.org/wiki/Blind_signature) (note: we use a different protocol & number system than the examples shown on the wiki page).

#### Example use
- You want to withdraw cash from your bank that runs a Parafa mintd service.
- Using a wallet, you request a withdrawal from the bank.
- Your bank debits your account and provides a stamp that your wallet turns into a valid voucher (serial + stamp).
- You can also give the voucher to your friends/family if you want to (through encrypted channels).
- They can swap the voucher for a new one so the previous owner's voucher becomes invalid.
- To redeem it, all they have to do is send a redemption request using their wallet. The bank then credits their account accordingly.

If someone can access your voucher(s), they can spend them.

### Components

[mintd](./mintd/) - The server operators run. Wallets talk to it to get vouchers issued and redeemed.

[Demo](./demo/) - A straightforward demonstration of serial generation, blinding, signing, verification & spent check.

### Status

Early development. NOT ready for production. See [mintd](./mintd/) for what works and what doesn't.

### AI use

No AI-generated code is in this repository. AI (Claude, model: Opus 5 and 5.5) was used for learning, research about design choices and read-only code review. It contributed to the following files by structuring, rewording parts, and making terms consistent:
- /README.md
- /mintd/README.md
- /demo/README.md

### License

AGPL-3.0. See [LICENSE](./LICENSE).