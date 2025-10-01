module SSTorytime.Api exposing (searchN4L, Msg(..))

import Http
import Json.Decode as Decode
import SSTorytime.Types exposing (SSToryTimeResponse)
import SSTorytime.Decoder exposing (responseDecoder)


type Msg
    = SearchCompleted (Result Http.Error SSToryTimeResponse)


-- Search the N4L knowledge base
searchN4L : String -> Cmd Msg
searchN4L query =
    Http.post
        { url = "http://localhost:8080/searchN4L"
        , body = Http.multipartBody 
            [ Http.stringPart "name" query ]
        , expect = Http.expectJson SearchCompleted responseDecoder
        }