package edge_tts_go

import (
	"fmt"
	"io"

	"github.com/pp-group/edge-tts-go/biz/service/tts/edge"
)

type Speech struct {
	*edge.Communicate
	Rw io.ReadWriter
}

func (s *Speech) GenTTS() error {
	// fileName := s.Folder + "/" + generateHashName(s.Text, s.VoiceLangRegion) + ".mp3"
	return s.gen()
}

func (s *Speech) gen() error {
	op, err := s.Stream()
	if err != nil {
		return err
	}
	defer s.CloseOutput()
	solveCount := 0
	audioData := make([][][]byte, s.AudioDataIndex)
	for i := range op {
		if _, ok := i["end"]; ok {
			solveCount++
			if solveCount == s.AudioDataIndex {
				break
			}
		}
		t, ok := i["type"]
		if ok && t == "audio" {
			data := i["data"].(edge.AudioData)
			audioData[data.Index] = append(audioData[data.Index], data.Data)
		}
		e, ok := i["error"]
		if ok {
			fmt.Printf("has error err: %v\n", e)
		}
	}
	// write data, sort by index
	for _, v := range audioData {
		for _, data := range v {
			if _, err := s.Rw.Write(data); err != nil {
				return err
			}
		}
	}
	return nil
}
