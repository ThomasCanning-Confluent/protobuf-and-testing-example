package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
)

var points map[string]int

func askForPredictions(data []byte) []byte {
	reader := bufio.NewReader(os.Stdin)
	matchInstance := &Match{}

	err := proto.Unmarshal(data, matchInstance)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
		return nil
	}

	if matchInstance.Predictions == nil {
		matchInstance.Predictions = make(map[string]string)
	}

	for {
		fmt.Println("Enter participant name (enter exit to finish)")
		playerName, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading participant name input")
			return nil
		}
		playerName = strings.TrimSpace(playerName)
		if playerName == "exit" {
			fmt.Println("Finished entering participants")
			break
		}

		fmt.Println("Enter participant's score prediction in form of 'team1:team2'")
		playerGuess, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading participant's input")
			return nil
		}
		playerGuess = strings.TrimSpace(playerGuess)

		matchInstance.Predictions[playerName] = playerGuess

		data, err = proto.Marshal(matchInstance)
		if err != nil {
			log.Fatal("marshaling error: ", err)
			return nil
		}
	}
	return data
}

func askForActualScore(data []byte) []byte {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter actual score in form of 'team1:team2'")
	actualScore, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading actual score input")
		return nil
	}
	matchInstance := &Match{}
	err = proto.Unmarshal(data, matchInstance)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
		return nil
	}
	matchInstance.ActualScore = strings.TrimSpace(actualScore)

	// Serialize match instance to bytes
	data, err = proto.Marshal(matchInstance)
	if err != nil {
		log.Fatal("marshaling error: ", err)
		return nil
	}

	return data
}

func calculatePoints(guess string, score string) int {
	guessSplit := strings.Split(guess, ":")
	scoreSplit := strings.Split(score, ":")

	if guessSplit[0] == scoreSplit[0] && guessSplit[1] == scoreSplit[1] {
		return 3 // exact score
	} else if guessSplit[0] == scoreSplit[0] || guessSplit[1] == scoreSplit[1] {
		return 1 // correct result for draw
	} else if (guessSplit[0] > guessSplit[1] && scoreSplit[0] > scoreSplit[1]) || (guessSplit[0] < guessSplit[1] && scoreSplit[0] < scoreSplit[1]) {
		return 1 // correct result for team 1 win
	} else {
		return 0
	}
}

func main() {
	points = make(map[string]int)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Enter match data? (y/n)")
		cont, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading cont input")
			return
		}
		if strings.TrimSpace(cont) == "n" {
			fmt.Println("Final points are: ")
			for player, point := range points {
				fmt.Println(player, point)
			}
			break
		}

		matchInstance := &Match{
			Predictions: make(map[string]string),
		}
		data, err := proto.Marshal(matchInstance)
		data = askForPredictions(data)
		data = askForActualScore(data)

		// Deserialize data back to matchInstance
		var updatedMatch Match
		err = proto.Unmarshal(data, &updatedMatch)
		if err != nil {
			log.Fatal("unmarshaling error: ", err)
			return
		}

		// Use updatedMatch.ActualScore for further processing
		for player, guess := range updatedMatch.Predictions {
			points[player] += calculatePoints(guess, updatedMatch.ActualScore)
		}

		fmt.Println("The current points are: ")
		for player, pointCount := range points {
			fmt.Println(player, " has a score of ", pointCount)
		}
	}
}
