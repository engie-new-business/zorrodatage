# ZorroDatage

Showcase application that demonstrates the backend integration of [Rockside](https://rockside.io) **gasless interactions with Ethereum**!

This repository contains the full code that runs the application at [https://timestamping-showcase.rockside.io](https://timestamping-showcase.rockside.io)

You can explore the code to understand how the Rockside integration is done into a GO backend.

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

1. Install `abigen` with `go install github.com/ethereum/go-ethereum/cmd/abigen` 
2. Add `//go:generate abigen --sol contract.sol --pkg main --out abi.go` at the top of your `main.go` file 
3. Run `go generate` and the `abi.go` file will be generated and contract bingings will be part of your project.

Also at some point you might want to deploy your contract on Ethereum. See below on how to deploy your contract easily with Rockside.

## Deploy your contract

There are various ways to deploy a smart contract. Here we will use Rockside and the Rockside CLI.

1. Install the [Rockside CLI](https://github.com/rocksideio/rockside-sdk-go/#command-line-interface-usage)
2. Get an API Key from Rockside
3. Create a Rockside identity: `rockside identities create` (then list it with `rockside identities ls`)
4. Deploy the contract with: `rockside deploy-contract`  (`cd` into the directory containing your contract)

Your contract is deployed on the Ethereum blockchain! 

*Note that to use Rospten instead add the flag `--testnet` each time you use the CLI.*

## Store fingerprints on the blockchain

Using the ABI GO bindings and the Rockside client and transactor we can send transactions to our contract without worrying about the nonce or gas.

To see how it is done, either look into this repository's code in `main.go` or examples in the [GO doc reference](https://pkg.go.dev/github.com/rocksideio/rockside-sdk-go?tab=doc#pkg-overview).

## Read fingerprints on the blockchain

From our backend we want to read from our deployed contract, to check if a fingerprint is stored.

Using the Rockside GO client we can easily obtain a JSON RPC client, then feed it into a `ContractCaller` to call read-only method on our contract.

To see how it is done, either look into this repository's code in `main.go` or examples in the [GO doc reference](https://pkg.go.dev/github.com/rocksideio/rockside-sdk-go?tab=doc#pkg-overview).