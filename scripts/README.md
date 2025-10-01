# SSTorytime Installation Scripts

This directory contains automated installation scripts to simplify setting up SSTorytime.

## Available Scripts

### `install-sstorytime.sh`
Complete installation script that sets up everything:
- PostgreSQL database
- Go programming language
- SSTorytime tools
- Sample data

**Usage:**
```bash
# Download and run
curl -fsSL https://raw.githubusercontent.com/markburgess/SSTorytime/main/scripts/install-sstorytime.sh | bash

# Or run locally
./install-sstorytime.sh

# With options
./install-sstorytime.sh --skip-postgres --skip-go
```

**Options:**
- `--skip-postgres`: Skip PostgreSQL installation
- `--skip-go`: Skip Go installation  
- `--skip-build`: Skip building SSTorytime
- `--skip-test`: Skip installation tests
- `--help`: Show help

### `setup-postgres.sh`
PostgreSQL-only installation and configuration:
- Installs PostgreSQL and extensions
- Creates database and user
- Configures authentication
- Tests connection

**Usage:**
```bash
# Download and run
curl -fsSL https://raw.githubusercontent.com/markburgess/SSTorytime/main/scripts/setup-postgres.sh | bash

# Or run locally
./setup-postgres.sh

# With custom settings
./setup-postgres.sh --db-name mydb --db-user myuser --db-pass mypass
```

**Options:**
- `--skip-install`: Skip PostgreSQL installation (configure only)
- `--db-name NAME`: Database name (default: sstoryline)
- `--db-user USER`: Database user (default: sstoryline)
- `--db-pass PASS`: Database password (default: sst_1234)
- `--help`: Show help

## Requirements

These scripts work on:
- Ubuntu 18.04+ / Debian 10+
- Fedora 30+ / CentOS 8+ / RHEL 8+
- Arch Linux / Manjaro

Requirements:
- `sudo` access (scripts will prompt when needed)
- Internet connection
- Basic development tools (git, curl, wget, make)

## What Gets Installed

After running the complete installation:

### PostgreSQL
- PostgreSQL server and client
- `postgresql-contrib` extensions
- Database: `sstoryline`
- User: `sstoryline` (password: `sst_1234`)
- Configured for local password authentication

### Go
- Go compiler and tools
- PostgreSQL driver (`github.com/lib/pq`)

### SSTorytime
- All command-line tools: `N4L`, `searchN4L`, `http_server`, etc.
- Sample data loaded from `examples/tutorial.n4l`
- Web server ready to run on port 8080

## After Installation

1. **Start the web interface:**
   ```bash
   cd SSTorytime
   ./http_server
   # Open http://localhost:8080
   ```

2. **Load your own data:**
   ```bash
   ./N4L -u your-file.n4l
   ```

3. **Search from command line:**
   ```bash
   ./searchN4L your-search-terms
   ```

## Troubleshooting

### Script Fails
- Run with individual components: use `--skip-*` options
- Check system logs: `sudo journalctl -f` 
- Try manual installation: see [docs/quick-install.md](../docs/quick-install.md)

### PostgreSQL Issues
- Check service status: `sudo systemctl status postgresql`
- Test connection: `psql -h localhost -U sstoryline -d sstoryline`
- View logs: `sudo journalctl -u postgresql`

### Permission Errors
- Don't run scripts as root (they use sudo when needed)
- Ensure your user has sudo privileges
- Check file ownership: `ls -la`

### Network Issues
- Verify internet connection
- Try manual download of prerequisites
- Use local package manager instead

## Manual Installation

If the scripts don't work for your system, follow the manual guides:
- [PostgreSQL Setup](../docs/postgresql-setup.md)
- [Go Setup](../docs/go-setup.md)
- [Quick Installation Guide](../docs/quick-install.md)

## Contributing

To improve these scripts:
1. Test on your Linux distribution
2. Report issues with system details
3. Submit pull requests with fixes
4. Add support for new distributions

## Security Notes

These scripts are designed for development/personal use:
- Use simple passwords by default
- Configure local-only database access
- Don't enforce security hardening

For production use, customize the database credentials and security settings.