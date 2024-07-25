package chat

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
)

const (
	adjectives  = "files/adjectives.txt"
	celebrities = "files/celebrities.txt"
	fantasies   = "files/fantasies.txt"
	foods       = "files/foods.txt"
	objects     = "files/objects.txt"
	profissions = "files/profissions.txt"

	namesDelimiter = ";"
)

type NameGenerator struct {
	adjectives    []string
	miscellaneous []string
	logger        *zap.SugaredLogger
}

func NewNameGenerator(logger *zap.SugaredLogger) *NameGenerator {
	files := map[string][]string{
		adjectives:  {},
		celebrities: {},
		fantasies:   {},
		foods:       {},
		objects:     {},
		profissions: {},
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(len(files))
	for path := range files {
		go readFiles(path, files, &wg, &mu, logger)
	}

	wg.Wait()

	return &NameGenerator{
		logger:     logger,
		adjectives: files[adjectives],
		miscellaneous: merge(
			files[celebrities], files[fantasies], files[foods], files[objects], files[profissions]),
	}
}

func readFiles(
	path string,
	files map[string][]string,
	wg *sync.WaitGroup,
	mu *sync.Mutex,
	logger *zap.SugaredLogger,
) {
	defer wg.Done()
	words, err := readFile(path)
	if err != nil {
		logger.Errorw("failed to read file", "file", path, zap.Error(err))
	}

	mu.Lock()
	files[path] = strings.Split(words, namesDelimiter)
	mu.Unlock()
}

func merge(lists ...[]string) []string {
	var result []string
	for _, list := range lists {
		result = append(result, list...)
	}
	return result
}

func readFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("erro ao abrir arquivo %s", path)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("erro ao ler arquivo %s", path)
	}

	return string(data), nil
}

func random(list []string) string {
	randomIndex := rand.Intn(len(list))
	return list[randomIndex]
}

func (n *NameGenerator) Generate() string {
	return fmt.Sprint(random(n.adjectives), random(n.miscellaneous))
}
