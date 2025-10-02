module Main exposing (main)

import Browser
import Html exposing (..)
import Html.Attributes exposing (placeholder, value)
import Html.Events exposing (onClick, onInput)
import Http
import SSTorytime.Types exposing (..)
import SSTorytime.Visualization as Viz



-- DEMO DATA


createDemoData : SSToryTimeResponse
createDemoData =
    { intent = "Demo visualization of Promise Theory concepts"
    , time = "2024-10-02T12:00:00Z"
    , ambient = "Interactive 3D exploration environment"
    , chapters =
        [ { title = "Promise Theory Demo"
          , position = Position3D 0.0 0.0 0.0 2.0 0.0 0.0
          , context =
                [ { text = "Central Agent"
                  , nodeType = "chapter"
                  , position = Position3D 0.0 0.0 0.0 1.0 0.0 0.0
                  , timeContext =
                        Just
                            { duration = 365.0
                            , startMoment = "2024-01-01"
                            , endMoment = "2024-12-31"
                            , continuity = "continuous"
                            , contextType = "persistent"
                            }
                  , connections =
                        [ { fromNode = "Central Agent"
                          , toNode = "Knowledge Node"
                          , arrowType = Causality
                          , sphericalCoord = SphericalCoordinate 0.785 1.57 1.0
                          , strength = 0.8
                          }
                        ]
                  }
                ]
          , single =
                [ { text = "Knowledge Node"
                  , nodeType = "single"
                  , position = Position3D 1.0 0.5 0.3 0.8 0.5 0.3
                  , timeContext =
                        Just
                            { duration = 180.0
                            , startMoment = "2024-06-01"
                            , endMoment = "2024-12-01"
                            , continuity = "periodic"
                            , contextType = "learning"
                            }
                  , connections =
                        [ { fromNode = "Knowledge Node"
                          , toNode = "Research Node"
                          , arrowType = Similarity
                          , sphericalCoord = SphericalCoordinate 2.35 1.0 1.2
                          , strength = 0.6
                          }
                        ]
                  }
                ]
          , common =
                Just
                    [ { text = "Shared Context"
                      , nodeType = "common"
                      , position = Position3D 0.2 -1.0 0.8 0.6 -0.8 0.4
                      , timeContext =
                            Just
                                { duration = 365.0
                                , startMoment = "2024-01-01"
                                , endMoment = "2024-12-31"
                                , continuity = "background"
                                , contextType = "ambient"
                                }
                      , connections = []
                      }
                    ]
          }
        ]
    }



-- MAIN


main : Program () Model Msg
main =
    Browser.element
        { init = init
        , update = update
        , subscriptions = subscriptions
        , view = view
        }



-- MODEL


type alias Model =
    { searchQuery : String
    , searchResult : String
    , status : String
    , showDemo : Bool
    , demoData : SSToryTimeResponse
    }


init : () -> ( Model, Cmd Msg )
init _ =
    ( { searchQuery = ""
      , searchResult = ""
      , status = "Ready to explore Semantic SpaceTime"
      , showDemo = True
      , demoData = createDemoData
      }
    , Cmd.none
    )



-- UPDATE


type Msg
    = UpdateSearchQuery String
    | PerformSearch
    | GotSearchResult (Result Http.Error String)
    | ToggleDemo


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        UpdateSearchQuery query ->
            ( { model | searchQuery = query }, Cmd.none )

        PerformSearch ->
            ( { model | status = "Searching the knowledge graph..." }
            , searchSSTorytime model.searchQuery
            )

        GotSearchResult result ->
            case result of
                Ok response ->
                    ( { model
                        | searchResult = response
                        , status = "Search completed - exploring semantic connections"
                        , showDemo = False
                      }
                    , Cmd.none
                    )

                Err error ->
                    ( { model
                        | searchResult = "Error: " ++ httpErrorString error
                        , status = "Error occurred"
                      }
                    , Cmd.none
                    )

        ToggleDemo ->
            ( { model | showDemo = not model.showDemo }, Cmd.none )



-- HTTP


searchSSTorytime : String -> Cmd Msg
searchSSTorytime query =
    Http.post
        { url = "http://localhost:8080/searchN4L"
        , body =
            Http.multipartBody
                [ Http.stringPart "name" query ]
        , expect = Http.expectString GotSearchResult
        }


httpErrorString : Http.Error -> String
httpErrorString error =
    case error of
        Http.BadUrl url ->
            "Bad URL: " ++ url

        Http.Timeout ->
            "Request timed out"

        Http.NetworkError ->
            "Network error"

        Http.BadStatus status ->
            "Bad status: " ++ String.fromInt status

        Http.BadBody body ->
            "Bad response body: " ++ body



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.none



-- VIEW


view : Model -> Html Msg
view model =
    div []
        [ div []
            [ h2 []
                [ text "SSTorytime - Semantic SpaceTime Explorer" ]
            , p []
                [ text "Promise Theory Graph Visualization with Spherical Connections" ]
            , button [ onClick ToggleDemo ]
                [ text
                    (if model.showDemo then
                        "Hide Demo"

                     else
                        "Show Demo"
                    )
                ]
            ]
        , div []
            [ div []
                [ input
                    [ placeholder "Enter search query (e.g., 'moon', 'reasoning', 'brains')"
                    , value model.searchQuery
                    , onInput UpdateSearchQuery
                    ]
                    []
                , button
                    [ onClick PerformSearch
                    ]
                    [ text "Explore SpaceTime" ]
                ]
            ]
        , div []
            [ p []
                [ text ("Status: " ++ model.status) ]
            ]
        , if model.showDemo then
            div []
                [ h3 [] [ text "3D Promise Theory Visualization" ]
                , Viz.view3D model.demoData
                ]

          else if String.isEmpty model.searchResult then
            welcomeSection

          else
            resultsSection model.searchResult
        ]


welcomeSection : Html Msg
welcomeSection =
    div []
        [ h1 []
            [ text "Welcome to Semantic SpaceTime" ]
        , p []
            [ text "Explore knowledge through Promise Theory with 3D nodes, spherical connections, and temporal dimensions." ]
        , div []
            [ h1 []
                [ text "Features:" ]
            , div []
                [ p [] [ text "• 3D nodes positioned in XYZ coordinates" ]
                , p [] [ text "• Spherical connections with longitude/latitude positioning" ]
                , p [] [ text "• 5 Arrow types from Promise Theory (Causality, Dependency, Similarity, Hierarchy, Association)" ]
                , p [] [ text "• Time representation with continuity context" ]
                , p [] [ text "• Interactive exploration of semantic relationships" ]
                ]
            ]
        ]


resultsSection : String -> Html Msg
resultsSection result =
    div []
        [ h3 []
            [ text "Semantic SpaceTime Results:" ]
        , pre []
            [ text result ]
        ]
