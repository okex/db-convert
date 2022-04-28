package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	tmdb "github.com/tendermint/tm-db"
)

var wg sync.WaitGroup

func main() {
	fmt.Println(os.Args)

	sdbType := os.Args[1]
	ddbType := os.Args[2]
	sdbDir := os.Args[3]
	ddbDir := os.Args[4]

	// */data/*
	fs, err := ioutil.ReadDir(sdbDir)
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		if f.IsDir() && strings.HasSuffix(f.Name(), ".db") {
			name := strings.Split(f.Name(), ".")[0]
			wg.Add(1)
			go func(name, leveldbDir, rocksdbDir string) {
				defer wg.Done()

				log.Printf("convert %s start...\n", name)
				sdb, err := tmdb.NewDB(name, tmdb.BackendType(sdbType), sdbDir)
				if err != nil {
					panic(err)
				}

				if _, err := os.Stat(ddbDir); os.IsNotExist(err) {
					if err := os.MkdirAll(ddbDir, 0755); err != nil {
						panic(err)
					}
				}
				ddb, err := tmdb.NewDB(name, tmdb.BackendType(ddbType), ddbDir)
				if err != nil {
					panic(err)
				}

				convert(sdb, ddb)
				log.Printf("convert %s end.\n", name)
				//log.Printf("compact %s start...\n", name)
				//rdb.DB().CompactRange(gorocksdb.Range{})
				//log.Printf("compact %s end.\n", name)
			}(name, sdbDir, ddbDir)
		}
	}

	wg.Wait()
}

func convert(sdb, ddb tmdb.DB) {
	iter, err := sdb.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}

	for ; iter.Valid(); iter.Next() {
		ddb.Set(iter.Key(), iter.Value())
	}
	iter.Close()
}
