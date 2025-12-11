# Chargetop 🔋

A high-performance, single-binary Terminal User Interface (TUI) application for monitoring your macOS battery status. Built with Go and Bubbletea.

![Battery Monitor Screenshot](https://via.placeholder.com/800x400?text=Battery+Monitor+Go+Screenshot)

## Features

- **Blazing Fast**: Native binary, instant startup, low memory footprint.
- **Real-time Monitoring**: Live updates of battery percentage and charging status.
- **Detailed Stats**: Cycle count, battery health, maximum capacity, wattage, and serial number.
- **Beautiful UI**: Modern, responsive terminal interface powered by Bubbletea.

## Installation

### Method 1: From Source (Go Required)

1.  Clone the repository:
    ```bash
    git clone https://github.com/aezizhu/chargetop.git
    cd chargetop
    ```

2.  Build the binary:
    ```bash
    go build -o chargetop
    ```

3.  Run:
    ```bash
    ./chargetop
    ```

### Method 2: Homebrew (Coming Soon)

```bash
brew tap aezizhu/chargetop
brew install chargetop
```

## Usage

- **Run**: `./chargetop`
- **Quit**: Press `q` or `Ctrl+C`
- **Refresh**: Press `r`

## Requirements

- macOS (tested on Sonoma/Ventura)
- Go 1.19+ (for building from source)
