package main

import (
	"flag"
	"os"
	"time"

	"github.com/unknwon/log"
	"github.com/volkstrader/fscrawler/crawler"
)

func main() {
	start := time.Now()
	rootDir := flag.String("r", ".", "directory to crawl")
	maxCrawler := flag.Int("c", 20, "maximum number crawlers")
	flag.Parse()
	*rootDir = os.ExpandEnv(*rootDir)
	log.Info("Crawling %v with %v crawlers", *rootDir, *maxCrawler)
	crawler.Crawl(*rootDir, *maxCrawler)
	log.Info("fscrawler took %s", time.Since(start))
}
