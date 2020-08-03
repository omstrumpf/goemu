package audio

import (
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	console_audio "github.com/omstrumpf/goemu/internal/app/console/audio"
	"github.com/omstrumpf/goemu/internal/app/log"
)

// TODO playing to speaker is super delayed

// Player plays audio to the speakers
type Player struct {
	InputChannel chan console_audio.Sample

	muted bool
}

// NewPlayer constructs a valid Player struct
func NewPlayer(bitrate int) *Player {
	p := &Player{
		InputChannel: make(chan console_audio.Sample, bitrate/6),
	}

	sampleRate := beep.SampleRate(bitrate)
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))

	streamer := beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		numStreamed := 0

		for i := range samples {
			select {
			case sample := <-p.InputChannel:
				if !p.muted {
					samples[i][0] = sample.L()
					samples[i][1] = sample.R()
				} else {
					samples[i][0] = 0
					samples[i][1] = 0
				}
				numStreamed++
			default:
				log.Tracef("ran out of samples")
				break
			}
		}

		return numStreamed, true
	})

	speaker.Play(streamer)

	return p
}

// SetMute sets the muted setting on the Player
func (p *Player) SetMute(muted bool) {
	p.muted = muted
}
