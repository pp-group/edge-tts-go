package edge_tts_go

import (
	"testing"

	"edge-tts-go/biz/service/tts/edge"
)

func TestGenTTS(t *testing.T) {
	fileName, err := genTTS("cpdd,cpdd,cpdd")
	if err != nil {
		t.Errorf("genTTS fail, err: %v", err)
		return
	}
	t.Logf("genTTS success, fileName: %s", fileName)
}

// genTTS template use edge-tts-go
func genTTS(text string) (string, error) {
	c, err := edge.NewCommunicate(text)
	if err != nil {
		return "", err
	}

	speech := &Speech{
		Communicate: c,
		Folder:      "templateaudio",
	}
	err = speech.GenTTS()
	if err != nil {
		return "", err
	}
	return speech.GetFileName(), nil
}
