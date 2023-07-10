package handler

import (
	"edge-tts-go/biz/service/tts"
	"edge-tts-go/biz/service/tts/edge"
)

// GenTTS template use edge-tts-go
func GenTTS(text string) (string, error) {
	c, err := edge.NewCommunicate(text)
	if err != nil {
		return "", err
	}

	speech := &tts.Speech{
		Communicate: c,
		Folder:      "audio",
	}
	err = speech.GenTTS()
	if err != nil {
		return "", err
	}
	return speech.GetFileName(), nil
}
