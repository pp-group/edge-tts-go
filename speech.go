package edge_tts_go

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/pp-group/edge-tts-go/biz/service/tts/edge"
)

type Speech struct {
	*edge.Communicate
	file     *os.File
	Folder   string
	fileName string
	storage  IStorage
}

func (s *Speech) GetFileName() string {
	return s.fileName
}

func (s *Speech) GenTTS() error {
	fileName := s.Folder + "/" + generateHashName(s.Text, s.VoiceLangRegion) + ".mp3"
	s.fileName = fileName
	isExist, err := s.storage.Exist(fileName)
	if err != nil {
		return err
	}
	if isExist {
		return nil
	}
	err = s.storage.Create(fileName)
	if err != nil {
		return err
	}
	defer s.file.Close()

	err = s.gen()
	if err != nil {
		return err
	}
	return nil
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
			s.storage.Write(data)
		}
	}
	s.storage.Close()
	return nil

}

func generateHashName(name, voice string) string {
	hash := sha256.Sum256([]byte(name))
	return fmt.Sprintf("%s_%s", voice, hex.EncodeToString(hash[:]))
}
