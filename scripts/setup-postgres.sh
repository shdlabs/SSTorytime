#!/bin/bash
#
# SSTorytime PostgreSQL Setup Script
# Automatically installs and configures PostgreSQL for SSTorytime
#

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default database configuration
DB_NAME="sstoryline"
DB_USER="sstoryline"
DB_PASS="sst_1234"

# Print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect Linux distribution
detect_distro() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        DISTRO=$ID
        VERSION=$VERSION_ID
    else
        print_error "Cannot detect Linux distribution"
        exit 1
    fi
    print_status "Detected distribution: $DISTRO $VERSION"
}

# Install PostgreSQL based on distribution
install_postgresql() {
    print_status "Installing PostgreSQL..."
    
    case "$DISTRO" in
        ubuntu|debian)
            sudo apt update
            sudo apt install -y postgresql postgresql-contrib
            ;;
        fedora)
            sudo dnf install -y postgresql postgresql-server postgresql-contrib
            sudo postgresql-setup --initdb
            ;;
        centos|rhel)
            sudo yum install -y postgresql postgresql-server postgresql-contrib
            sudo postgresql-setup initdb
            ;;
        arch|manjaro)
            sudo pacman -S --noconfirm postgresql
            sudo -u postgres initdb -D /var/lib/postgres/data
            ;;
        *)
            print_error "Unsupported distribution: $DISTRO"
            print_error "Please install PostgreSQL manually and run this script with --skip-install"
            exit 1
            ;;
    esac
    
    print_success "PostgreSQL installed"
}

# Start and enable PostgreSQL service
start_postgresql() {
    print_status "Starting PostgreSQL service..."
    
    sudo systemctl enable postgresql
    sudo systemctl start postgresql
    
    # Wait a moment for service to fully start
    sleep 2
    
    if sudo systemctl is-active --quiet postgresql; then
        print_success "PostgreSQL service started"
    else
        print_error "Failed to start PostgreSQL service"
        print_error "Check logs with: sudo journalctl -u postgresql"
        exit 1
    fi
}

# Create database and user
setup_database() {
    print_status "Creating database user and database..."
    
    # Create user with password
    sudo -u postgres psql -c "CREATE USER $DB_USER WITH PASSWORD '$DB_PASS' SUPERUSER;" 2>/dev/null || {
        print_warning "User $DB_USER already exists, updating password..."
        sudo -u postgres psql -c "ALTER USER $DB_USER PASSWORD '$DB_PASS';"
    }
    
    # Create database
    sudo -u postgres createdb -O $DB_USER $DB_NAME 2>/dev/null || {
        print_warning "Database $DB_NAME already exists"
    }
    
    # Add required extension
    sudo -u postgres psql -d $DB_NAME -c "CREATE EXTENSION IF NOT EXISTS unaccent;"
    
    print_success "Database setup completed"
}

# Configure authentication
configure_auth() {
    print_status "Configuring PostgreSQL authentication..."
    
    # Find pg_hba.conf location
    HBA_FILE=$(sudo -u postgres psql -t -P format=unaligned -c 'SHOW hba_file;' | tr -d ' ')
    
    if [ ! -f "$HBA_FILE" ]; then
        print_error "Could not find pg_hba.conf file at: $HBA_FILE"
        exit 1
    fi
    
    # Backup original file
    sudo cp "$HBA_FILE" "${HBA_FILE}.backup"
    
    # Update authentication method to md5 for password authentication
    sudo sed -i 's/local   all             all                                     peer/local   all             all                                     md5/' "$HBA_FILE"
    sudo sed -i 's/host    all             all             127.0.0.1\/32            ident/host    all             all             127.0.0.1\/32            md5/' "$HBA_FILE"
    sudo sed -i 's/host    all             all             ::1\/128                 ident/host    all             all             ::1\/128                 md5/' "$HBA_FILE"
    
    print_success "Authentication configured"
}

# Restart PostgreSQL to apply configuration changes
restart_postgresql() {
    print_status "Restarting PostgreSQL to apply configuration..."
    sudo systemctl restart postgresql
    sleep 2
    print_success "PostgreSQL restarted"
}

# Test the database connection
test_connection() {
    print_status "Testing database connection..."
    
    # Test connection with password
    export PGPASSWORD="$DB_PASS"
    if psql -h localhost -U $DB_USER -d $DB_NAME -c "SELECT version();" >/dev/null 2>&1; then
        print_success "Database connection test passed!"
        print_success "You can connect with: psql -h localhost -U $DB_USER -d $DB_NAME"
        print_success "Password: $DB_PASS"
    else
        print_error "Database connection test failed"
        print_error "Try connecting manually: psql -h localhost -U $DB_USER -d $DB_NAME"
        print_error "Check PostgreSQL logs: sudo journalctl -u postgresql"
        exit 1
    fi
    unset PGPASSWORD
}

# Print final instructions
print_instructions() {
    echo ""
    print_success "PostgreSQL setup completed successfully!"
    echo ""
    echo "Database Configuration:"
    echo "  Host: localhost"
    echo "  Port: 5432"
    echo "  Database: $DB_NAME"
    echo "  User: $DB_USER"
    echo "  Password: $DB_PASS"
    echo ""
    echo "Next steps:"
    echo "  1. Install Go: https://golang.org/dl/"
    echo "  2. Build SSTorytime: cd /path/to/SSTorytime && make all"
    echo "  3. Try examples: cd examples && make && ../N4L -u tutorial.n4l"
    echo ""
    echo "To connect manually:"
    echo "  psql -h localhost -U $DB_USER -d $DB_NAME"
    echo ""
}

# Command line options
SKIP_INSTALL=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-install)
            SKIP_INSTALL=true
            shift
            ;;
        --db-name)
            DB_NAME="$2"
            shift 2
            ;;
        --db-user)
            DB_USER="$2"
            shift 2
            ;;
        --db-pass)
            DB_PASS="$2"
            shift 2
            ;;
        -h|--help)
            echo "SSTorytime PostgreSQL Setup Script"
            echo ""
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --skip-install    Skip PostgreSQL installation (assume already installed)"
            echo "  --db-name NAME    Database name (default: sstoryline)"
            echo "  --db-user USER    Database user (default: sstoryline)"
            echo "  --db-pass PASS    Database password (default: sst_1234)"
            echo "  -h, --help        Show this help"
            echo ""
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Main execution
main() {
    echo "SSTorytime PostgreSQL Setup Script"
    echo "=================================="
    echo ""
    
    # Check if running as root
    if [ "$EUID" -eq 0 ]; then
        print_error "Please do not run this script as root"
        print_error "The script will use sudo when needed"
        exit 1
    fi
    
    # Check if sudo is available
    if ! command -v sudo &> /dev/null; then
        print_error "sudo is required but not installed"
        exit 1
    fi
    
    detect_distro
    
    if [ "$SKIP_INSTALL" = false ]; then
        install_postgresql
    else
        print_status "Skipping PostgreSQL installation"
    fi
    
    start_postgresql
    setup_database
    configure_auth
    restart_postgresql
    test_connection
    print_instructions
}

# Run main function
main "$@"