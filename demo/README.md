# Parafa demo (early development version)

This demonstration goes through a few simple steps. (server) generate keys from seed; (client) generate serial and blind it; (server) sign; (client) unblind; (server) redeem/verify.

### Run it

Clone the repo, then:

```sh
make build-demo
./bin/demo
```

Or without make:

```sh
go build -o bin/demo ./demo
./bin/demo
```

### What happens when you run this demo?

#### Setup

- `seed` - 32 random bytes, encrypted on disk
- `k = derive(seed, 1, "2026Q3")` - private key
- `K = kG` - public key

#### Client makes a serial


- `s` - 32 random bytes
- `S = HashToCurve(s)` - serial as a point
- `r` - random blinding factor

#### Client blinds

- `B = S + rG` - this gets sent to mintd

#### Mintd signs

```
C' = kB
   = k(S + rG)
   = kS + krG
```

#### Client unblinds

- `rK = rkG`
- `C = C' + mirror of rK = kS` - signed S

#### Mintd verifies it

- `HashToCurve(s) -> S` - compute from serial
- `note is valid if kS = C`