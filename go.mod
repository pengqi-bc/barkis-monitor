module github.com/barkisnet/monitor

go 1.13

require (
	github.com/btcsuite/btcd v0.0.0-20190523000118-16327141da8c
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/barkisnet/barkis v0.0.0-20200307124653-9882dc2ed461
	github.com/tendermint/tendermint v0.32.2
	golang.org/x/crypto v0.0.0-20190313024323-a1f597ede03a
)

replace github.com/barkisnet/barkis => github.com/barkisnet/barkis2.0 v0.0.0-20200307124653-9882dc2ed461

replace github.com/tendermint/iavl => github.com/barkisnet/iavl v0.12.4-barkis
