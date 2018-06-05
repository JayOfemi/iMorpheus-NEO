
blockChainMorp
===

The demo is to create the accounts and each account includes a private key with the memorizing words, a public key and an address. These data is saved in disk, then one can use it to send coins to another one, to get balance of it's account, and to mine some new coins in the internal test net.

## Requirements

- Ubuntu 16.04

- Go version >= 1.10.

- golang.org/x/crypto/ripemd160

- github.com/visualfc/goqt

- libsecp256k1

```bash
go get github.com/visualfc/goqt
cd $GOPATH/src/github.com/visualfc/goqt/qtdrv
qmake "CONFIG+=release" && make
cd $GOPATH/src/github.com/visualfc/goqt/tools/rcc
qmake "CONFIG+=release" && make
cd $GOPATH/src/github.com/visualfc/goqt/ui
go install -v
sudo cp $GOPATH/src/github.com/visualfc/goqt/bin/libqtdrv.ui.so.1.0.0 /usr/lib/libqtdrv.ui.so.1
```

## Usage
When first run this demo, it will create a account, and exit normally. Then use command below to test each module.
- createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS
- createwallet - Generates a new key-pair and saves it into the wallet file
- getbalance -address ADDRESS - Get balance of ADDRESS
- listaddresses - Lists all addresses from the wallet file
- printchain - Print all the blocks of the blockchain
- reindexutxo - Rebuilds the UTXO set
- send -from FROM -to TO -amount AMOUNT -mine - Send AMOUNT of coins from FROM address to TO. Mine on the same node, when -mine is set
- startnode - Start a node to sync data without mining enabled
- startnode -miner ADDRESS - Start a node with mining enabled
- gui open the demo with GUI mode
