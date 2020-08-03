package inspector

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	width  = 600
	height = 400
)

// AudioInspector displays an oscilloscope-like view of live audio data
type AudioInspector struct {
	InputChannel chan [4]float64

	scope1 *audioScope
	scope2 *audioScope
	scope3 *audioScope
	scope4 *audioScope

	win *pixelgl.Window
	pic *pixel.PictureData
	buf []color.RGBA
}

// NewAudioInspector constructs a valid AudioInspector struct
func NewAudioInspector() *AudioInspector {
	buf := make([]color.RGBA, width*height)

	ai := &AudioInspector{
		InputChannel: make(chan [4]float64, 43690),
		buf:          buf,
		scope1:       newAudioScope(width, height/4, buf[180000:240000]),
		scope2:       newAudioScope(width, height/4, buf[120000:180000]),
		scope3:       newAudioScope(width, height/4, buf[60000:120000]),
		scope4:       newAudioScope(width, height/4, buf[0:60000]),
	}

	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:     "Audio Inspector",
		Bounds:    pixel.R(0, 0, float64(width), float64(height)),
		Resizable: true,
	})
	if err != nil {
		panic(err)
	}
	ai.win = win

	ai.pic = &pixel.PictureData{
		Pix:    ai.buf,
		Stride: width,
		Rect:   pixel.R(0, 0, float64(width), float64(height)),
	}

	return ai
}

// Render computes a new scope frame for each channel, and renders it to the inspector window
func (ai *AudioInspector) Render() {
	ai.updateScopes()

	picture := pixel.Picture(ai.pic)
	sprite := pixel.NewSprite(picture, picture.Bounds())
	sprite.Draw(ai.win, pixel.IM)

	shift := ai.win.Bounds().Size().Scaled(0.5).Sub(pixel.ZV)
	mat := pixel.IM.ScaledXY(pixel.ZV, pixel.V(1, 1)).Moved(shift)
	ai.win.SetMatrix(mat)

	ai.win.Update()
}

func (ai *AudioInspector) updateScopes() {
	var data1 []float64
	var data2 []float64
	var data3 []float64
	var data4 []float64

	numSamples := 0

consume:
	for {
		if numSamples > 1024 {
			break consume
		}

		select {
		case s := <-ai.InputChannel:
			data1 = append(data1, s[0]*4)
			data2 = append(data2, s[1]*4)
			data3 = append(data3, s[2]*4)
			data4 = append(data4, s[3]*4)

			numSamples++
		default:
			break consume
		}
	}

	if numSamples > 0 {
		ai.scope1.updateFrame(data1)
		ai.scope2.updateFrame(data2)
		ai.scope3.updateFrame(data3)
		ai.scope4.updateFrame(data4)
	}
}
