package dunwich

import (
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"unicode"
)

//go:generate curl https://www.gutenberg.org/cache/epub/50133/pg50133.txt -o dunwich.txt
//go:embed dunwich.txt
var text string

var start_marker string = "*** START OF THE PROJECT GUTENBERG EBOOK THE DUNWICH HORROR ***"
var end_marker string = "*** END OF THE PROJECT GUTENBERG EBOOK THE DUNWICH HORROR ***"

var corpus []string

func init() {
	buildCorpus(&corpus)
}

func buildCorpus(target *[]string) {

	slog.Debug("starting to build book")

	_, int0, _ := strings.Cut(text, start_marker)
	int1, _, _ := strings.Cut(int0, end_marker)
	int2 := strings.ReplaceAll(int1, "\r\n", "\n")
	_, int3, _ := strings.Cut(int2, "1\n")

	chunks := strings.Split(int3, "\n")

	paragraph := ""

	for i, line := range chunks {
		if line == "" {
			slog.Debug(
				"line is empty",
				slog.Int("num", i),
			)

			if paragraph != "" {
				*target = append(*target, paragraph)
				paragraph = ""
			}
		} else if isChapterMarker(line) {
			slog.Debug(
				"line is a chapter marker, ignoring",
				slog.Int("num", i),
				slog.String("content", line),
			)
		} else {
			// Not an empty line or chapter marker, so append the line to the current paragraph
			paragraph = paragraph + strings.TrimSpace(line) + "\n"
		}
	}

	slog.Debug("finished building book")
}

func isChapterMarker(line string) bool {

	line = strings.TrimSpace(line)

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

func GetNumSegments() int {
	return len(corpus)
}

func GetSegment(id int) (*string, error) {
	if id < 0 || id >= len(corpus) {
		return nil, errors.New(fmt.Sprintf("index %d is not in range 0:%d", id, len(corpus)))
	}

	return &corpus[id], nil
}
