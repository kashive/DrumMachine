package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

//main function that is responsible for de-serializing, validating the input, building track and then running the player
func main() {

	arguments := os.Args[1:]
	absPathToInputFile, err := filepath.Abs("../" +  arguments[0])
	failIfError(err, "Encountered error finding absolute path to input file")
	trackInput, err := deserializeTrackInput(absPathToInputFile)
	failIfError(err, fmt.Sprintf("Encountered error when reading input file in path: %v", absPathToInputFile))

	track, err := BuildTrack(trackInput)
	failIfError(err, fmt.Sprintf("Failed while building track in path: %v. ", absPathToInputFile))
	var player Player
	player.Play(track)
}


//Unmarshal the JSON input to TrackInput struct
func deserializeTrackInput(absPath string) (TrackInput, error) {
	var trackInput TrackInput
	file, err := ioutil.ReadFile(absPath)
	if err != nil {
		return trackInput, err
	}
	err = json.Unmarshal(file, &trackInput)
	if err != nil {
		return trackInput, err
	}
	return trackInput, nil
}

func failIfError(err error, message string) {
	if err != nil {
		log.Fatal(message, err)
	}
}

//builds track from track input
func BuildTrack(trackInput TrackInput) (Track, error) {
	var err error
	var track Track
	if len(trackInput.Name) < 1 {
		err = errors.New("track input does not have a name")
		return track, err
	}
	if trackInput.NumberOfSteps < 1 {
		err = errors.New("track input does not have number of steps")
		return track, err
	}
	//grouping means that if there are duplicate sounds in the pattern then the last one wins
	soundsByStepId, err := groupByStepId(trackInput)
	if err != nil {
		return track, err
	}
	track = Track{
		Name:  trackInput.Name,
		Steps: []Step{},
		TrackInput: trackInput,
	}
	numberOfSteps := trackInput.NumberOfSteps
	for stepId := 1; stepId <= numberOfSteps; stepId++ {
		step := Step{Id: stepId}
		sounds := soundsByStepId[stepId]
		step.Sounds = sounds
		steps := append(track.Steps, step)
		track.Steps = steps
	}
	return track, err
}

//groups sounds in track input by step id aka {1 -> [{name: 'Hi Hat', pathToFile: '/tmp/'}]
func groupByStepId(trackInput TrackInput) (map[int][]Sound, error) {
	res := make(map[int][]Sound)
	patterns := trackInput.Patterns
	for _, pattern := range patterns {
		sound := pattern.Sound
		for _, step := range pattern.Steps {
			if step > trackInput.NumberOfSteps {
				return res, errors.New(
					fmt.Sprintf("Found stepId: %v which is greater than total number of steps: %v",
						step, trackInput.NumberOfSteps))
			}
			if val, ok := res[step]; ok {
				res[step] = append(val, sound)
			} else {
				res[step] = []Sound{sound}
			}
		}
	}
	return res, nil
}

