# How to check ETH node

```

root@noxon-testnet ~ # su docker
docker@noxon-testnet:/root$ cd ~/geth
docker@noxon-testnet:~/geth$ docker-compose exec geth geth attach
Welcome to the Geth JavaScript console!

instance: Geth/v1.10.6-unstable-f0b1bdda-20210714/linux-amd64/go1.16.6
at block: 0 (Thu Jan 01 1970 00:00:00 GMT+0000 (UTC))
 datadir: /root/.ethereum
 modules: admin:1.0 debug:1.0 eth:1.0 ethash:1.0 miner:1.0 net:1.0 personal:1.0 rpc:1.0 txpool:1.0 web3:1.0

To exit, press ctrl-d
> eth.syncing
{
  currentBlock: 12856462,
  highestBlock: 12856581,
  knownStates: 632104722,
  pulledStates: 632057184,
  startingBlock: 6193069
}
>


```
