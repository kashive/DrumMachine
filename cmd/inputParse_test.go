package main

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestDeserialize(t *testing.T) {
	var trackInput TrackInput
	t.Run("path to input file does not exist", testDeserializeFunc("not existing path", trackInput))
	t.Run("json is malformed", testDeserializeFunc("/sampleInput/malformed.json", trackInput))
	t.Run("json does not match struct", testDeserializeFunc("/sampleInput/does_not_match_track_input.json", trackInput))
	trackInput = getValidTrackInput()
	t.Run("json is valid", testDeserializeFunc("/sampleInput/test-input.json", trackInput))
}

func TestBuildTrack(t *testing.T) {
	var track Track
	absPathTopLevelDir, _ := filepath.Abs("../")
	t.Run("track name does not exist", testBuildTrackFunc(absPathTopLevelDir + "/sampleInput/track_name_not_present.json", track))
	t.Run("stepId in sound pattern larger than total number of steps", testBuildTrackFunc(absPathTopLevelDir + "/sampleInput/step_num_in_pattern_greater_than_total.json", track))
	track = getValidTrack(absPathTopLevelDir + "/sampleInput/track_is_valid.json")
	t.Run("track is valid", testBuildTrackFunc("/sampleInput/track_is_valid.json", track))
}

func getValidTrack(path string) Track {
	trackInput, _ := deserializeTrackInput(path)
	return Track{
		Name:       "Awesome Track",
		Steps:      []Step{
			{
				Id:     1,
				Sounds: []Sound{
					{
						Name:       "Hi Hat",
						PathToFile: "/sampleInput/hi-hat.wav",
					},
					{
						Name:       "Snare Drum",
						PathToFile: "/sampleInput/snare.wav",
					},
				},
			},
			{
				Id:     2,
				Sounds: []Sound{
					{
						Name:       "Snare Drum",
						PathToFile: "/sampleInput/snare.wav",
					},
				},
			},

		},
		TrackInput: trackInput,
	}
}

func getValidTrackInput() TrackInput {

	return TrackInput{
		Name:          "Awesome Track",
		NumberOfSteps: 16,
		PlaybackConfig: PlaybackConfig{
			Bpm:  60,
			Loop: false,
		},
		Patterns: []Pattern{
			{
				Sound: Sound{
					Name:       "Hi Hat",
					PathToFile: "/sampleInput/hi-hat.wav",
				},
				Steps: []int{3, 7, 11, 15},
			},
			{
				Sound: Sound{
					Name:       "Snare Drum",
					PathToFile: "/sampleInput/snare.wav",
				},
				Steps: []int{5, 13},
			},
			{
				Sound: Sound{
					Name:       "Bass Drum",
					PathToFile: "/sampleInput/kick.wav",
				},
				Steps: []int{1, 5, 9, 13},
			},
		},
	}
}

func testBuildTrackFunc(pathToFile string, expected Track) func(*testing.T) {
	return func(t *testing.T) {
		absPath, _ := filepath.Abs("../" + pathToFile)
		trackInput, _ := deserializeTrackInput(absPath)
		actual, _ := BuildTrack(trackInput)
		assert.Equal(t, expected, actual)
	}
}

func testDeserializeFunc(pathToFile string, expected TrackInput) func(*testing.T) {
	return func(t *testing.T) {
		absPath, _ := filepath.Abs("../" + pathToFile)
		actual, _ := deserializeTrackInput(absPath)
		assert.Equal(t, expected, actual)
	}
}
