# tmp2p

A simple CLI to validate tendermint/cometbft p2p addresses, both for reachability and correct Node ID (ed25519 authentication)

### Install

```bash
git clone https://github.com/strangelove-ventures/tmp2p.git
cd tmp2p
make install
```

### Usage

``` bash
$ tmp2p validate -h
Validate list of peers, optionally with limit

Usage:
  tmp2p validate [peers] [limit]

Aliases:
  validate, v

Examples:
$ tmp2p validate 17bfb555c37b79e89af31342f4e068bf4f93e144@65.108.137.39:26656,efa6e21632ca4c7070c28fb244d9079a92dce67d@65.21.134.202:26616
$ tmp2p v 17bfb555c37b79e89af31342f4e068bf4f93e144@65.108.137.39:26656,efa6e21632ca4c7070c28fb244d9079a92dce67d@65.21.134.202:26616 10
```

### Build static bins

```
$ make build-static
building tmp2p amd64 static binary...
building tmp2p arm64 static binary...
$ ls build
tmp2p-amd64  tmp2p-arm64
```
