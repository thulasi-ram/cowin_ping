package main

import (
	"context"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"log"
	"os"
	"time"
)

func makeSound(ctx context.Context) {
	streamer, format := readSound()
	_ = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)

	speaker.Play(
		beep.Seq(
			streamer,
			beep.Callback(func() {
				cancel()
			}), // cancel after streaming is done
		),
	)
	select {
	case <-ctx.Done():
		streamer.Close()
		return
	}

}

func readSound() (beep.StreamSeekCloser, beep.Format) {

	f, err := os.Open("alarm.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	return streamer, format
}
