# ZorroDatage

Showcase application that demonstrates the backend integration of [Rockside](https://rockside.io) **gasless interactions with Ethereum**!

This repository contains the full code that runs the application at [https://timestamping-showcase.rockside.io](https://timestamping-showcase.rockside.io)

# How does it work?

Let's say you want to easily build a GO backend that performs Ethereum proof of existence on documents. 

## Project Layout Example

Your project layout might look a bit like this:

```sh
my_app
  ├── go.mod
  ├── go.sum
  ├── abi.go
  ├── main.go
  └── contract.sol
```

* `contract.sol` is your timestamping solidity contract.
* `main.go` contains your server logic
* `abi.go` contains all your generated contract bindings in GO. Handy!

To generate the `abi.go` file, use [abigen](https://github.com/ethereum/go-ethereum/tree/master/cmd):

1. Add `//go:generate abigen --sol mycontract.sol --pkg main --out abi.go` at the top of your `main.go` file 
2. Run `go generate` and the `abi.go` file will appear.

Also at some point you might want to deploy your contract on Ethereum. See below on how to deploy your contract easily with Rockside.

## Deploy your contract

There are various ways to deploy a smart contract. Here we will use Rockside and the Rockside CLI.

1. Install the [Rockside CLI](https://github.com/rocksideio/rockside-sdk-go/#command-line-interface-usage)
2. Get an API Key from Rockside
3. Create a Rockside identity: `rockside identities create` then list it with `rockside identities ls`
4. Deploy the contract with: `rockside deploy-contract`  (`cd` into the directory containing your contract)

Your contract is deployed on the Ethereum blockchain! 

*Note that to use Rospten instead add the flag `--testnet` each time you use the CLI.*

## Store fingerprint on the blockchain

From our backend we want to store fingerprint of incoming files in Ethereum.

```go
rocksideClient, err = rockside.NewClient(apiURL, apiKey, network)
...
transactor := rockside.NewTransactor(rocksideIdentityAddress, rocksideClient)
...
contract, err := NewContractTransactor(contractAddress, rocksideTransactor)
...

tx, err := contract.Register(rockside.TransactOpts(), fingerprint)
...
``` 

## Read fingerprints on the blockchain

From our backend we want to read from our contract, see if a fingerprint is stored.

```go
rocksideClient, err = rockside.NewClient(apiURL, apiKey, network)
...
contract, err := NewContractCaller(contractAddress, rocksideClient.RPCClient)
if err != nil {
...
}

timestamp, err := contract.Fingerprints(&bind.CallOpts{}, fingerprint))
```