package search

// I'm Shodikhuja :)

import (
	"bufio"
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

//Result описывает один результат поиска
type Result struct {
	// фраза, которую искали
	Phrase string
	// целиком вся строкка, где нашли фразу (без \n или \r\n в конце)
	Line string
	// номер строки (начиная с 1)
	LineNum int64
	// номер позиции (начиная с 1) - или символа?
	ColNum int64
}

//All ищет все появления фразы в текстовых файлах files
func All(ctx context.Context, phrase string, files []string) <-chan []Result {
	ch := make(chan []Result)
	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(ctx)

	for i := 0; i < len(files); i++ {
		wg.Add(1)

		go func(ctx context.Context, path string, i int, ch chan<- []Result) {
			defer wg.Done()

			res, err := FindAll(phrase, path)
			if err != nil {
				log.Println("error not opened file err => ", err)
				return
			}

			if len(res) > 0 {
				ch <- res
			}

		}(ctx, files[i], i, ch)
	}

	go func() {
		defer close(ch)
		wg.Wait()

	}()

	cancel()
	return ch
}

//Any ...
func Any(ctx context.Context, phrase string, files []string) <-chan Result {
	ch := make(chan Result)
	wg := sync.WaitGroup{}
	result := Result{}

	ctx, cancel := context.WithCancel(ctx)

	for _, f := range files {
		file, err := ioutil.ReadFile(f)
		if err != nil {
			log.Println("error while open file: ", err)
		}

		if strings.Contains(string(file), phrase) {
			res, err := FindAny(phrase, string(file))
			if err != nil {
				log.Println("error while open file: ", err)
			}

			if (Result{}) != res {
				result = res
				break
			}
		}

	}

	wg.Add(1)
	go func(ctx context.Context, ch chan<- Result) {
		defer wg.Done()
		if (Result{}) != result {
			ch <- result
		}
	}(ctx, ch)

	go func() {
		defer close(ch)
		wg.Wait()

	}()

	cancel()
	return ch
}

//FindAll ...
func FindAll(phrase, path string) ([]Result, error) {
	res := []Result{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for i := 0; i < len(lines); i++ {

		if strings.Contains(lines[i], phrase) {

			r := Result{
				Phrase:  phrase,
				Line:    lines[i],
				LineNum: int64(i + 1),
				ColNum:  int64(strings.Index(lines[i], phrase)) + 1,
			}

			res = append(res, r)
		}
	}
	return res, nil
}

//FindAny ...
func FindAny(phrase, path string) (Result, error) {
	var lines []string
	res := Result{}
	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], phrase) {
			res = Result{
				Phrase:  phrase,
				Line:    lines[i],
				LineNum: int64(i + 1),
				ColNum:  int64(strings.Index(lines[i], phrase)) + 1,
			}
		}
	}
	return res, nil
}