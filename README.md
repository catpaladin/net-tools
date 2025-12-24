# 🌐 net-tools

A comprehensive CLI network utility written in Go that provides essential network
diagnostic commands with both traditional CLI interface and a modern Terminal User
Interface (TUI). Perfect for network administrators, developers, and anyone who needs
quick network diagnostics.

## ✨ Features

### 🔍 DNS Lookup

- Comprehensive DNS record lookup (A, MX, NS, CNAME, TXT records)
- Clean, formatted output with detailed information
- Support for any domain name resolution

### 🔌 Port Connectivity Testing

- Test TCP port connectivity (netcat-like functionality)
- Connection status validation with detailed error reporting
- Support for any host and port combination

### 🌐 IP Address Information

- Get your private IP address
- Retrieve your public IP address
- Display both private and public IPs simultaneously

### 📊 Network Connections (Netstat)

- List active network connections
- Display local addresses, ports, and associated processes
- Responsive table formatting that adapts to terminal width

### 🖥️ Terminal User Interface (TUI)

- Modern, tabbed interface built with Bubble Tea
- Responsive design that adapts to any terminal size
- Text-selectable output for easy copying
- Intuitive keyboard navigation
- Real-time loading indicators and status updates

## 🚀 Installation

### Install Both Binaries from GitHub Releases

```bash
# Install net-tools CLI
go install github.com/catpaladin/net-tools/cmd/net-tools@latest

# Install net-tui TUI
go install github.com/catpaladin/net-tools/cmd/net-tui@latest
```

### Install via Script

```bash
curl -L https://raw.githubusercontent.com/catpaladin/net-tools/main/scripts/install.sh | bash
```

### Build from Source

```bash
git clone https://github.com/catpaladin/net-tools.git
cd net-tools
make install          # Install net-tools CLI
make install-tui      # Install net-tui TUI
```

### Cross-platform Installation

```bash
# Build for specific platforms
make build-linux      # Linux binary
make build-windows    # Windows binary
make build-mac        # macOS binary
make build-arm        # ARM binary

# Universal macOS binary (arm64 + amd64)
make build-darwin-universal
```

## 🎯 Usage

### net-tools CLI (Command Line Interface)

#### DNS Lookup

```bash
# Look up all DNS records for a domain
net-tools dig example.com

# Examples
net-tools dig google.com
net-tools dig github.com
```

#### Port Connectivity Test

```bash
# Test if a port is open on a host
net-tools nc <host> <port>

# Examples
net-tools nc google.com 80
net-tools nc localhost 3000
net-tools nc 192.168.1.1 22
```

#### IP Address Information

```bash
# Get both private and public IP (default)
net-tools ip

# Get only private IP
net-tools ip -t private

# Get only public IP
net-tools ip -t public

# Get both private and public IP
net-tools ip -t both
```

#### Network Connections

```bash
# List all active network connections
net-tools netstat
```

### net-tui (Terminal User Interface)

Launch the standalone TUI with all network tools:

```bash
net-tui
```

The TUI provides:

- **Tab navigation**: Use `←` `→` arrow keys to switch between tools
- **Field navigation**: Use `Tab` or `↑` `↓` to navigate form fields
- **Execution**: Press `Enter` to run commands
- **Scrolling**: Use `↑` `↓` to scroll through results
- **Text selection**: Click and drag to select and copy output text
- **Real-time status**: Loading indicators and error messages

## 📋 Examples

### DNS Lookup Example

```bash
$ net-tools dig example.com

🔍 DNS Lookup Results for example.com

A Records:
  93.184.216.34

MX Records:
  No MX records found

NS Records:
  a.iana-servers.net
  b.iana-servers.net

CNAME Record:
  No CNAME record found

TXT Records:
  v=spf1 -all
```

### Port Test Example

```bash
$ net-tools nc google.com 80

🔌 Port Connectivity Test

✓ Connection successful
Target: google.com:80
Status: Port is open and accepting connections
```

### IP Information Example

```bash
$ net-tools ip -t both

🌐 IP Address Information

Private IP: 192.168.1.100
Public IP: 203.0.113.45
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes,
please open an issue first to discuss what you would like to change.

- Fork this repo
- Create a PR from your feature branch using
  [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#summary)
- Please include a detailed description of your changes
- PR will be reviewed

### Development Guidelines

- Follow the existing code style and patterns
- Add tests for new functionality
- Update documentation as needed
- Ensure TUI changes maintain text selectability and responsive design

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for
details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - Powerful CLI framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Modern TUI framework
- [Charm](https://charm.sh) - Beautiful tools for the command line
