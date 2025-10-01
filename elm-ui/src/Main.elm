module Main exposing (main)

import Browser
import Html exposing (Html, div, h1, input, button, text, pre, p)
import Html.Attributes exposing (placeholder, value, style)
import Html.Events exposing (onInput, onClick)
import Http
import Json.Decode as Decode

-- MAIN

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
    }


init : () -> (Model, Cmd Msg)
init _ =
    ( { searchQuery = ""
      , searchResult = ""
      , status = "Ready to explore Semantic SpaceTime"
      }
    , Cmd.none
    )


-- UPDATE

type Msg
    = UpdateSearchQuery String
    | PerformSearch
    | GotSearchResult (Result Http.Error String)


update : Msg -> Model -> (Model, Cmd Msg)
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
                      }
                    , Cmd.none
                    )

                Err error ->
                    ( { model 
                        | searchResult = "Error: " ++ (httpErrorString error)
                        , status = "Error occurred"
                      }
                    , Cmd.none
                    )


-- HTTP

searchSSTorytime : String -> Cmd Msg
searchSSTorytime query =
    Http.post
        { url = "http://localhost:8080/searchN4L"
        , body = Http.multipartBody 
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
subscriptions model =
    Sub.none


-- VIEW

view : Model -> Html Msg
view model =
    div 
        [ style "padding" "30px"
        , style "font-family" "Inter, -apple-system, BlinkMacSystemFont, sans-serif"
        , style "max-width" "1200px"
        , style "margin" "0 auto"
        , style "background" "linear-gradient(135deg, #667eea 0%, #764ba2 100%)"
        , style "min-height" "100vh"
        , style "color" "white"
        ]
        [ div 
            [ style "background" "rgba(255, 255, 255, 0.1)"
            , style "padding" "25px"
            , style "border-radius" "15px"
            , style "backdrop-filter" "blur(10px)"
            , style "margin-bottom" "30px"
            , style "text-align" "center"
            ]
            [ h1 
                [ style "font-size" "2.5rem"
                , style "margin-bottom" "10px"
                ] 
                [ text "SSTorytime - Semantic SpaceTime Explorer" ]
            , p 
                [ style "font-size" "1.2rem"
                , style "margin-bottom" "20px"
                , style "opacity" "0.9"
                ]
                [ text "Promise Theory Graph Visualization with Spherical Connections" ]
            ]
        
        , div 
            [ style "background" "rgba(255, 255, 255, 0.1)"
            , style "padding" "25px"
            , style "border-radius" "15px"
            , style "backdrop-filter" "blur(10px)"
            , style "margin-bottom" "30px"
            ]
            [ div 
                [ style "display" "flex"
                , style "gap" "15px"
                , style "align-items" "center"
                , style "justify-content" "center"
                , style "flex-wrap" "wrap"
                ]
                [ input
                    [ placeholder "Enter search query (e.g., 'moon', 'reasoning', 'brains')"
                    , value model.searchQuery
                    , onInput UpdateSearchQuery
                    , style "padding" "15px"
                    , style "width" "400px"
                    , style "border" "none"
                    , style "border-radius" "10px"
                    , style "font-size" "16px"
                    , style "background" "rgba(255, 255, 255, 0.9)"
                    , style "color" "#333"
                    ]
                    []
                , button
                    [ onClick PerformSearch
                    , style "padding" "15px 25px"
                    , style "background" "rgba(255, 255, 255, 0.2)"
                    , style "color" "white"
                    , style "border" "2px solid rgba(255, 255, 255, 0.3)"
                    , style "border-radius" "10px"
                    , style "cursor" "pointer"
                    , style "font-size" "16px"
                    , style "font-weight" "600"
                    , style "transition" "all 0.3s ease"
                    ]
                    [ text "Explore SpaceTime" ]
                ]
            ]
        
        , div 
            [ style "background" "rgba(255, 255, 255, 0.05)"
            , style "padding" "20px"
            , style "border-radius" "10px"
            , style "margin-bottom" "20px"
            ]
            [ p 
                [ style "margin" "0"
                , style "font-size" "18px"
                , style "font-weight" "500"
                ]
                [ text ("Status: " ++ model.status) ]
            ]
        
        , if String.isEmpty model.searchResult then
            welcomeSection
          else
            resultsSection model.searchResult
        ]


welcomeSection : Html Msg
welcomeSection =
    div 
        [ style "background" "rgba(255, 255, 255, 0.1)"
        , style "padding" "30px"
        , style "border-radius" "15px"
        , style "text-align" "center"
        ]
        [ h1 
            [ style "font-size" "1.8rem"
            , style "margin-bottom" "20px"
            ]
            [ text "Welcome to Semantic SpaceTime" ]
        , p 
            [ style "font-size" "1.1rem"
            , style "margin-bottom" "25px"
            , style "line-height" "1.6"
            ]
            [ text "Explore knowledge through Promise Theory with 3D nodes, spherical connections, and temporal dimensions." ]
        , div 
            [ style "text-align" "left"
            , style "max-width" "600px"
            , style "margin" "0 auto"
            ]
            [ h1 
                [ style "font-size" "1.4rem"
                , style "margin-bottom" "15px"
                ]
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
    div 
        [ style "background" "rgba(255, 255, 255, 0.1)"
        , style "padding" "25px"
        , style "border-radius" "15px"
        ]
        [ h1 
            [ style "font-size" "1.5rem"
            , style "margin-bottom" "20px"
            , style "color" "#fff"
            ]
            [ text "Semantic SpaceTime Results:" ]
        , pre 
            [ style "background" "rgba(0, 0, 0, 0.3)"
            , style "padding" "20px"
            , style "border-radius" "10px"
            , style "overflow" "auto"
            , style "max-height" "600px"
            , style "font-family" "Monaco, Consolas, monospace"
            , style "font-size" "14px"
            , style "line-height" "1.4"
            , style "white-space" "pre-wrap"
            , style "color" "#f8f9fa"
            ]
            [ text result ]
        ]