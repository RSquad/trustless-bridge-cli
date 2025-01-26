# Trustless Bridge Utility Suite

This project is a suite of command-line utilities designed to interact with a trustless bridge. It provides tools for fetching and processing blockchain data.

## Features

- **Fetch Block**: Retrieve a block from the blockchain using its seqno and workchain ID. Supports output in JSON, binary, or hexadecimal formats.
- **Prune Block**: Remove unnecessary data from a blockchain block to reduce its size. Supports output in binary or hexadecimal formats.

## Configuration

The application uses a configuration file to set up necessary parameters. By default, it looks for a file named `.trustless-bridge-cli.yaml` in the user's home directory. The configuration file should include the following:

```yaml
ton_config_url: "https://ton-blockchain.github.io/testnet-global.config.json"
```

This URL is used to configure the TON client with the necessary blockchain settings.

## Documentation

The CLI provides built-in help documentation. You can access it by using the `--help` flag with any command. For example:

```bash
go run main.go block --help
```

## Usage

To run the utility, use the `go run` command:

```bash
go run main.go block fetch --s <seqno> --w <workchain> --f <output-format>
```

Replace `<seqno>`, `<workchain>`, and `<output-format>` with the desired block sequence number, workchain ID, and output format (`json`, `bin`, or `hex`), respectively.

To prune a block, use the following command:

```bash
go run main.go block prune --i <input-file> --f <output-format>
```

Replace `<input-file>` with the path to the block file you want to prune and `<output-format>` with the desired output format (`bin` or `hex`).

## Building

To build the project, use the following command:

```bash
go build -o trustless-bridge-cli main.go
```

This will create an executable named `trustless-bridge-cli` that you can run with the same options as above.
