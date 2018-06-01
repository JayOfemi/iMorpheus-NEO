parse signature scripts
===

Currently , some blockchain system,like bitcoin ethere,or other network based on them need a credit system to be built.
Credit provement needs a lots of things to be done in which a most important is proving your identity (ID) 
This part will help you validate a network node who communicating with you by parsing signature scripts.

try it
===

You need a completing bitcoin client. Be patient, sync transaction data may need some time to check every transaction up to now.

$./bitcoin-qt -server -rpcuser=rpc -rpcpassword=rpc -txindex=1 -printtoconsole -deprecatedrpc=accounts
$ export GOPATH=/yourdir
$ go build signature

you can choose your like IDE (eg liteIDE) 
