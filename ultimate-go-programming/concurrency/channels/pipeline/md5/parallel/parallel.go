package parallel

import (
	"crypto/md5"
	"sync"
	"path/filepath"
	"os"
	"io/ioutil"
	"errors"
)

type result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

func sumFiles(done <-chan struct{}, root string) (<-chan result, <-chan error) {
	c := make(chan result)
	errc := make(chan error, 1)

	go func() {
		var wg sync.WaitGroup

		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.Mode().IsRegular() {
				return nil
			}

			wg.Add(1)

			go func() {
				data, readErr := ioutil.ReadFile(path)
				select {
				case c <- result{path: path, sum: md5.Sum(data), err: readErr}:
				case <-done:
				}
				wg.Done()
			}()

			select {
			case <-done:
				return errors.New("walk cancelled")
			default:
				return nil
			}
		})

		go func() {
			wg.Wait()
			close(c)
		}()

		errc <- err
	}()

	return c, errc
}

func MD5All(root string) (map[string][md5.Size]byte, error) {
	done := make(chan struct{})
	defer close(done)

	c, errc := sumFiles(done, root)

	m := make(map[string][md5.Size]byte)
	for r := range c {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}

	if err := <-errc; err != nil {
		return nil, err
	}

	return m, nil
}
