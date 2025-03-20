


package main

import (
	"fmt"
	"strings"


	"github.com/go-tts/tts/pkg/audio"
	"github.com/go-tts/tts/pkg/speech"
)

// *********************************************

func main () {


	audioIn, err := speech.FromText(text, speech.LangEn)

	err := speech.WriteToAudioStream(textIn, audioOut, "it")

}