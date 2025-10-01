# Simple Go Setup for SSTorytime

This guide will help you quickly install and configure Go (Golang) for building SSTorytime.

## Quick Installation

### Option 1: Package Manager (Recommended)

#### Ubuntu/Debian
```bash
# Install Go from official repos (usually recent enough)
sudo apt update
sudo apt install golang-go

# Verify installation
go version
```

#### Fedora/CentOS/RHEL
```bash
# Install Go
sudo dnf install golang
# or for older systems: sudo yum install golang

# Verify installation
go version
```

#### Arch Linux
```bash
# Install Go
sudo pacman -S go

# Verify installation
go version
```

### Option 2: Official Installer (Latest Version)

If you need the latest Go version:

```bash
# Download and install (replace with latest version)
GO_VERSION="1.21.5"  # Check https://golang.org/dl/ for latest
wget "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz"

# Remove old installation and install new
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"

# Add to PATH in your shell config file
echo 'export PATH=/usr/local/go/bin:$PATH' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
```

## Configure Go Environment

### Simple Setup (Recommended)

Modern Go (1.13+) uses modules by default, which simplifies things:

```bash
# Set up workspace directory
mkdir -p ~/go-workspace
cd ~/go-workspace

# Clone SSTorytime
git clone https://github.com/markburgess/SSTorytime.git
cd SSTorytime

# Initialize Go module (if not already done)
go mod tidy

# Test that Go can build the project
make http_server
```

### Legacy Setup (If You Have Issues)

If you encounter module-related problems:

```bash
# Disable Go modules (fallback to old GOPATH mode)
go env -w GO111MODULE=off

# Set up traditional GOPATH
export GOPATH=~/go
export PATH=$PATH:$GOPATH/bin

# Create workspace structure
mkdir -p ~/go/{bin,src,pkg}

# Link SSTorytime to GOPATH
ln -s ~/path/to/SSTorytime ~/go/src/SSTorytime

# Install required dependencies
go get github.com/lib/pq
```

## Test Your Setup

```bash
cd /path/to/SSTorytime

# Test compilation
make clean
make all

# If successful, you should see executables:
ls -la N4L searchN4L http_server text2N4L

# Test a simple build
echo "package main; import \"fmt\"; func main() { fmt.Println(\"Go works!\") }" > test.go
go run test.go
rm test.go
```

## Troubleshooting

### Common Issues

**"go: command not found"**
```bash
# Check if Go is in your PATH
echo $PATH | grep go

# Add Go to PATH (choose the appropriate path)
echo 'export PATH=/usr/local/go/bin:$PATH' >> ~/.bashrc
# or
echo 'export PATH=/usr/bin/go:$PATH' >> ~/.bashrc

# Reload shell
source ~/.bashrc
```

**Module errors**
```bash
# Try disabling modules
go env -w GO111MODULE=off

# Or ensure you're in the right directory
cd SSTorytime
go mod tidy
```

**Build errors about missing packages**
```bash
# Install PostgreSQL driver
go get github.com/lib/pq

# Or if using modules
go mod download
```

**Permission errors**
```bash
# Don't use sudo with go commands
# Make sure you own your GOPATH/workspace
chown -R $USER:$USER ~/go-workspace
```

### Verify Everything Works

```bash
# Check Go installation
go version

# Check environment
go env GOPATH GOROOT

# Test SSTorytime build
cd /path/to/SSTorytime
make clean && make http_server

# Test the server
./http_server &
curl "http://localhost:8080/?search=test"
kill %1  # Stop background server
```

## Next Steps

Once Go is working:

1. **Set up PostgreSQL**: Follow the [PostgreSQL setup guide](postgresql-setup.md)
2. **Build SSTorytime**: Run `make all` to build all tools
3. **Try examples**: Upload sample data and test the system
4. **Read documentation**: Check out the [Tutorial](Tutorial.md)

## Advanced Configuration

### Performance Tuning

```bash
# Enable Go proxy for faster downloads
go env -w GOPROXY=https://proxy.golang.org,direct

# Enable checksum verification
go env -w GOSUMDB=sum.golang.org

# Set module cache location (optional)
go env -w GOMODCACHE=/path/to/large/disk/go-mod-cache
```

### IDE Setup

For development with VS Code:
```bash
# Install Go extension
code --install-extension golang.go

# Or for other editors, make sure Go is in PATH
which go
```

## Getting Help

If you're still having issues:

1. Check Go installation: `go version`
2. Check SSTorytime build: `cd SSTorytime && make clean && make all`
3. Ask for help on [SSTorytime discussions](https://github.com/markburgess/SSTorytime/discussions)
4. Check [Go documentation](https://golang.org/doc/) for Go-specific issues