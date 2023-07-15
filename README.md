# Ethereum Tracker
The Ethereum Tracker is pre-configured to work with the Sepolia testnet by default.


## Prerequisites

- Go 1.16 or higher

## Quick Start

```shell
# Clone this repo
git clone https://github.com/devson2561/eth-tracker.git && cd eth-tracker

# Install the dependencies
go mod download

# Add address to track (replace <address> with the actual address you want to track)
go run main.go add <address>

# start tracking
go run main.go start

```


## Configuration
To update the Ethereum node URL and database file path, please modify the corresponding constants in the `main.go` file as per your configuration.

```shell
  const (
    rpcUrl  = "wss://eth-node-url.com"  // Replace with your Ethereum node URL
    dbFile  = "./ethereum-watcher.db"  // Replace with your desired database file path
  )
```

## Note
If you have stopped the system and want to run it again, the system will prompt you whether you want to continue syncing from where you left off. You can choose to sync or not based on your preference.