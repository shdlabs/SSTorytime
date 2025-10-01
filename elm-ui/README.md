# SSTorytime Elm UI - Semantic SpaceTime Visualization

This is an Elm-based visualization frontend for SSTorytime, implementing true Semantic SpaceTime concepts based on Promise Theory.

## Features

### Semantic SpaceTime Concepts
- **3D Graph Theory with Vectors**: Nodes positioned in XYZ coordinates
- **Spherical Connections**: Connections between nodes represented on sphere surfaces using longitude/latitude
- **5 Arrow Types from Promise Theory**:
  - Causality (Cause → Effect)
  - Dependency (Requirement → Dependent)  
  - Similarity (Concept ~ Concept)
  - Hierarchy (Parent → Child)
  - Association (Loose connection)
- **Time Representation with Continuity**: Time ranges that maintain context of when/where they occurred

### Technical Features
- **Elm Architecture**: Type-safe, functional reactive programming
- **3D Rendering**: WebGL-based 3D visualization using elm-3d-scene
- **Interactive Exploration**: Navigate through the semantic space
- **Real-time API Integration**: Connect to SSTorytime backend

## Quick Start

### Prerequisites
- [Elm](https://elm-lang.org/) (0.19.1 or later)
- SSTorytime backend running on `localhost:8080`

### Installation & Setup

1. **Install Elm dependencies**:
   ```bash
   cd elm-ui
   make install
   ```

2. **Build the application**:
   ```bash
   make build
   ```

3. **Start development server**:
   ```bash
   make dev
   ```
   
   This will build the app and serve it on `http://localhost:3002`

### Development Workflow

- **Watch mode** (auto-rebuild on changes):
  ```bash
  make watch  # Requires elm-live: npm install -g elm-live
  ```

- **Manual development**:
  ```bash
  make dev    # Build and serve
  ```

- **Production build**:
  ```bash
  make build-prod
  ```

## Architecture

### Core Modules

- **`SSTorytime.Types`**: Core data types for Semantic SpaceTime
- **`SSTorytime.Decoder`**: JSON decoders for API responses
- **`SSTorytime.Api`**: HTTP client for backend communication
- **`SSTorytime.Visualization`**: 3D rendering and visualization logic
- **`Main`**: Application entry point and state management

### Data Model

The visualization represents knowledge as:

```elm
type alias Node =
    { text : String
    , position : Position3D  -- XYZ coordinates + sphere info
    , connections : List Connection
    , timeContext : Maybe TimeContext
    , nodeType : String
    }

type alias Connection =
    { fromNode : String
    , toNode : String
    , arrowType : ArrowType  -- 5 Promise Theory types
    , sphericalCoord : SphericalCoordinate  -- Position on sphere wall
    , strength : Float
    }
```

### 3D Visualization

- **Nodes**: Rendered as spheres with different colors/sizes based on type
- **Spheres**: Each node has a surrounding sphere showing connection space
- **Connections**: Rendered as arrows with colors indicating relationship type
- **Camera**: Interactive 3D camera for exploration

## API Integration

The Elm frontend communicates with the SSTorytime backend via HTTP POST to `/searchN4L`:

```elm
searchN4L : String -> Cmd Msg
searchN4L query =
    Http.post
        { url = "http://localhost:8080/searchN4L"
        , body = Http.multipartBody [ Http.stringPart "name" query ]
        , expect = Http.expectJson SearchCompleted responseDecoder
        }
```

## Semantic SpaceTime Concepts

### Graph Theory with Vectors
Each node exists in 3D space with:
- `X, Y, Z`: Spatial coordinates in the semantic space
- `R`: Radius of the node's connection sphere
- `Lat, Lon`: Used for spherical coordinate calculations

### Promise Theory Arrow Types
1. **Causality**: Direct cause-effect relationships (red arrows)
2. **Dependency**: One concept requires another (blue arrows)
3. **Similarity**: Related or analogous concepts (green arrows)
4. **Hierarchy**: Parent-child or categorical relationships (yellow arrows)
5. **Association**: General connections (gray arrows)

### Spherical Connections
Unlike traditional graph edges, connections are positioned on the surface of spheres around nodes:
- Each node has a sphere of radius `R`
- Connected nodes appear on this sphere's surface
- Position determined by longitude/latitude coordinates
- Preserves spatial relationships and maintains Promise Theory semantics

### Time Representation
Time contexts preserve continuity:
- Duration: How long something lasted
- Start/End moments: When it occurred
- Continuity: Where in the ongoing time stream
- Context type: What kind of temporal event

This allows representing "5 hours from 14:00-19:00 on Tuesday in the context of project development" rather than just "5 hours".

## Development

### Adding New Features

1. **New data types**: Add to `SSTorytime.Types`
2. **API changes**: Update decoders in `SSTorytime.Decoder`
3. **Visualization**: Extend `SSTorytime.Visualization`
4. **UI**: Modify `Main.elm`

### Debugging

Elm's compiler provides excellent error messages. Common issues:
- **Type mismatches**: Check decoder definitions match API response
- **Missing imports**: Ensure all modules are properly imported
- **JSON decoding**: Use browser dev tools to inspect API responses

## Deployment

For production deployment:

1. Build optimized version:
   ```bash
   make build-prod
   ```

2. Serve the `public/` directory with any web server
3. Ensure SSTorytime backend is accessible at the configured URL

## Contributing

This implementation focuses on the mathematical and conceptual accuracy of Semantic SpaceTime visualization. Key areas for contribution:
- Enhanced 3D interactions (zoom, pan, rotate)
- Additional Promise Theory relationship types
- Time dimension visualization improvements
- Performance optimizations for large knowledge graphs
- Mobile/responsive design enhancements

## License

Part of the SSTorytime project.