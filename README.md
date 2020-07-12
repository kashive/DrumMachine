![Demo](https://github.com/kashive/DrumMachine/blob/master/demo.gif)
#### Features
* Animated visualization on the console
* Audio output as well
* Highly configurable
    * can change track name, bpm
    * can easily add new sounds and patterns
    * supports continuous playback
    
#### How to run
* Compile the package and run `drumMachine.go` with one program argument i.e. relative path to the input json file. eg: `/sampleInput/sample_input.json`

#### Design
* The program start code is at `drumMachine.go`
* Uses JSON file as input. Used JSON because Go has a built in support for unmarshalling JSON to struct. Also, JSON allows easy addition of attributes in case we wanted to expose more functionality to the user.
* The `drumMachine.go` is responsible for de-serializing, validating the input, building track struct and invoking the player.
* All the domain structs can be found in the `types.go` file.
* `Track` struct represents the sounds the player needs to play at each step. It is denormalized in order to make the abstractions clearer and the code simpler. The trade off being higher memory usage. 
* The `Player.go` outputs to both the console and the speaker.
* To print on the console player uses `github.com/jedib0t/go-pretty/table` and for audio out it uses `github.com/buger/goterm` library. The footer of the table is used to indicate the step that the player is on.
* On every step the entire table is printed on the console to give an animation kind of effect.
* `player.go` uses the `github.com/faiface/beep` library to play the audio. I found the `portaudio` library a bit low level. So, used beep instead.
* There are unit tests around input parsing and track building which can be found at `inputParse_test.go`

#### Assumptions/Possible Next Steps
* The player when playing the audio file creates a new stream and closes it on every step which adds to the latency. Can optimize by re-using the same stream and re-setting it to the initial position.
* The player assumes a sampleRate of 44.1 kHz. All the audio samples in the examples are at 44.1 kHz. Beep library does support re-sampling.
* There is a latency of around 230-400 ms when processing each step in the player which limits the max BPM the program can achieve. The above mentioned idea of reusing the same stream may help.
* In the input.json the pattern struct is a list of sound. If a sound repeats then the last one wins.
* The code for console player and audio player is in the same file. It would be ideal to separate them and somehow still be able to run them simultaneously.
* We can provide different ways to define a pattern in the input file. One idea is to provide a `startAtStepId` and `repeatAfterSteps`. This would be especially useful if we needed to create a pattern across many steps.
