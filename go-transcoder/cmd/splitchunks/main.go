package splitchunks

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

const chunkSize = 1 * 1024 * 1024 // 1 MB

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: \n\tgo run main.go <MP4-file-path>")
		return
	}

	mp4Path := os.Args[1]
	file, err := os.Open(mp4Path)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		return
	}
	defer file.Close()

	chunkDir := filepath.Dir(mp4Path)
	fmt.Printf("Chunk directory: %s\n", chunkDir)

	chunkCount := 0
	for {
		chunkFileName := filepath.Join(chunkDir, strconv.Itoa(chunkCount)+".chunk")
		chunkFile, err := os.Create(chunkFileName)
		if err != nil {
			fmt.Printf("error creating chunk file: %v\n", err)
			return
		}

		if _, err := io.CopyN(chunkFile, file, chunkSize); err != nil {
			if err == io.EOF {
				chunkFile.Close()
				fmt.Println("finished splitting the file into chunks.")
				break
			} else {
				fmt.Printf("error copying chunk: %v\n", err)
				chunkFile.Close()
				return
			}
		}
		chunkFile.Close()
		fmt.Printf("created chunk: %s\n", chunkFileName)
		chunkCount++
	}
}
