package crawler

import (
	"io/ioutil"
	"sync"

	"path/filepath"

	"os"

	"crypto/md5"
	"io"

	"encoding/hex"

	"time"

	"github.com/unknwon/log"
)

type fileNode struct {
	path string
	size int64
	hash string
}

var dirCh chan *fileNode
var filesCh chan *fileNode

var nDir, nFile int
var totalSize int64

func monitor(wg *sync.WaitGroup) {
	wg.Wait()
	close(dirCh)
}

func Crawl(root string, n int) {
	wg := &sync.WaitGroup{}
	dirCh = make(chan *fileNode, n)
	filesCh = make(chan *fileNode, n)

	dirCh <- &fileNode{
		path: root,
	}

	go monitor(wg)

loop:
	for {
		select {
		case subdir, ok := <-dirCh:
			if !ok {
				log.Info("crawled %v directories and %v files, total size %v", nDir-1, nFile, totalSize)
				break loop
			}
			nDir++
			wg.Add(1)
			go func() {
				defer wg.Done()
				crawl(subdir.path)
			}()
		case f := <-filesCh:
			nFile++
			totalSize += f.size
			go func(f *fileNode) {
				wg.Add(1)
				defer wg.Done()

				start := time.Now()
				log.Debug("hashing(%s) ...", f.path)
				f.hash = hex.EncodeToString(hashFile(f.path))
				log.Debug("hash(%s) %s in %s", f.path, f.hash, time.Since(start))
			}(f)
		}
	}
}

func crawl(dir string) {
	if dir == "" || dir == "." || dir == ".." {
		return
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Error("%v", err)
		return
	}

	for _, file := range files {
		node := fileNode{
			path: filepath.Join(dir, file.Name()),
		}

		if file.IsDir() {
			dirCh <- &node
		} else {
			node.size = file.Size()
			filesCh <- &node
		}

		log.Debug("%v", node)
	}
}

func hashFile(file string) []byte {
	f, err := os.Open(file)
	if err != nil {
		log.Error("%s", err)
		return nil
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Error("%s", err)
		return nil
	}

	return h.Sum(nil)
}
