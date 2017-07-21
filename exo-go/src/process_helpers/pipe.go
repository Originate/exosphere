package processHelpers

import (
	"bufio"
	"fmt"
	"io"
)

func duplicateReader(reader io.ReadCloser) (io.ReadCloser, io.ReadCloser) {
	reader1, writer1 := io.Pipe()
	reader2, writer2 := io.Pipe()
	multiWriter := io.MultiWriter(writer1, writer2)
	go func() {
		if _, err := io.Copy(multiWriter, reader); err != nil {
			fmt.Println("Error copying stream", err)
		}
	}()
	return reader1, reader2
}

func readPipe(reader io.Reader, log func(string)) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log(scanner.Text())
	}
}
