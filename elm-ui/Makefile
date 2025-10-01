# Elm SSTorytime Visualization Makefile

.PHONY: build dev clean install serve

# Build the Elm application with Vite
build:
	npx vite build

# Build with optimization for production
build-prod:
	npx vite build --mode production

# Development mode with Vite dev server
dev:
	npx vite dev --port 3002 --host asus

# Install Elm dependencies
install:
	elm install

# Serve the application locally with Vite preview
serve:
	npx vite preview --port 3002 --host asus

# Clean build artifacts
clean:
	rm -rf dist/
	rm -rf elm-stuff/

# Development watch mode with Vite
watch:
	npx vite dev --port 3002 --host asus

# Help
help:
	@echo "Available commands:"
	@echo "  build      - Build the Elm application with Vite"
	@echo "  build-prod - Build with production optimization"
	@echo "  dev        - Start Vite development server with hot reload on asus:3002"
	@echo "  install    - Install Elm dependencies"
	@echo "  serve      - Serve built application on asus:3002"
	@echo "  watch      - Same as dev - Vite development server with hot reload"
	@echo "  clean      - Clean build artifacts"
	@echo "  help       - Show this help message"