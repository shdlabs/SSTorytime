#!/bin/bash
#
# SSTorytime Complete Installation Script
# Installs PostgreSQL, Go, and builds SSTorytime
#

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Install prerequisites
install_prerequisites() {
    print_status "Installing prerequisites..."
    
    case "$DISTRO" in
        ubuntu|debian)
            sudo apt update
            sudo apt install -y git curl wget make
            ;;
        fedora)
            sudo dnf install -y git curl wget make
            ;;
        centos|rhel)
            sudo yum install -y git curl wget make
            ;;
        arch|manjaro)
            sudo pacman -S --noconfirm git curl wget make
            ;;
        *)
            print_warning "Unknown distribution, assuming prerequisites are installed"
            ;;
    esac
    
    print_success "Prerequisites installed"
}

# Install Go
install_go() {
    print_status "Installing Go..."
    
    if command -v go &> /dev/null; then
        print_success "Go already installed: $(go version)"
        return
    fi
    
    case "$DISTRO" in
        ubuntu|debian)
            sudo apt install -y golang-go
            ;;
        fedora)
            sudo dnf install -y golang
            ;;
        centos|rhel)
            sudo yum install -y golang
            ;;
        arch|manjaro)
            sudo pacman -S --noconfirm go
            ;;
        *)
            print_error "Cannot install Go automatically for $DISTRO"
            print_error "Please install Go manually and re-run this script"
            exit 1
            ;;
    esac
    
    # Verify Go installation
    if command -v go &> /dev/null; then
        print_success "Go installed: $(go version)"
    else
        print_error "Go installation failed"
        exit 1
    fi
}

# Install PostgreSQL using our setup script
install_postgresql() {
    print_status "Installing and configuring PostgreSQL..."
    
    # Download and run PostgreSQL setup script
    if [ -f "./scripts/setup-postgres.sh" ]; then
        # Use local script if available
        bash ./scripts/setup-postgres.sh
    else
        # Download from GitHub
        curl -fsSL https://raw.githubusercontent.com/markburgess/SSTorytime/main/scripts/setup-postgres.sh | bash
    fi
    
    print_success "PostgreSQL setup completed"
}

# Clone SSTorytime repository
clone_sstorytime() {
    print_status "Cloning SSTorytime repository..."
    
    if [ -d "SSTorytime" ]; then
        print_warning "SSTorytime directory already exists, updating..."
        cd SSTorytime
        git pull origin main || print_warning "Failed to update repository"
    else
        git clone https://github.com/markburgess/SSTorytime.git
        cd SSTorytime
    fi
    
    print_success "SSTorytime repository ready"
}

# Build SSTorytime
build_sstorytime() {
    print_status "Building SSTorytime..."
    
    # Ensure we're in the right directory
    if [ ! -f "Makefile" ]; then
        print_error "Not in SSTorytime directory or Makefile missing"
        exit 1
    fi
    
    # Clean and build
    make clean || print_warning "Clean failed (continuing...)"
    
    # Install Go dependencies
    if [ -f "go.mod" ]; then
        print_status "Installing Go dependencies..."
        go mod tidy
        go mod download
    else
        print_status "Installing PostgreSQL driver..."
        go get github.com/lib/pq
    fi
    
    # Build all tools
    make all
    
    print_success "SSTorytime built successfully"
}

# Test the installation
test_installation() {
    print_status "Testing installation..."
    
    # Test database connection
    export PGPASSWORD="sst_1234"
    if psql -h localhost -U sstoryline -d sstoryline -c "SELECT version();" >/dev/null 2>&1; then
        print_success "Database connection test passed"
    else
        print_error "Database connection test failed"
        return 1
    fi
    unset PGPASSWORD
    
    # Test that executables exist
    for tool in N4L searchN4L http_server text2N4L; do
        if [ -x "./$tool" ]; then
            print_success "$tool built successfully"
        else
            print_error "$tool not found or not executable"
            return 1
        fi
    done
    
    # Test web server briefly
    print_status "Testing web server..."
    ./http_server &
    SERVER_PID=$!
    sleep 2
    
    if curl -s "http://localhost:8080/?search=test" >/dev/null 2>&1; then
        print_success "Web server test passed"
    else
        print_warning "Web server test failed (this may be normal if no data is loaded)"
    fi
    
    # Stop test server
    kill $SERVER_PID 2>/dev/null || true
    sleep 1
    
    print_success "Installation test completed"
}

# Load sample data
load_examples() {
    print_status "Loading sample data..."
    
    if [ -d "examples" ]; then
        cd examples
        if make >/dev/null 2>&1; then
            # Try to load tutorial data
            if ../N4L -u tutorial.n4l >/dev/null 2>&1; then
                print_success "Sample data loaded (tutorial.n4l)"
            else
                print_warning "Failed to load tutorial.n4l (continuing...)"
            fi
        else
            print_warning "Failed to build examples (continuing...)"
        fi
        cd ..
    else
        print_warning "Examples directory not found"
    fi
}

# Print final instructions
print_final_instructions() {
    echo ""
    echo "================================================"
    print_success "SSTorytime installation completed successfully!"
    echo "================================================"
    echo ""
    echo "What you can do now:"
    echo ""
    echo "1. Start the web interface:"
    echo "   ./http_server"
    echo "   Then open: http://localhost:8080"
    echo ""
    echo "2. Load your own data:"
    echo "   ./N4L -u your-file.n4l"
    echo ""
    echo "3. Search from command line:"
    echo "   ./searchN4L your-search-terms"
    echo ""
    echo "4. Learn more:"
    echo "   Read docs/Tutorial.md"
    echo "   Check docs/N4L.md for data format"
    echo "   See docs/search_examples.md for search examples"
    echo ""
    echo "Database connection info:"
    echo "   Host: localhost:5432"
    echo "   Database: sstoryline"
    echo "   User: sstoryline"
    echo "   Password: sst_1234"
    echo ""
    echo "For help: https://github.com/markburgess/SSTorytime/discussions"
    echo ""
}

# Command line options
SKIP_POSTGRES=false
SKIP_GO=false
SKIP_BUILD=false
SKIP_TEST=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-postgres)
            SKIP_POSTGRES=true
            shift
            ;;
        --skip-go)
            SKIP_GO=true
            shift
            ;;
        --skip-build)
            SKIP_BUILD=true
            shift
            ;;
        --skip-test)
            SKIP_TEST=true
            shift
            ;;
        -h|--help)
            echo "SSTorytime Complete Installation Script"
            echo ""
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --skip-postgres   Skip PostgreSQL installation"
            echo "  --skip-go         Skip Go installation"
            echo "  --skip-build      Skip building SSTorytime"
            echo "  --skip-test       Skip installation tests"
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
    echo "SSTorytime Complete Installation Script"
    echo "======================================"
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
    install_prerequisites
    
    if [ "$SKIP_GO" = false ]; then
        install_go
    else
        print_status "Skipping Go installation"
    fi
    
    if [ "$SKIP_POSTGRES" = false ]; then
        install_postgresql
    else
        print_status "Skipping PostgreSQL installation"
    fi
    
    clone_sstorytime
    
    if [ "$SKIP_BUILD" = false ]; then
        build_sstorytime
    else
        print_status "Skipping build"
    fi
    
    if [ "$SKIP_TEST" = false ]; then
        test_installation
    else
        print_status "Skipping tests"
    fi
    
    load_examples
    print_final_instructions
}

# Run main function
main "$@"