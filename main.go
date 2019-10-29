package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatalln("Você precisa informar o diretório inicial e o caminho absoluto do polyfill")
	}

	var wg sync.WaitGroup

	filepath.Walk(args[0], func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		name := info.Name()
		if info.IsDir() && (name == ".git" || name == "node_modules") {
			return filepath.SkipDir
		}

		if idx := len(name) - 4; idx > 0 && name[idx:] == ".php" {
			wg.Add(1)
			go includePolyfill(path, args[1], &wg)
		}

		return nil
	})

	wg.Wait()

	log.Printf("Concluído em %v\n", time.Since(start))
}

func includePolyfill(path, polyfill string, wg *sync.WaitGroup) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return
	}
	content := fmt.Sprintf("<?php include_once('%v'); ?>\r\n", polyfill)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content = content + scanner.Text() + "\r\n"
	}
	file.Close()
	ioutil.WriteFile(path, []byte(content), 0644)
	wg.Done()
}
