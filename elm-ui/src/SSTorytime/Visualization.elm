module SSTorytime.Visualization exposing (nodeToMesh, sphericalConnection, view3D)

import Html exposing (Html)
import Html.Attributes
import SSTorytime.Types exposing (..)



-- Create a 3D CSS-based visualization using CSS transforms


nodeToMesh : Node -> Html msg
nodeToMesh node =
    let
        -- Convert position to CSS transform
        transform =
            "translate3d("
                ++ String.fromFloat (node.position.x * 100)
                ++ "px, "
                ++ String.fromFloat (node.position.y * 100)
                ++ "px, "
                ++ String.fromFloat (node.position.z * 100)
                ++ "px)"

        nodeColor =
            case node.nodeType of
                "chapter" ->
                    "#3b82f6"

                -- blue
                "context" ->
                    "#10b981"

                -- green
                "single" ->
                    "#f59e0b"

                -- orange
                "common" ->
                    "#8b5cf6"

                -- purple
                _ ->
                    "#6b7280"

        -- gray
    in
    Html.div
        [ Html.Attributes.class "node-3d"
        , Html.Attributes.style "transform" transform
        , Html.Attributes.style "background-color" nodeColor
        , Html.Attributes.attribute "data-node-type" node.nodeType
        ]
        [ Html.div [ Html.Attributes.class "node-label" ]
            [ Html.text node.text ]
        ]


sphericalConnection : Node -> Connection -> List (Html msg)
sphericalConnection fromNode connection =
    let
        -- Calculate spherical connection endpoint
        sphereRadius =
            fromNode.position.r

        lon =
            connection.sphericalCoord.longitude

        lat =
            connection.sphericalCoord.latitude

        -- Convert to cartesian
        toX =
            fromNode.position.x + sphereRadius * cos lat * cos lon

        toY =
            fromNode.position.y + sphereRadius * cos lat * sin lon

        toZ =
            fromNode.position.z + sphereRadius * sin lat

        -- Create CSS line between points
        fromTransform =
            "translate3d("
                ++ String.fromFloat (fromNode.position.x * 100)
                ++ "px, "
                ++ String.fromFloat (fromNode.position.y * 100)
                ++ "px, "
                ++ String.fromFloat (fromNode.position.z * 100)
                ++ "px)"

        toTransform =
            "translate3d("
                ++ String.fromFloat (toX * 100)
                ++ "px, "
                ++ String.fromFloat (toY * 100)
                ++ "px, "
                ++ String.fromFloat (toZ * 100)
                ++ "px)"

        arrowColor =
            case connection.arrowType of
                Causality ->
                    "#ef4444"

                -- red
                Dependency ->
                    "#3b82f6"

                -- blue
                Similarity ->
                    "#10b981"

                -- green
                Hierarchy ->
                    "#eab308"

                -- yellow
                Association ->
                    "#6b7280"

        -- gray
    in
    [ Html.div
        [ Html.Attributes.class "connection-3d"
        , Html.Attributes.style "transform" fromTransform
        , Html.Attributes.style "background-color" arrowColor
        , Html.Attributes.attribute "data-arrow-type" (arrowTypeToString connection.arrowType)
        ]
        []
    , Html.div
        [ Html.Attributes.class "connection-endpoint"
        , Html.Attributes.style "transform" toTransform
        , Html.Attributes.style "background-color" arrowColor
        ]
        []
    ]


view3D : SSToryTimeResponse -> Html msg
view3D response =
    let
        -- Collect all nodes from all chapters
        allNodes =
            List.concatMap
                (\chapter ->
                    chapter.context
                        ++ chapter.single
                        ++ Maybe.withDefault [] chapter.common
                )
                response.chapters

        -- Create 3D node visualizations
        nodeViews =
            List.map nodeToMesh allNodes

        -- Create connection visualizations
        connectionViews =
            List.concatMap
                (\node ->
                    List.concatMap (sphericalConnection node) node.connections
                )
                allNodes

        -- Combine all 3D elements
        allElements =
            nodeViews ++ connectionViews
    in
    Html.div
        [ Html.Attributes.class "scene-3d" ]
        [ Html.div
            [ Html.Attributes.class "viewport-3d" ]
            [ Html.h2 [] [ Html.text "3D Promise Theory Visualization" ]
            , Html.div
                [ Html.Attributes.class "world-3d" ]
                allElements
            ]
        ]
