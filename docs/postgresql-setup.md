# Simple PostgreSQL Setup for SSTorytime

This guide will help you quickly set up PostgreSQL for SSTorytime with minimal complexity. We focus on getting it working first - you can secure it later if needed.

## Quick Setup (Recommended)

For most users, this automated setup will get you running in minutes:

### 1. Download and run the setup script

```bash
# Download and run the PostgreSQL setup script
curl -fsSL https://raw.githubusercontent.com/markburgess/SSTorytime/main/scripts/setup-postgres.sh | bash

# Or if you prefer to review first:
wget https://raw.githubusercontent.com/markburgess/SSTorytime/main/scripts/setup-postgres.sh
chmod +x setup-postgres.sh
./setup-postgres.sh
```

### 2. Test your setup

```bash
# Test that everything works
psql -h localhost -U sstoryline -d sstoryline -c "SELECT version();"
# Enter password: sst_1234
```

If you see PostgreSQL version information, you're ready to go! Skip to [Next Steps](#next-steps).

---

## Manual Setup (Step by Step)

If you prefer to understand each step or the automated script didn't work:

### Step 1: Install PostgreSQL

Choose your Linux distribution:

#### Ubuntu/Debian
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

#### Fedora/CentOS/RHEL
```bash
sudo dnf install postgresql postgresql-server postgresql-contrib
# or for older systems:
sudo yum install postgresql postgresql-server postgresql-contrib
```

#### Arch Linux
```bash
sudo pacman -S postgresql
```

### Step 2: Initialize and Start PostgreSQL

#### Ubuntu/Debian (usually automatic)
```bash
sudo systemctl enable postgresql
sudo systemctl start postgresql
```

#### Fedora/CentOS/RHEL
```bash
sudo postgresql-setup --initdb
sudo systemctl enable postgresql
sudo systemctl start postgresql
```

#### Arch Linux
```bash
sudo -u postgres initdb -D /var/lib/postgres/data
sudo systemctl enable postgresql
sudo systemctl start postgresql
```

### Step 3: Create Database and User

```bash
# Switch to postgres user and create our database
sudo -u postgres createuser -P sstoryline
# Enter password when prompted: sst_1234
# Confirm password: sst_1234
# Superuser? y

# Create database
sudo -u postgres createdb -O sstoryline sstoryline

# Add required extension
sudo -u postgres psql -d sstoryline -c "CREATE EXTENSION IF NOT EXISTS unaccent;"
```

### Step 4: Configure Authentication

```bash
# Find and edit the pg_hba.conf file
sudo -u postgres psql -t -P format=unaligned -c 'SHOW hba_file;'
# This shows the location of pg_hba.conf file

# Edit the file (location from above command)
sudo nano /var/lib/pgsql/data/pg_hba.conf  # Common location

# Find these lines and change 'ident' to 'md5':
# local   all             all                                     peer
# host    all             all             127.0.0.1/32            ident
# host    all             all             ::1/128                 ident

# Change to:
# local   all             all                                     md5
# host    all             all             127.0.0.1/32            md5
# host    all             all             ::1/128                 md5
```

### Step 5: Restart PostgreSQL

```bash
sudo systemctl restart postgresql
```

### Step 6: Test Connection

```bash
psql -h localhost -U sstoryline -d sstoryline -c "SELECT version();"
# Enter password: sst_1234
```

---

## Testing Your Setup

Verify everything works:

```bash
# Test basic connection
psql -h localhost -U sstoryline -d sstoryline

# Inside psql, run:
\dt                    # List tables (should be empty initially)
\q                     # Quit psql

# Test from SSTorytime directory:
cd /path/to/SSTorytime
make http_server
./http_server &        # Start server in background
curl "http://localhost:8080/?search=test"  # Should return JSON response
```

---

## Custom Configuration (Optional)

If you want to use different database credentials, create `$HOME/.SSTorytime`:

```bash
cat > $HOME/.SSTorytime << 'EOF'
dbname: my_custom_db
user: my_user
passwd: my_secure_password
EOF
```

Then create your custom database:
```bash
sudo -u postgres createuser -P my_user
sudo -u postgres createdb -O my_user my_custom_db
sudo -u postgres psql -d my_custom_db -c "CREATE EXTENSION IF NOT EXISTS unaccent;"
```

---

## Troubleshooting

### Common Issues

**Connection refused errors:**
```bash
# Check if PostgreSQL is running
sudo systemctl status postgresql

# Check if it's listening on the right port
sudo netstat -nltp | grep 5432
```

**Authentication failed:**
```bash
# Verify user exists
sudo -u postgres psql -c "\du"

# Reset password
sudo -u postgres psql -c "ALTER USER sstoryline PASSWORD 'sst_1234';"
```

**Permission denied errors:**
```bash
# Check pg_hba.conf authentication methods
sudo -u postgres grep -v '^#' /var/lib/pgsql/data/pg_hba.conf | grep -v '^$'

# Restart after changes
sudo systemctl restart postgresql
```

**Can't find pg_hba.conf:**
```bash
# Find the file location
sudo -u postgres psql -t -P format=unaligned -c 'SHOW hba_file;'
```

### Getting Help

If you're still having issues:

1. Check PostgreSQL logs:
   ```bash
   sudo journalctl -u postgresql -f
   ```

2. Check if SSTorytime can connect:
   ```bash
   cd /path/to/SSTorytime
   go run cmd/tools/text2N4L.go --help
   ```

3. Ask for help on the [SSTorytime discussions](https://github.com/markburgess/SSTorytime/discussions)

---

## Performance Options

### RAM-based PostgreSQL (Advanced)

For better performance, you can run PostgreSQL in memory. **Warning: All data is lost on reboot!**

```bash
# Create memory filesystem
sudo mkdir -p /mnt/pg_ram
sudo mount -t tmpfs -o size=1G tmpfs /mnt/pg_ram
sudo chown postgres:postgres /mnt/pg_ram

# Stop regular PostgreSQL
sudo systemctl stop postgresql

# Initialize and start in-memory instance
sudo -u postgres /usr/lib/postgresql/*/bin/initdb -D /mnt/pg_ram/pgdata
sudo -u postgres /usr/lib/postgresql/*/bin/pg_ctl -D /mnt/pg_ram/pgdata -l /mnt/pg_ram/logfile start

# Repeat the database setup steps from above
```

---

## Next Steps

Once PostgreSQL is working:

1. **Install Go**: Follow the [Go installation guide](go-setup.md)
2. **Build SSTorytime**: Run `make all` in the SSTorytime directory
3. **Try examples**: Upload sample data with `make && ../N4L -u examples/tutorial.n4l`
4. **Start exploring**: Read the [Tutorial](Tutorial.md) or [Quick Start](README.md)

---

## Security Notes

This setup uses simple passwords and local-only access for ease of use. If you're running on a server or shared machine:

- Change the default password
- Configure PostgreSQL for your specific security requirements
- Consider using connection encryption
- Review PostgreSQL security documentation

For personal laptop use, this simple setup is usually fine.