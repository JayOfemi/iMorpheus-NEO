
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
