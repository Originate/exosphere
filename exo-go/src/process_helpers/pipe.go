package processHelpers

import (
	"fmt"
	"io"
)

func splitReader(reader io.ReadCloser) (io.ReadCloser, io.ReadCloser) {
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
