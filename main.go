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

// Create a wait group
var waitGroup sync.WaitGroup
var letterCounter = make(map[string]int)

func main() {
	port := 3000
	address := fmt.Sprintf(":%d", port)

	app := fiber.New()

	app.Post("/count-letters", CounterLettersHandler)

	log.Fatal(app.Listen(address))
}

func CounterLettersHandler(context *fiber.Ctx) error {
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
	contentInLowerCase := strings.ToLower(contentWithoutSpaces)

	// Split the content into paragraphs
	arrayParagraphs := strings.Split(contentInLowerCase, "\n")

	// mutex
	var mutex sync.Mutex

	// Add the number of paragraphs to the wait group
	waitGroup.Add(len(arrayParagraphs))

	for i := 0; i < len(arrayParagraphs); i++ {
		// Count the letters in the paragraph
		go CountLetters(arrayParagraphs[i], &mutex)
	}

	waitGroup.Wait()

	return context.Status(http.StatusOK).JSON(letterCounter)
}

func CountLetters(paragraph string, m *sync.Mutex) {
	defer waitGroup.Done()

	runeRegex := regexp.MustCompile("[a-z]")

	letters := runeRegex.FindAllString(paragraph, -1)

	for _, letter := range letters {
		m.Lock()
		letterCounter[letter]++
		m.Unlock()
	}
}
