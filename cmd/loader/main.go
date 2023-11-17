package main

import (
	"bufio"
	_ "embed"
	"log/slog"
	"os"
	"strings"
	"unicode"
)

//go:generate curl https://www.gutenberg.org/cache/epub/50133/pg50133.txt -o dunwich.txt

//go:embed dunwich.txt
var text string

var start_marker string = "*** START OF THE PROJECT GUTENBERG EBOOK THE DUNWICH HORROR ***"
var end_marker string = "*** END OF THE PROJECT GUTENBERG EBOOK THE DUNWICH HORROR ***"

func isChapterMarker(line string) bool {
	if line == "" {
		return false
	}

	for _, v := range line {
		if !unicode.IsDigit(v) {
			return false
		}
	}

	return true
}

func main() {
	slog.Info("Hi")

	_, int0, _ := strings.Cut(text, start_marker)
	int1, _, _ := strings.Cut(int0, end_marker)
	int2 := strings.ReplaceAll(int1, "\r\n", "\n")
	_, int3, _ := strings.Cut(int2, "1\n")

	chunks := strings.Split(int3, "\n")

	corpus := []string{}

	paragraph := ""

	for i, line := range chunks {
		if line == "" {
			slog.Debug(
				"line is empty",
				slog.Int("num", i),
			)

			if paragraph != "" {
				corpus = append(corpus, paragraph)
				paragraph = ""
			}
		} else if isChapterMarker(line) {
			slog.Info(
				"line is a chapter marker, ignoring",
				slog.Int("num", i),
				slog.String("content", line),
			)
		}

		// Not an empty line or chapter marker, so append the line to the current paragraph
		paragraph = paragraph + strings.TrimSpace(line) + "\n"
	}

	f, err := os.Create("corpus.txt")
	if err != nil {
		slog.Info(
			"failed to create output file",
		)
		slog.Error(
			"failed to create output file",
			slog.String("err", err.Error()),
		)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(strings.Join(corpus, "======"))
	w.Flush()

}
