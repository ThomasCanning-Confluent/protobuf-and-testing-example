package main

import (
	"github.com/golang/protobuf/proto"
	"os"
	"testing"
)

func Test_calculatePoints(t *testing.T) {
	//calculatePoints(guess string, score string)
	type args struct {
		guess string
		score string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "incorrect", args: args{guess: "3:1", score: "1:3"}, want: 0},
		{name: "winner correct", args: args{guess: "3:1", score: "2:1"}, want: 1},
		{name: "score correct", args: args{guess: "3:1", score: "3:1"}, want: 3},
		{name: "incorrect2", args: args{guess: "2:4", score: "4:2"}, want: 0},
		{name: "winner correct2", args: args{guess: "2:3", score: "2:4"}, want: 1},
		{name: "score correct2", args: args{guess: "0:0", score: "0:0"}, want: 3},
		{name: "score correct3", args: args{guess: "1:2", score: "1:2"}, want: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculatePoints(tt.args.guess, tt.args.score); got != tt.want {
				t.Errorf("calculatePoints with prediction %v and actual score %v returned %v points, wanted %v points", tt.args.guess, tt.args.score, got, tt.want)
			}
		})
	}
}

func TestAskForPredictions(t *testing.T) {
	// Create a pipe to simulate os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	// Save the original os.Stdin
	origStdin := os.Stdin
	defer func() { os.Stdin = origStdin }()

	// Redirect os.Stdin to the pipe
	os.Stdin = r

	// Write test input to the pipe
	testInput := "Tom C\n1:0\nElliot B\n2:1\nexit\n"

	// Uses a goroutine to write to the pipe asynchronously
	go func() {
		w.Write([]byte(testInput))
		w.Close()
	}()

	// Prepare a Match instance and call the method
	match := &Match{
		Predictions: make(map[string]string)}
	data, err := proto.Marshal(match)
	data = askForPredictions(data)
	updatedMatch := &Match{}
	err = proto.Unmarshal(data, updatedMatch)

	// Validate the results
	expectedPredictions := map[string]string{
		"Tom C":    "1:0",
		"Elliot B": "2:1",
	}

	for name, score := range expectedPredictions {
		if updatedMatch.Predictions[name] != score {
			t.Errorf("Expected %s to predict %s, got %s", name, score, updatedMatch.Predictions[name])
		}
	}

	if len(updatedMatch.Predictions) != len(expectedPredictions) {
		t.Errorf("Expected %d predictions, got %d", len(expectedPredictions), len(updatedMatch.Predictions))
	}
}

func TestAskForActualScore(t *testing.T) {
	// Create a pipe to simulate os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	// Save the original os.Stdin
	origStdin := os.Stdin
	defer func() { os.Stdin = origStdin }()

	// Redirect os.Stdin to the pipe
	os.Stdin = r

	// Write test input to the pipe
	testInput := "3:2\n"

	// Uses a goroutine to write to the pipe asynchronously
	go func() {
		w.Write([]byte(testInput))
		w.Close()
	}()

	// Prepare a Match instance and call the method
	match := &Match{}
	data, err := proto.Marshal(match)
	data = askForActualScore(data)
	updatedMatch := &Match{}
	err = proto.Unmarshal(data, updatedMatch)

	// Validate the results
	if updatedMatch.ActualScore != "3:2" {
		t.Errorf("Expected actual score to be 3:2, got %s", updatedMatch.ActualScore)
	}
}
