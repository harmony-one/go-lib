# Harmony: go-lib

go-lib is a library used to interact with Harmony's RPC layer as well as adding a lot of utility functions used by various frameworks, e.g. harmony-tf, harmony-tests and harmony-stress.

While go-sdk is an actual program/CLI this library is solely designed to be used/referenced by other tools and applications.

It tries to use go-sdk as much as possible, but given go-sdk's heavy reliance on CLI/Cobra, go-lib implements a few workarounds to enable certain RPC access and functionality from outside the scope of go-sdk.

go-lib also provides extra layers of data marshalling/unmarshalling, logic and other functionality.

# Build

```
go build ./...
```

# Usage & Examples

TODO!
