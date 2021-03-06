package rules

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/romainmenke/pusher/common"
)

type settings struct {
	path        string
	rulesReader io.Reader
}

// Option -> Functional Options : Awesome
type Option func(*settings)

// FileOption is used to pass the location of your Rules File to the Handler func
func FileOption(path string) func(*settings) {
	return func(s *settings) {
		s.path = path
	}
}

// ReaderOption is used to pass an io.Reader which contains your rules to the Handler func
// This will take precendent over the File Option
func ReaderOption(reader io.Reader) func(*settings) {
	return func(s *settings) {
		s.rulesReader = reader
	}
}

// Handler wraps an http.Handler.
// It sets Link Headers for paths declared in the Rules File.
func Handler(handler http.Handler, options ...Option) http.Handler {

	s := &settings{}
	m := &sync.RWMutex{}

	for _, option := range options {
		option(s)
	}

	if s.path == "" && s.rulesReader == nil {
		return handler
	}

	var (
		reader io.Reader = s.rulesReader
		file   *os.File
		err    error
	)

	if s.path != "" && reader == nil {
		file, err = os.Open(s.path)
		if err != nil {
			return handler
		}
		reader = file
	}

	if file != nil {
		defer file.Close()
	}

	linkMap, assetMap, err := read(reader)
	if err != nil {
		return handler
	}

	mux := http.NewServeMux()

	for path := range linkMap {
		scopedPath := path
		h := wrap(scopedPath, assetMap, linkMap, m, handler)
		mux.Handle(scopedPath, h)
	}

	return mux
}

func read(rules io.Reader) (map[string][]string, map[string]struct{}, error) {

	pathMap := make(map[string][]string)
	headerMap := make(map[string]struct{})

	scanner := bufio.NewScanner(rules)

	currentPath := ""
	currentHeaders := []string{}

	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) > 0 && txt[:1] == "/" && currentPath == "" {
			currentHeaders = []string{}
			currentPath = txt
			continue
		}
		if txt == "-" {
			continue
		}
		if currentPath != "" && len(txt) > 0 {
			link := common.ParseLinkHeader(txt)
			if link == "" {
				continue
			}
			currentHeaders = append(currentHeaders, txt)
			headerMap[link] = struct{}{}
			continue
		}
		if txt == "" {
			if currentPath != "" {
				pathMap[currentPath] = currentHeaders
			}
			currentPath = ""
			currentHeaders = []string{}
		}

	}

	if currentPath != "" {
		pathMap[currentPath] = currentHeaders
	}

	if err := scanner.Err(); err != nil {
		if err == io.EOF {
			return pathMap, headerMap, nil
		}
		return nil, nil, err
	}

	return pathMap, headerMap, nil

}
