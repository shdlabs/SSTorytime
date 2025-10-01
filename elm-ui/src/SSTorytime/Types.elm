module SSTorytime.Types exposing 
    ( Node
    , Connection
    , ArrowType(..)
    , SphericalCoordinate
    , TimeContext
    , SSToryTimeResponse
    , Chapter
    , Position3D
    , arrowTypeToColor
    , arrowTypeToString
    , sphericalToCartesian
    , cartesianToSpherical
    , calculateSphericalPath
    , calculateSemanticDistance
    )

import Json.Decode as Decode


-- 3D Position in Semantic SpaceTime
type alias Position3D =
    { x : Float
    , y : Float
    , z : Float
    , r : Float  -- Radius of the sphere
    , lat : Float  -- Latitude on sphere
    , lon : Float  -- Longitude on sphere
    }


-- The 5 Arrow Types from Promise Theory
type ArrowType
    = Causality        -- Cause -> Effect
    | Dependency       -- Requirement -> Dependent
    | Similarity       -- Concept ~ Concept  
    | Hierarchy        -- Parent -> Child
    | Association      -- Loose connection


-- Spherical coordinates for connections
type alias SphericalCoordinate =
    { longitude : Float  -- Direction on sphere wall
    , latitude : Float   -- Altitude on sphere wall
    , radius : Float     -- Distance to connected node
    }


-- Connection between nodes with Promise Theory semantics
type alias Connection =
    { fromNode : String
    , toNode : String
    , arrowType : ArrowType
    , sphericalCoord : SphericalCoordinate
    , strength : Float  -- Connection strength (0.0 - 1.0)
    }


-- Time representation that maintains continuity and context
type alias TimeContext =
    { duration : Float        -- How long (in hours/minutes)
    , startMoment : String    -- When it started (context preserved)
    , endMoment : String      -- When it ended
    , continuity : String     -- Where in the continuum
    , contextType : String    -- What kind of time period
    }


-- Node in Semantic SpaceTime
type alias Node =
    { text : String
    , position : Position3D
    , connections : List Connection
    , timeContext : Maybe TimeContext
    , nodeType : String  -- "chapter", "context", "single", "common"
    }


-- Chapter containing nodes
type alias Chapter =
    { title : String
    , position : Position3D
    , context : List Node
    , single : List Node
    , common : Maybe (List Node)
    }


-- Full SSTorytime API Response
type alias SSToryTimeResponse =
    { intent : String
    , chapters : List Chapter
    , time : String
    , ambient : String
    }


-- Helper Functions for Promise Theory Visualization

arrowTypeToColor : ArrowType -> String
arrowTypeToColor arrowType =
    case arrowType of
        Causality ->
            "#ff6b6b"  -- Red for causal relationships
        
        Dependency ->
            "#4ecdc4"  -- Teal for dependencies
        
        Similarity ->
            "#45b7d1"  -- Blue for similarities
        
        Hierarchy ->
            "#96ceb4"  -- Green for hierarchical relationships
        
        Association ->
            "#feca57"  -- Yellow for associations


arrowTypeToString : ArrowType -> String
arrowTypeToString arrowType =
    case arrowType of
        Causality -> "Causality"
        Dependency -> "Dependency" 
        Similarity -> "Similarity"
        Hierarchy -> "Hierarchy"
        Association -> "Association"


-- Spherical Mathematics for Node Connections
sphericalToCartesian : SphericalCoordinate -> ( Float, Float, Float )
sphericalToCartesian { latitude, longitude, radius } =
    let
        lat = degrees latitude
        lon = degrees longitude
        x = radius * cos lat * cos lon
        y = radius * cos lat * sin lon
        z = radius * sin lat
    in
    ( x, y, z )


cartesianToSpherical : Float -> Float -> Float -> SphericalCoordinate
cartesianToSpherical x y z =
    let
        radius = sqrt (x * x + y * y + z * z)
        latitude = asin (z / radius) |> (\rad -> rad * 180 / pi)
        longitude = atan2 y x |> (\rad -> rad * 180 / pi)
    in
    { latitude = latitude
    , longitude = longitude  
    , radius = radius
    }


-- Calculate connection path on sphere surface
calculateSphericalPath : Position3D -> Position3D -> List SphericalCoordinate
calculateSphericalPath from to =
    let
        fromSphere = cartesianToSpherical from.x from.y from.z
        toSphere = cartesianToSpherical to.x to.y to.z
        
        -- Create intermediate points for smooth arc on sphere surface
        steps = 10
        stepSize = 1.0 / toFloat steps
    in
    List.range 0 steps
        |> List.map (\i ->
            let
                t = toFloat i * stepSize
                -- Spherical linear interpolation (slerp)
                lat = fromSphere.latitude + t * (toSphere.latitude - fromSphere.latitude)
                lon = fromSphere.longitude + t * (toSphere.longitude - fromSphere.longitude)
                -- Use average radius for connection path
                radius = (fromSphere.radius + toSphere.radius) / 2
            in
            { latitude = lat, longitude = lon, radius = radius }
        )


-- Promise Theory Semantic Distance Calculation
calculateSemanticDistance : Node -> Node -> Float
calculateSemanticDistance node1 node2 =
    let
        -- Euclidean distance in 3D space
        dx = node1.position.x - node2.position.x
        dy = node1.position.y - node2.position.y  
        dz = node1.position.z - node2.position.z
        
        spatialDistance = sqrt (dx * dx + dy * dy + dz * dz)
    in
    spatialDistance