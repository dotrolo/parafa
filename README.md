# Parafa

**Parafa** is a collection of programs and packages written in Go that make up a lightweight and privacy-preserving digital cash system. It is designed to be universal and not limited to any currency or platform. The same binary can be used for any fiat currency, cryptocurrency or anything else, because the mintd service doesn't handle funds. Debiting/crediting has to be implemented by the operator.

### What is it for?

When you pay by card or transfer, the bank logs both the sender and the receiver. Cash doesn't work that way, but cash only moves in person.

Parafa lets a licensed operator issue digital notes that behave like cash. The operator knows you withdrew, and knows someone redeemed, but cannot connect the two.

Both ends stay identified. The operator does full KYC and can answer a lawful request about either. It works the same way a bank hands out cash; the only thing they don't see is how you spend it.

### Example use

- You want to withdraw cash from your bank that runs a parafa mintd service.
- Using a wallet, you request a withdrawal from the bank.
- Your bank debits your account and provides a signature that your wallet turns into a valid note.
- You can also give the note to your friends/family if you want to (through encrypted channels).
- They can swap the note to a new one with the same value so the previous owner's note becomes invalid.
- To redeem it, all they have to do is send a redemption request using their wallet. Bank then credits their account accordingly.

A note is just a piece of data that by design isn't linked to your identity. Works exactly like physical cash. If someone can access your note(s), they can spend them.

### Components

[Mintd](./mintd/) - The server operators run. Wallets talk to it to get notes issued and redeemed.

[Demo](./demo/) - A straightforward demonstration of serial generation, blinding, signing, verification & spent check.

### Status

Early development. It is NOT ready for real funds. See [mintd](./mintd/) for what works and what doesn't.

### AI use

No AI-generated code is in this repository. An LLM was used for learning and research about design choices. It contributed substantially to the wording of the README files. Earlier AI-assisted test files (*_test.go) were removed and are not part of the codebase.

### License

AGPL-3.0. See [LICENSE](./LICENSE).