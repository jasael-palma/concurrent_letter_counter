package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
)

func main() {
	port := 3000
	address := fmt.Sprintf(":%d", port)

	app := fiber.New()

	app.Post("/count-letters", CounterLettersHandler)

	log.Fatal(app.Listen(address))
}

func CounterLettersHandler(context *fiber.Ctx) error {
	// Create a wait group
	var waitGroup sync.WaitGroup
	letterCounter := map[string]int{
		"A": 0,
		"B": 0,
		"C": 0,
		"D": 0,
		"E": 0,
		"F": 0,
		"G": 0,
		"H": 0,
		"I": 0,
		"J": 0,
		"K": 0,
		"L": 0,
		"M": 0,
		"N": 0,
		"O": 0,
		"P": 0,
		"Q": 0,
		"R": 0,
		"S": 0,
		"T": 0,
		"U": 0,
		"V": 0,
		"W": 0,
		"X": 0,
		"Y": 0,
		"Z": 0,
	}

	context.Accepts("multipart")

	// Parse the multipart form
	file, error := context.FormFile("file")

	if error != nil {
		return fiber.NewError(500, error.Error())
	}

	if file.Header.Get("Content-Type") != "text/plain" {
		return fiber.NewError(400, "El archivo debe ser de tipo text/plain")
	}

	if file.Filename[len(file.Filename)-4:] != ".txt" {
		return fiber.NewError(400, "El archivo debe tener la extensión .txt")
	}

	// Open the file
	fileOpened, error := file.Open()

	if error != nil {
		return fiber.NewError(500, error.Error())
	}

	// Close the file when we finish
	defer fileOpened.Close()

	buffer, error := io.ReadAll(fileOpened)

	if error != nil {
		return fiber.NewError(500, error.Error())
	}

	content := string(buffer)

	contentWithoutSpaces := strings.ReplaceAll(content, " ", "")

	// Split the content into paragraphs
	arrayParagraphs := strings.Split(contentWithoutSpaces, "\n")

	// mutex
	var mutex sync.Mutex

	// Add the number of paragraphs to the wait group
	waitGroup.Add(len(arrayParagraphs))

	for i := 0; i < len(arrayParagraphs); i++ {
		// Count the letters in the paragraph
		go CountLetters(arrayParagraphs[i], &mutex, &waitGroup, letterCounter)
	}

	waitGroup.Wait()

	return context.Status(http.StatusOK).JSON(letterCounter)
}

func CountLetters(paragraph string, m *sync.Mutex, waitGroup *sync.WaitGroup, letterCounter map[string]int) {
	defer waitGroup.Done()

	runeRegex := regexp.MustCompile("(?i)[a-z]")

	letters := runeRegex.FindAllString(paragraph, -1)

	for _, letter := range letters {
		m.Lock()
		letterCounter[strings.ToUpper(letter)]++
		m.Unlock()
	}
}
