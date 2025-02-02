# Trustless Bridge CLI

A set of command-line tools for interacting with a **trustless bridge** in the TON blockchain. This project includes utilities for fetching, processing, and verifying blockchain data.

## Features

- **Fetch Block**  
  Retrieves a block from the blockchain by its `seqno` and `workchain`. Supports output in `json`, `bin`, or `hex` format.
- **Prune Block**  
  Removes unnecessary data from a block to reduce its size. Supports output in `bin` or `hex` format.
- **Block Proof**  
  Generates a proof from one block to another, loaded from a specified file. Supports output in `json`, `bin`, or `hex` format. Currently works only with blocks from the masterchain.
- **Block Signatures**  
  Extracts and returns block signatures for a specified block in the masterchain. Supports output in `json`, `bin`, or `hex` format.
- **Transaction Proof**  
  Constructs a proof for a transaction contained within a specified block. Supports output in `hex` or `bin` format.
- **Deploy Contracts**  
  Deploy contracts using the `deploy` command.

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

Where:

- `<seqno>` is the block's sequence number.
- `<workchain>` is the workchain ID.
- `<output-format>` can be `json`, `bin`, or `hex`.
- `<network>` is the network selection, which can be `testnet` or `fastnet`. The default is `testnet`.

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
