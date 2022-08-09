# go-peercoin-rpc
[![Go Report Card](https://goreportcard.com/badge/github.com/saeveritt/go-peercoin-rpc)](https://goreportcard.com/report/github.com/saeveritt/go-peercoin-rpc)
## Setup your environment variables

First, ensure that both `PPC_RPCUSER` and `PPC_RPCPASSWORD` ENV variables exist
```bash
echo "export PPC_RPCUSER='<rpcusername>'" >> ~/.profile
echo "export PPC_RPCPASSWORD='<rpcpassword>'" >> ~/.profile
```
Then, update shell
```bash
source ~/.profile
```
## Setup your config.json
For example,
```json
{
    "testnet": true,
    "host": "localhost",
    "port": 9904
}
```

