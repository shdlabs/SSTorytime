# Docker Setup

All Docker-related files are in the `DockerSST/` directory.

## Quick Start

```bash
cd DockerSST
./start_docker.sh
```

## Directory Structure

```
SSTorytime/
├── DockerSST/
│   ├── Dockerfile
│   ├── docker-compose.yml
│   ├── .dockerignore
│   ├── init-db.sql
│   ├── start_docker.sh       # Main startup script
│   ├── Makefile.docker
│   ├── DOCKER.md             # Full documentation
│   ├── DOCKER_QUICKSTART.md  # Quick reference
│   └── README.md             # This file
└── ... (rest of project)
```

## Usage

All Docker commands should be run from the `DockerSST/` directory:

```bash
cd DockerSST

# Start services
./start_docker.sh
# or
docker compose up -d

# View logs
docker compose logs -f server

# Stop services
docker compose down

# Load example data
make -f Makefile.docker populate-db
```

See [DOCKER_QUICKSTART.md](DOCKER_QUICKSTART.md) for common commands or [DOCKER.md](DOCKER.md) for complete documentation.
