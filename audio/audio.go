package audio

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
	"github.com/kyeett/mogui/assets"
)

const (
	basePathAudio  = "assets/audio/"
	sampleRate     = 44100
	bytesPerSample = 4
)

var Sounds = map[string][]byte{}

// Player represents the current audio state.
type Player struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	current      time.Duration
	total        time.Duration
	seBytes      []byte
	volume128    int
}

var globalAudioContext *audio.Context

func LoadResources() {
	var err error
	globalAudioContext, err = audio.NewContext(sampleRate)
	if err != nil {
		log.Fatal(err)
	}

	for _, path := range assets.AssetNames() {
		if !strings.Contains(path, basePathAudio) || !strings.Contains(path, "mp3") {
			continue
		}

		fmt.Println(path)
		// s, err := vorbis.Decode(globalAudioContext, audio.BytesReadSeekCloser(assets.MustAsset(path)))
		s, err := mp3.Decode(globalAudioContext, audio.BytesReadSeekCloser(assets.MustAsset(path)))
		// s, err := wav.Decode(globalAudioContext, audio.BytesReadSeekCloser(assets.MustAsset(path)))
		if err != nil {
			log.Fatal("failed to load", err)
		}
		b, err := ioutil.ReadAll(s)
		if err != nil {
			log.Fatal("failed to read", err)
		}
		Sounds[path] = b
		fmt.Println("done1")
	}
	fmt.Println("done2")
	// fmt.Println(Sounds)
}

func Play(name string, volume ...float64) {
	b := Sounds[basePathAudio+name]
	fmt.Println("play audio!", len(b))
	if b == nil || len(b) == 0 {
		log.Println("tried to play empty bytes")
		return
	}
	tmpP, err := audio.NewPlayerFromBytes(globalAudioContext, b)
	if err != nil {
		log.Fatal(err)
	}
	if len(volume) > 0 {
		tmpP.SetVolume(volume[0])
	}
	tmpP.Play()
}
