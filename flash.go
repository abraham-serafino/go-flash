package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)

type options struct {
	fileName string
}

type flashcard struct {
	front string
	back  string
}

func getOptions() options {
	fileName1 := ""
	fileName2 := ""
	opts := options{fileName: "example.crd"}

	printUsageAndExit := func() {
		fmt.Println("Usage: flash {OPTIONS}")
		fmt.Println("")
		fmt.Println("    Options:")
		fmt.Println("\t-f, --filename [filename]\tpath to flashcard file")
		fmt.Println("")
		fmt.Println("    ex. flash --filename jlpt-3-kanji.crd")
		fmt.Println("")

		os.Exit(1)
	}

	flag.StringVar(&fileName1, "f", "", "Path to flashcard file")
	flag.StringVar(&fileName2, "filename", "", "Path to flashcard file")

	flag.Parse()

	if fileName1 != "" && fileName2 != "" {
		printUsageAndExit()
	}

	if fileName1 != "" {
		opts.fileName = fileName1

	} else if fileName2 != "" {
		opts.fileName = fileName2
	}

	return opts
}

func trim(str string) string {
	return strings.TrimLeft(
		strings.TrimLeft(
			strings.TrimRight(
				strings.TrimRight(str, "\r"),
				"\n"),
			"\r"),
		"\n")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadCardsFromFile(fileName string) []flashcard {

	fileFormatError := func() {
		fmt.Println(fileName, ": format error. File must have the following format:")
		fmt.Println("")
		fmt.Println("\t[front side text]")
		fmt.Println("\t%")
		fmt.Println("\t[back side text]")
		fmt.Println("\t%")
		fmt.Println("\t[front side text]")
		fmt.Println("\t&c.")
		os.Exit(1)
	}

	data, err := os.ReadFile(fileName)
	check(err)
	str := string(data)

	splitLines := strings.Split(str, "%")
	numberOfLines := len(splitLines)

	getFrontAndBack := func(frontIndex int, backIndex int) (string, string) {
		if frontIndex > numberOfLines || backIndex > numberOfLines {
			fileFormatError()
		}

		return trim(splitLines[frontIndex]), trim(splitLines[backIndex])
	}

	if numberOfLines%2 != 0 {
		fileFormatError()
	}

	var flashcards []flashcard

	for i := 0; i < len(splitLines); i += 2 {
		frontIndex := i
		backIndex := i + 1

		card := flashcard{front: "", back: ""}

		card.front, card.back = getFrontAndBack(frontIndex, backIndex)
		flashcards = append(flashcards, card)
	}

	return flashcards
}

func testCard(card flashcard) bool {
	fmt.Println(card.front)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	return trim(input) == trim(card.back)
}

func main() {
	options := getOptions()

	flashcards := loadCardsFromFile(options.fileName)

	shuffle := func() {
		rand.Shuffle(len(flashcards), func(i, j int) {
			flashcards[i], flashcards[j] = flashcards[j], flashcards[i]
		})
	}

	for len(flashcards) > 0 {
		shuffle()

		for i := 0; i < len(flashcards); i++ {
			if testCard(flashcards[i]) {
				fmt.Println("Right!")
				fmt.Println("")
				flashcards = append(flashcards[:i], flashcards[i+1:]...)

			} else {
				fmt.Println(flashcards[i].back + "\n")
			}
		}
	}
}
