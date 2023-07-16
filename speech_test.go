package edge_tts_go

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/pp-group/edge-tts-go/biz/service/tts/edge"
)

func TestGenTTS(t *testing.T) {
	dir := "x"
	text := "cpdd,cpdd,cpdd"
	fileName := generateHashName(text, "zh-CN-XiaoxiaoNeural")
	f, err := os.OpenFile(fmt.Sprintf("%s/%s.mp3", dir, fileName), os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		t.Errorf("open file error %v", err)
		return
	}
	defer f.Close()

	if err := genTTS(text, f); err != nil {
		t.Errorf("genTTS fail, err: %v", err)
		return
	}
	t.Logf("genTTS success, fileName: %s", fileName)
}

func generateHashName(name, voice string) string {
	hash := sha256.Sum256([]byte(name))
	return fmt.Sprintf("%s_%s", voice, hex.EncodeToString(hash[:]))
}

// genTTS template use github.com/pp-group/edge-tts-go
func genTTS(text string, rw io.ReadWriter) error {
	c, err := edge.NewCommunicate(text)
	if err != nil {
		return err
	}
	speech := &Speech{
		Communicate: c,
		Rw:          rw,
	}
	err = speech.GenTTS()
	if err != nil {
		return err
	}
	return nil
}
