package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/tecbot/gorocksdb"
	tmdb "github.com/tendermint/tm-db"
)

var wg sync.WaitGroup

func main() {
	fmt.Println(os.Args[1:])

	// /data/
	fs, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		if f.IsDir() {
			name := strings.Split(f.Name(), ".")[0]
			wg.Add(1)
			go func(name, leveldbDir, rocksdbDir string) {
				defer wg.Done()

				log.Printf("convert %s start...\n", name)
				ldb, err := tmdb.NewGoLevelDB(name, leveldbDir)
				if err != nil {
					panic(err)
				}
				rdb, err := tmdb.NewRocksDB(name, rocksdbDir)
				if err != nil {
					panic(err)
				}

				iter, err := ldb.Iterator(nil, nil)
				if err != nil {
					panic(err)
				}

				for ; iter.Valid(); iter.Next() {
					rdb.Set(iter.Key(), iter.Value())
				}
				iter.Close()
				log.Printf("convert %s end.\n", name)
				log.Printf("compact %s start...\n", name)
				rdb.DB().CompactRange(gorocksdb.Range{})
				log.Printf("compact %s end.\n", name)
			}(name, os.Args[1], os.Args[2])
		}
	}

	wg.Wait()
}
