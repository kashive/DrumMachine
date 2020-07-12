package main

import (
	"container/ring"
	tm "github.com/buger/goterm"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/jedib0t/go-pretty/table"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

//player domain object. Currently only has Play, but we can add more methods if we want to extend. Eg. stop, changeBpm etc.
type Player struct {
}

//plays the provided track aka prints the track on console and also plays audio
func (player *Player) Play(track Track) {
	stepsLen := len(track.Steps)
	circular := buildTrackRing(track, stepsLen)
	initSpeaker()
	tableWriter := buildTable(track.TrackInput)
	tm.Clear() //clearing any existing content
	playbackConfig := track.TrackInput.PlaybackConfig
	sleepDuration := time.Duration((60*1000)/playbackConfig.Bpm) * time.Millisecond
	for {
		start := time.Now()
		step := circular.Value.(Step)
		appendFooter(stepsLen, step, tableWriter)
		printTable(track.Name, playbackConfig.Bpm, tableWriter)
		playSoundsInStep(step)
		if !playbackConfig.Loop && step.Id == stepsLen {
			break
		}
		circular = circular.Next()
		since := time.Since(start)
		//subtract the time taken to play so that we are more consistent. May result in -ve in which case it will
		//return without sleeping
		time.Sleep(sleepDuration - since)
	}
}

//initializes speaker. Need to initialize only once
func initSpeaker() {
	sampleRate := beep.SampleRate(44100) //hard coding for now. All wav files used for testing have this sample rate
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))
}

//example at: https://github.com/faiface/beep/wiki/Hello,-Beep!
func playSoundsInStep(step Step) {
	if step.Sounds == nil || len(step.Sounds) == 0 {
		return
	}
	streamers, err := getSoundStreamersInStep(step)
	if err != nil {
		//todo: do more
		return
	}
	done := make(chan bool)
	mixedStreamer := beep.Mix(streamers...)
	speaker.Play(beep.Seq(mixedStreamer, beep.Callback(func() {
		done <- true
	})))
	<-done
	//closing streamers after playing sound for now
	//ideally we should be resetting the streamer to first position and reusing it
	for _, streamer := range streamers {
		streamer.(beep.StreamSeekCloser).Close()
	}
}

//returns sound streamers in the step
func getSoundStreamersInStep(step Step) ([]beep.Streamer, error) {
	//get streamers for each sound and play them simultaneously
	var streamers []beep.Streamer
	for _, sound := range step.Sounds {
		streamer, err := createStreamer(sound.PathToFile)
		if err != nil {
			return nil, err
		}
		streamers = append(streamers, streamer)
	}

	return streamers, nil
}

//prints table to console
func printTable(trackName string, bpm int, tableWriter table.Writer) {
	tm.MoveCursor(1, 1) //ensures that we overwrite existing
	tm.Println("Name: " + trackName)
	tm.Println("BPM: " + strconv.Itoa(bpm))
	tm.Println(tableWriter.Render())
	tm.Flush()
}

func appendFooter(numOfSteps int, step Step, tableWriter table.Writer) {
	tableWriter.ResetFooters()
	var footer table.Row
	for i := 0; i <= numOfSteps; i++ { //start at 0 to skip the first column
		if i == step.Id {
			footer = append(footer, "^")
		} else {
			footer = append(footer, "")
		}
	}
	tableWriter.AppendFooter(footer)
}

//creates the streamer if not created yet
func createStreamer(audioFilePath string) (beep.StreamSeekCloser, error){
	absPathToInputFile, _ := filepath.Abs("../" +  audioFilePath)
	//no need to close f here. Closing the streamer takes care of it. See: https://github.com/faiface/beep/wiki/Hello,-Beep!
	f, err := os.Open(absPathToInputFile)
	if err != nil {
		return nil, err
	}
	//assuming wav files only for now
	streamer, _, err := wav.Decode(f)
	if err != nil {
		return nil, err
	}
	return streamer, nil
}

//prints a textual representation of the track on the console
func buildTable(trackInput TrackInput) table.Writer{
	var header table.Row
	header = append(header, "STEP")
	for i := 1; i <= trackInput.NumberOfSteps; i++ {
		header = append(header, strconv.Itoa(i))
	}
	t := table.NewWriter()
	t.AppendHeader(header)
	for _, pattern := range trackInput.Patterns {
		var row table.Row
		sound := pattern.Sound
		steps := pattern.Steps
		row = append(row, sound.Name)
		for i := 1; i <= trackInput.NumberOfSteps; i++ {
			if contains(steps, i) {
				row = append(row, "X")
			}else {
				row = append(row, "_")
			}
		}
		t.AppendRow(row)
	}
	return t
}
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func buildTrackRing(track Track, numOfSteps int) *ring.Ring {
	circular := ring.New(numOfSteps)
	for _, step := range track.Steps {
		circular.Value = step
		circular = circular.Next()
	}
	return circular
}