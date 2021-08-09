package omx

import (
	"bufio"
	"os"
)

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func ReadFileLines(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	text := make([]string, 0)
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	return text
}