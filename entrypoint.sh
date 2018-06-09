#!/bin/bash
# set -e

# if [[ "$1" == "bitcoin-cli" || "$1" == "bitcoin-tx" || "$1" == "bitcoind" || "$1" == "test_bitcoin" ]]; then
# 	mkdir -p "$BCH_DATA"

# 	if [[ ! -s "$BCH_DATA/bitcoin.conf" ]]; then
# 		cat <<-EOF > "$BCH_DATA/bitcoin.conf"
# 		printtoconsole=1
# 		rpcallowip=::/0
# 		rpcpassword=${BITCOIN_RPC_PASSWORD:-password}
# 		rpcuser=${BITCOIN_RPC_USER:-bitcoin}
# 		EOF
# 		chown bitcoin:bitcoin "$BCH_DATA/bitcoin.conf"
# 	fi

# 	# ensure correct ownership and linking of data directory
# 	# we do not update group ownership here, in case users want to mount
# 	# a host directory and still retain access to it
# 	chown -R bitcoin "$BCH_DATA"
# 	ln -sfn "$BCH_DATA" /home/bitcoin/.bitcoin
# 	chown -h bitcoin:bitcoin /home/bitcoin/.bitcoin

# 	exec gosu bitcoin "$@"
# fi

# echo $GOPATH
# echo $BCH_DATA
# echo $1
# exec "$@"

# go build ~/go/src/github.com/rohenaz/dtvcash/
# ./dtvcash web&
# ./dtvcash main-node
bitcoind -datadir="$BCH_DATA" &
# bitcoind -daemon -conf=bitcoin.conf -datadir="$BCH_DATA" &
# bitcoind &
dtvcash web &
dtvcash main-node