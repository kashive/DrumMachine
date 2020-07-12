package main

//represents sound with name and path to the audio file
type Sound struct {
	Name       string //hiHat, bass drum etc.
	PathToFile string // path to the sound file
}

//represents the json input used as the program input
type TrackInput struct {
	Name          string //name of the track
	NumberOfSteps int //number of steps in the track
	PlaybackConfig PlaybackConfig //bpm and whether to continuously loop
	Patterns      []Pattern //sound name and patterns
}

//configuration regarding playing the track
type PlaybackConfig struct {
	Bpm  int  //bpm to play the song
	Loop bool //whether to loop
}

//represents a track with all the steps and the sounds in those steps
type Track struct {
	Name  string //name of the track
	Steps []Step //steps struct
	TrackInput TrackInput //input used to create this track
}

//represents the stepId and the sounds in that step
type Step struct {
	Id     int     //id identifying the step
	Sounds []Sound //sounds to play in this step
}

//represents the sound pattern
type Pattern struct {
	Sound Sound //sound struct
	Steps []int //sequence pattern in total number of steps
}
