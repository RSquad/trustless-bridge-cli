# Trustless Bridge CLI

A set of command-line tools for interacting with a **trustless bridge** in the TON blockchain. This project includes utilities for fetching, processing, and verifying blockchain data.

## Features

- **Fetch Block**: Retrieves a block from the blockchain.
- **Prune Block**: Removes unnecessary data from a block.
- **Block Proof**: Generates a proof from one block to another.
- **Block Signatures**: Extracts block signatures.
- **Transaction Proof**: Constructs a proof for a transaction.
- **Deploy Contracts**: Deploy contracts using the `deploy` command.
- **Check Transaction**: Sends a `check_transaction` message to verify transactions.

## Configuration

By default, the utilities look for a configuration file named `.trustless-bridge-cli.yaml` in the user's home directory. You can also specify a custom configuration file path using the global `--config` flag.

There is a sample `.trustless-bridge-cli.yml` file in the root of the repository that you can use as a reference.

### Configuration Keys

- **`wallet_mnemonic`**: This key is used to specify the mnemonic phrase for the wallet.
- **`wallet_version`**: This key indicates the version of the wallet being used, which include: v1r1, v1r2, v1r3, v2r1, v2r2, v3r1, v3r2, v3, v4r1, v4r2, v5r1beta, v5r1final.
- **`lite_client_code`**: This key is used to specify the code for the lite client, which is necessary for deploying the lite client contract.
- **`tx_checker_code`**: This key is used to specify the code for the transaction checker, which is necessary for deploying the transaction checker contract.

These keys are required for executing the `deploy`, `run`, and `get` commands. If you don't use such commands, you can leave them empty.

## Installation

Make sure you have Go installed (version 1.23.1 or later).

1. **Clone the Repository**

   ```bash
   git clone https://github.com/RSquad/trustless-bridge-cli.git
   cd trustless-bridge-cli
   ```

2. **Build the Project**

   ```bash
   go build -o trustless-bridge-cli main.go
   ```

   This command compiles the project into an executable named `trustless-bridge-cli` in the current directory.

3. **Install into Go Workspace (Optional)**

   To install the CLI into your `$GOPATH/bin`, run:

   ```bash
   go install
   ```

   Once installed, you can call the CLI from any terminal session without specifying the path.

## Usage

### Network Selection

You can specify the network using the global `--network` flag, which can be either `testnet` or `fastnet`. The default is `testnet`.

### Running the CLI

After building or installing, you can run the utility:

```bash
trustless-bridge-cli block fetch -s <seqno> -w <workchain> -f <output-format> --network <network>
```

Alternatively, you can run it directly with `go run` (no need to build beforehand):

```bash
go run main.go block fetch -s <seqno> -w <workchain> -f <output-format>
```

### Saving Output to a File

Use standard shell redirection to save output to a file. For example:

```bash
go run main.go block fetch -s 27450812 -f bin > block.boc
```

This fetches the block with sequence number 27450812 in binary format and writes it to `block.boc`.

### Help Flags

To view help for any command, use `--help`. For example:

```bash
trustless-bridge-cli block --help
```

## Main Scripts

### Deploy All

```bash
go run main.go deploy all -s 706883 --network testnet --config ./.trustless-bridge-cli.yaml
```

This command will fetch the nearest previous key block from the specified sequence number `706883` from **fastnet**, mark it as trusted, and deploy the **LiteClient** and **TxChecker** in the **testnet** with the trusted block.

**Note:** The command fetches data from **fastnet** and deploys it to **testnet** if the `--network` flag is specified as **testnet** and vice versa.

### Send Check Transaction

```bash
go run main.go send check-tx -a EQCzBNUbnja6DRzZYwPj6HXS2IwHE4Oz9zYpun9MxXNmsHJN -t 0908bfb9eb41b3186e63ab043142a3c4d493bfbaa3013094f17a15d3575a3138 -s 706883 --network testnet --config ./.trustless-bridge-cli.yaml
```

This command will fetch block `706883` and transaction `0908bfb9eb41b3186e63ab043142a3c4d493bfbaa3013094f17a15d3575a3138` from **fastnet**, build proof, construct the message body `check_transaction#91d555f7 transaction:^Cell proof:^Cell current_block:^Cell = InternalMsgBody;`, and send it to the **testnet** to the address `EQCzBNUbnja6DRzZYwPj6HXS2IwHE4Oz9zYpun9MxXNmsHJN`.

**Note:** The command fetches data from **fastnet** and sends it to **testnet** if the `--network` flag is specified as **testnet** and vice versa.
