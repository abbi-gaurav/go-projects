package queuing

import (
	"bufio"
	"github.com/abbi-gaurav/go-learning-projects/ultimate-go-programming/ccp-in-go/patterns/utils"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func BenchmarkUnbufferedWrite(b *testing.B) {
	performWrite(b, tempFileOrFatal())
}

func BenchmarkBufferedWrite(b *testing.B) {
	bufferedFile := bufio.NewWriter(tempFileOrFatal())
	performWrite(b, bufferedFile)
}

func tempFileOrFatal() *os.File {
	file, err := ioutil.TempFile("", "tmp")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return file
}

func performWrite(b *testing.B, writer io.Writer) {
	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()

	for bt := range utils.Take(done, utils.Repeat(done, byte(0)), b.N) {
		_, err := writer.Write([]byte{bt.(byte)})

		if err != nil {
			log.Fatalf("error write: %v", err)
		}
	}
}
