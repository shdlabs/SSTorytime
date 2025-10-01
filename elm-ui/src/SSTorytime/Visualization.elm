module SSTorytime.Visualization exposing (view3D, nodeToMesh, sphericalConnection)

import Html exposing (Html)
import Html.Attributes
import Scene3d
import Scene3d.Material as Material
import Camera3d
import Viewpoint3d
import Point3d
import Direction3d
import Length
import Angle
import Color
import Sphere3d
import Cylinder3d
import SSTorytime.Types exposing (..)
import Math.Vector3 as Vec3


-- Convert Node to 3D Mesh
nodeToMesh : Node -> Scene3d.Entity coordinates
nodeToMesh node =
    let
        position = Point3d.fromMeters 
            { x = node.position.x
            , y = node.position.y  
            , z = node.position.z
            }
        
        -- Different colors and sizes based on node type
        (nodeColor, nodeRadius) = case node.nodeType of
            "chapter" -> (Color.blue, 0.1)
            "context" -> (Color.green, 0.08)
            "single" -> (Color.orange, 0.06)
            "common" -> (Color.purple, 0.06)
            _ -> (Color.gray, 0.05)
        
        sphere = Sphere3d.atPoint position (Length.meters nodeRadius)
        material = Material.color nodeColor
    in
    Scene3d.sphere material sphere


-- Create spherical connection visualization
sphericalConnection : Node -> Connection -> List (Scene3d.Entity coordinates)
sphericalConnection fromNode connection =
    let
        fromPos = Point3d.fromMeters 
            { x = fromNode.position.x
            , y = fromNode.position.y
            , z = fromNode.position.z
            }
        
        -- Calculate position on sphere wall using spherical coordinates
        sphereRadius = fromNode.position.r
        lon = connection.sphericalCoord.longitude
        lat = connection.sphericalCoord.latitude
        
        -- Convert spherical to cartesian coordinates
        toX = fromNode.position.x + sphereRadius * cos(lat) * cos(lon)
        toY = fromNode.position.y + sphereRadius * cos(lat) * sin(lon)
        toZ = fromNode.position.z + sphereRadius * sin(lat)
        
        toPos = Point3d.fromMeters { x = toX, y = toY, z = toZ }
        
        -- Arrow color based on type
        arrowColor = case connection.arrowType of
            Causality -> Color.red
            Dependency -> Color.blue
            Similarity -> Color.green
            Hierarchy -> Color.yellow
            Association -> Color.gray
        
        -- Create cylinder connecting the points
        cylinder = Cylinder3d.from fromPos toPos (Length.meters 0.02)
        material = Material.color arrowColor
        
        -- Create arrow head (small cone at destination)
        arrowHead = Sphere3d.atPoint toPos (Length.meters 0.03)
    in
    [ Scene3d.cylinder material cylinder
    , Scene3d.sphere material arrowHead
    ]


-- Create sphere visualization around node
nodeSphere : Node -> Scene3d.Entity coordinates
nodeSphere node =
    let
        position = Point3d.fromMeters 
            { x = node.position.x
            , y = node.position.y
            , z = node.position.z
            }
        
        sphereRadius = node.position.r
        sphere = Sphere3d.atPoint position (Length.meters sphereRadius)
        
        -- Transparent sphere material
        material = Material.color (Color.rgba 0.5 0.5 0.5 0.1)
    in
    Scene3d.sphere material sphere


-- Main 3D Scene
view3D : SSToryTimeResponse -> Html msg
view3D response =
    let
        -- Camera setup
        camera = Camera3d.perspective
            { viewpoint = Viewpoint3d.lookAt
                { focalPoint = Point3d.origin
                , eyePoint = Point3d.fromMeters { x = 3, y = 3, z = 2 }
                , upDirection = Direction3d.positiveZ
                }
            , verticalFieldOfView = Angle.degrees 45
            }
        
        -- Collect all nodes from all chapters
        allNodes = List.concatMap 
            (\chapter -> 
                chapter.context 
                ++ chapter.single 
                ++ (Maybe.withDefault [] chapter.common)
            ) 
            response.chapters
        
        -- Create node meshes
        nodeMeshes = List.map nodeToMesh allNodes
        
        -- Create sphere visualizations
        sphereMeshes = List.map nodeSphere allNodes
        
        -- Create connection visualizations
        connectionMeshes = List.concatMap 
            (\node -> 
                List.concatMap (sphericalConnection node) node.connections
            ) 
            allNodes
        
        -- Combine all entities
        allEntities = nodeMeshes ++ sphereMeshes ++ connectionMeshes
    in
    Scene3d.sunny
        { camera = camera
        , clipDepth = Length.meters 0.1
        , dimensions = ( 800, 600 )
        , background = Scene3d.backgroundColor Color.black
        , entities = allEntities
        , upDirection = Direction3d.positiveZ
        , sunlightDirection = Direction3d.xyZ (Angle.degrees 45) (Angle.degrees 45)
        , shadows = False
        }