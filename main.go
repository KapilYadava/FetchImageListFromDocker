package main

import (
	"fmt"
	"log"

	git "github.com/go-git/go-git/v5"
)

func main() {
	fmt.Println("hello from go")
	repoPath := "/mnt/c/Users/A200154986/OSC/DiffCommitExample"
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatalf("could not open repository: %v", err)
	}

	fmt.Println(repo.Head())
}
