# Quick Start Installation Guide

Get SSTorytime running on your Linux system in minutes.

## One-Command Setup (Easiest)

For Ubuntu/Debian systems, this single command will set up everything:

```bash
curl -fsSL https://raw.githubusercontent.com/markburgess/SSTorytime/main/scripts/install-sstorytime.sh | bash
```

## Step-by-Step Setup

If you prefer to do it manually or the one-command setup doesn't work:

### 1. Install PostgreSQL

Choose your method:

**Option A: Automated script**
```bash
curl -fsSL https://raw.githubusercontent.com/markburgess/SSTorytime/main/scripts/setup-postgres.sh | bash
```

**Option B: Manual setup**  
Follow the [detailed PostgreSQL guide](postgresql-setup.md)

### 2. Install Go

**Quick install:**
```bash
sudo apt install golang-go   # Ubuntu/Debian
sudo dnf install golang      # Fedora
sudo pacman -S go            # Arch
```

**Detailed setup:**  
Follow the [Go installation guide](go-setup.md)

### 3. Get SSTorytime

```bash
git clone https://github.com/markburgess/SSTorytime.git
cd SSTorytime
```

### 4. Build Everything

```bash
make all
```

### 5. Test Your Installation

```bash
# Start the web server
./http_server &

# Test it works
curl "http://localhost:8080/?search=test"

# Stop the server
kill %1
```

## Try the Examples

Load some sample data:

```bash
cd examples
make
../N4L -u tutorial.n4l
```

Browse your data:
```bash
../http_server &
# Open http://localhost:8080 in your browser
```

## What Next?

- **Learn the basics**: Read the [Tutorial](Tutorial.md)
- **Create your own data**: Learn [N4L syntax](N4L.md)
- **Search your data**: Try [search examples](search_examples.md)
- **Use the API**: Check the [API documentation](API.md)

## Need Help?

- **Common issues**: Check [Troubleshooting](#troubleshooting) below
- **PostgreSQL problems**: See [PostgreSQL setup guide](postgresql-setup.md)
- **Go problems**: See [Go setup guide](go-setup.md)
- **Ask questions**: [GitHub Discussions](https://github.com/markburgess/SSTorytime/discussions)

## Troubleshooting

### PostgreSQL Issues

**Connection refused:**
```bash
# Check if PostgreSQL is running
sudo systemctl status postgresql

# Start it if needed
sudo systemctl start postgresql
```

**Authentication failed:**
```bash
# Test connection
psql -h localhost -U sstoryline -d sstoryline
# Password: sst_1234
```

### Build Issues

**Go command not found:**
```bash
# Install Go
sudo apt install golang-go  # or your distro's package manager

# Verify
go version
```

**Build errors:**
```bash
# Clean and rebuild
make clean
make all

# Check for missing dependencies
go mod tidy
```

### Server Issues

**Port already in use:**
```bash
# Find what's using port 8080
sudo lsof -i :8080

# Kill it or use a different port
./http_server -port 8081
```

**Can't connect to database:**
```bash
# Test database connection
psql -h localhost -U sstoryline -d sstoryline -c "SELECT version();"

# If this fails, re-run PostgreSQL setup
```

## Custom Configuration

### Different Database Settings

Create `~/.SSTorytime` with your custom settings:
```
dbname: my_database
user: my_user
passwd: my_password
```

### Different Web Server Port

```bash
./http_server -port 9090
```

## Advanced Options

### Memory-Based PostgreSQL (Faster, Data Lost on Reboot)

For development or testing with better performance:

```bash
# Setup described in postgresql-setup.md under "Performance Options"
```

### Development Setup

For developers wanting to modify SSTorytime:

```bash
# Install development tools
go install golang.org/x/tools/...

# Run tests
make test

# Format code
go fmt ./...
```

## Architecture Overview

SSTorytime consists of:

- **PostgreSQL**: Stores your knowledge graph data
- **Go tools**: Command-line programs for data management
  - `N4L`: Load data from N4L files
  - `searchN4L`: Command-line search tool
  - `http_server`: Web interface
- **Web interface**: Browser-based exploration of your data
- **N4L language**: Simple syntax for creating knowledge maps

## Files and Directories

After installation you'll have:

```
SSTorytime/
├── N4L              # Data loader
├── searchN4L        # Command-line search
├── http_server      # Web server
├── examples/        # Sample data files
├── docs/            # Documentation
└── cmd/             # Source code
```

Your data files (`.n4l` format) can be stored anywhere and loaded with `N4L -u yourfile.n4l`.

---

*This simplified guide focuses on getting you up and running quickly. For detailed information about any component, see the specific guides linked above.*