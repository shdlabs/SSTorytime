module SSTorytime.Decoder exposing (chapterDecoder, nodeDecoder, responseDecoder)

import Json.Decode as Decode exposing (Decoder)
import Json.Decode.Pipeline as Pipeline
import SSTorytime.Types exposing (..)



-- Decoder for 3D Position


position3DDecoder : Decoder Position3D
position3DDecoder =
    Decode.succeed Position3D
        |> Pipeline.required "X" Decode.float
        |> Pipeline.required "Y" Decode.float
        |> Pipeline.required "Z" Decode.float
        |> Pipeline.required "R" Decode.float
        |> Pipeline.required "Lat" Decode.float
        |> Pipeline.required "Lon" Decode.float



-- Decoder for Arrow Types


arrowTypeDecoder : String -> ArrowType
arrowTypeDecoder str =
    case str of
        "causality" ->
            Causality

        "dependency" ->
            Dependency

        "similarity" ->
            Similarity

        "hierarchy" ->
            Hierarchy

        _ ->
            Association



-- Decoder for Spherical Coordinates


sphericalCoordinateDecoder : Decoder SphericalCoordinate
sphericalCoordinateDecoder =
    Decode.succeed SphericalCoordinate
        |> Pipeline.required "longitude" Decode.float
        |> Pipeline.required "latitude" Decode.float
        |> Pipeline.required "radius" Decode.float



-- Decoder for Connections


connectionDecoder : Decoder Connection
connectionDecoder =
    Decode.succeed Connection
        |> Pipeline.required "fromNode" Decode.string
        |> Pipeline.required "toNode" Decode.string
        |> Pipeline.required "arrowType" (Decode.map arrowTypeDecoder Decode.string)
        |> Pipeline.required "sphericalCoord" sphericalCoordinateDecoder
        |> Pipeline.required "strength" Decode.float



-- Decoder for Time Context


timeContextDecoder : Decoder TimeContext
timeContextDecoder =
    Decode.succeed TimeContext
        |> Pipeline.required "duration" Decode.float
        |> Pipeline.required "startMoment" Decode.string
        |> Pipeline.required "endMoment" Decode.string
        |> Pipeline.required "continuity" Decode.string
        |> Pipeline.required "contextType" Decode.string



-- Decoder for Nodes


nodeDecoder : Decoder Node
nodeDecoder =
    Decode.succeed Node
        |> Pipeline.required "Text" Decode.string
        |> Pipeline.required "XYZ" position3DDecoder
        |> Pipeline.optional "Reln" (Decode.list connectionDecoder) []
        |> Pipeline.optional "TimeContext" (Decode.maybe timeContextDecoder) Nothing
        |> Pipeline.hardcoded "node"



-- Default node type
-- Decoder for Chapter


chapterDecoder : Decoder Chapter
chapterDecoder =
    Decode.succeed Chapter
        |> Pipeline.required "Chapter" Decode.string
        |> Pipeline.required "XYZ" position3DDecoder
        |> Pipeline.optional "Context" (Decode.list nodeDecoder) []
        |> Pipeline.optional "Single" (Decode.list nodeDecoder) []
        |> Pipeline.optional "Common" (Decode.maybe (Decode.list nodeDecoder)) Nothing



-- Main Response Decoder


responseDecoder : Decoder SSToryTimeResponse
responseDecoder =
    Decode.field "Response" <|
        (Decode.succeed SSToryTimeResponse
            |> Pipeline.required "Intent" Decode.string
            |> Pipeline.required "Orbits" (Decode.list chapterDecoder)
            |> Pipeline.required "Time" Decode.string
            |> Pipeline.required "Ambient" Decode.string
        )
