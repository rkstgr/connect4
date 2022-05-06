package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

type TestBoardPosition struct {
	positionString string
	evaluation     int
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

// This is a simple test program to test the connect4 package.

// Load the test data from the file "test_data/End-Easy.txt"
// The file contains the following:
// <positionString> <evaluation>

func loadTestData(file string) []TestBoardPosition {
	var testData []TestBoardPosition
	fileHandle, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer func(fileHandle *os.File) {
		err := fileHandle.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(fileHandle)
	scanner := bufio.NewScanner(fileHandle)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		testData = append(testData, TestBoardPosition{fields[0], atoi(fields[1])})
	}
	return testData
}

func evaluateTestData(testData []TestBoardPosition, maxEvalTime time.Duration, t *testing.T) {
	positionCounters := make([]int, len(testData))
	elapsedTimes := make([]time.Duration, len(testData))
	for i, testPosition := range testData {
		position := createBoard(testPosition.positionString)
		positionsCount := Counter{}
		// start a goroutine that with a timer of maxEvalTime
		// if the evaluation is not finished in time, the test fails
		evaluated := make(chan bool)
		go func() {
			select {
			case <-time.After(maxEvalTime):
				t.Errorf("Test timed out")
				os.Exit(1)
			case <-evaluated:
				return
			}
		}()

		startTime := time.Now()
		evaluation := solve(position, &positionsCount)
		elapsedTime := time.Since(startTime)
		evaluated <- true

		elapsedTimes[i] = elapsedTime
		count := positionsCount.get()
		positionCounters[i] = count
		if evaluation != testPosition.evaluation {
			t.Errorf("❌ Evaluation of %s is %d, but should be %d",
				testPosition.positionString, evaluation, testPosition.evaluation)
		}
	}
	// Print average positions count
	var totalPositions int
	for _, count := range positionCounters {
		totalPositions += count
	}
	var totalTime time.Duration
	for _, elapsedTime := range elapsedTimes {
		totalTime += elapsedTime
	}
	var averageTime time.Duration
	if len(elapsedTimes) > 0 {
		averageTime = totalTime / time.Duration(len(elapsedTimes))
	}
	average := float64(totalPositions) / float64(len(positionCounters))
	t.Logf("Mean elapsed time: %s", averageTime)
	t.Logf("Average positions count: %.2f", average)
	t.Logf("K Pos. per second: %.2f", float64(totalPositions)/1000/totalTime.Seconds())
}

// Create a test for every test file:
// - End-Easy.txt
// - Middle-Easy.txt
// - Middle-Medium.txt
// - Start-Easy.txt
// - Start-Medium.txt
// - Start-Hard.txt

func TestEndEasy(t *testing.T) {
	testDataEndEasy := loadTestData("test_data/End-Easy.txt")
	evaluateTestData(testDataEndEasy, time.Second, t)
}

func TestMiddleEasy(t *testing.T) {
	testData := loadTestData("test_data/Middle-Easy.txt")
	evaluateTestData(testData, 30*time.Second, t)
}

func TestMiddleMedium(t *testing.T) {
	testData := loadTestData("test_data/Middle-Medium.txt")
	evaluateTestData(testData, time.Minute, t)
}

func TestStartEasy(t *testing.T) {
	testData := loadTestData("test_data/Start-Easy.txt")
	evaluateTestData(testData, time.Minute, t)
}

func TestStartMedium(t *testing.T) {
	testData := loadTestData("test_data/Start-Medium.txt")
	evaluateTestData(testData, time.Minute, t)
}

func TestStartHard(t *testing.T) {
	testData := loadTestData("test_data/Start-Hard.txt")
	evaluateTestData(testData, time.Minute, t)
}
