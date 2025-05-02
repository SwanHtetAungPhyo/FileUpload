package cmd

import (
	"bufio"
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path"
	"sort"
)

var (
	filePath  string // Path of the directory to read
	remoteUrl string
	token     string
)

var initCommand = &cobra.Command{
	Use:   "init",
	Short: "Initialize tracking files in the directory",
	Long:  `This command initializes the tracking of files in the specified directory by calculating their hashes and storing metadata.`,
	Run: func(cmd *cobra.Command, args []string) {
		Start()
	},
}

var statusCommand = &cobra.Command{
	Use:   "status",
	Short: "Show the status of files (modified, new, etc.)",
	Long:  `This command shows the status of files, including those that have been modified or added since the last commit.`,
	Run: func(cmd *cobra.Command, args []string) {
		ShowStatus()
	},
}

var commitCommand = &cobra.Command{
	Use:   "commit",
	Short: "Commit changes to the directory",
	Long:  `This command commits the current state of the directory to the metadata file.`,
	Run: func(cmd *cobra.Command, args []string) {
		CommitChanges()
	},
}

var remoteUrlCommand = &cobra.Command{
	Use:   "remote-url",
	Short: "Show the remote URL of a remote repository",
	Long:  `This command shows the remote URL of a remote repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		RemoteUploadUrlAdd()
	},
}

type FileInfo struct {
	FileName string
	FileSize float64
	FileType string
	FileHash string
}

func Start() {
	if filePath == "." {
		var err error
		filePath, err = os.Getwd()
		if err != nil {
			panic(fmt.Sprintf("Error getting current working directory: %v", err))
		}
	}

	// Print the file path being accessed for debugging
	fmt.Printf("Reading files from: %s\n", filePath)

	// Read the directory contents
	files, err := os.ReadDir(filePath)
	if err != nil {
		panic(fmt.Sprintf("Error reading directory: %v", err))
	}

	var fileContentHashes []*FileInfo

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fullFilePath := path.Join(filePath, file.Name())
		fileType, _ := os.Lstat(fullFilePath)
		fileMeta, err := os.Stat(fullFilePath)
		if os.IsNotExist(err) {
			fmt.Printf("File does not exist: %s\n", fullFilePath)
			continue
		}

		// Open the file
		f, err := os.Open(fullFilePath)
		if err != nil {
			panic(fmt.Sprintf("Error opening file %s: %v", fullFilePath, err))
		}

		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				panic(fmt.Sprintf("Error closing file %s: %v", fullFilePath, err))
			}
		}(f)

		fileHash := sha1.New()

		if err != nil {
			panic(fmt.Sprintf("Error reading file %s: %v", fullFilePath, err))
		}
		reader := bufio.NewReader(f)
		_, err = io.Copy(fileHash, reader)
		if err != nil {
			panic(fmt.Sprintf("Error reading file %s: %v", fullFilePath, err))
		}

		fileHashString := hex.EncodeToString(fileHash.Sum(nil))

		fileInfo := &FileInfo{
			FileName: file.Name(),
			FileHash: fileHashString,
			FileSize: float64(fileMeta.Size() / 1024),
			FileType: fileType.Name(),
		}
		fileContentHashes = append(fileContentHashes, fileInfo)
	}

	fmt.Printf("Total files processed: %d\n", len(fileContentHashes))
	if len(fileContentHashes) == 0 {
		fmt.Println("No files found")
		return
	}

	sort.Slice(fileContentHashes, func(i, j int) bool {
		return fileContentHashes[i].FileHash < fileContentHashes[j].FileHash
	})

	err = SaveToMetaFile(".meta", fileContentHashes)
	if err != nil {
		fmt.Printf("Error saving metadata file: %v", err)
		return
	}
	fmt.Printf("Total files processed and committed: %d\n", len(fileContentHashes))
}

func ShowStatus() {
	var fileContentHashes []*FileInfo
	err := LoadMetaFile(".meta", &fileContentHashes)
	if err != nil {
		fmt.Printf("Error loading metadata file: %v\n", err)
		return
	}

	currentFiles := make(map[string]string)
	for _, fileInfo := range fileContentHashes {
		currentFiles[fileInfo.FileName] = fileInfo.FileHash
	}

	files, err := os.ReadDir(filePath)
	if err != nil {
		fmt.Printf("Error reading directory: %v", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fullFilePath := path.Join(filePath, file.Name())
		_, err := os.Stat(fullFilePath)
		if err != nil {
			fmt.Printf("Error getting file stats: %v\n", err)
			continue
		}

		f, err := os.Open(fullFilePath)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			continue
		}
		defer f.Close()

		fileHash := sha1.New()
		_, err = io.Copy(fileHash, f)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			continue
		}
		currentHash := hex.EncodeToString(fileHash.Sum(nil))

		// Check if the file has been modified or added
		storedHash, exists := currentFiles[file.Name()]
		if !exists {
			fmt.Printf("New file: %s\n", file.Name())
		} else if storedHash != currentHash {
			fmt.Printf("Modified file: %s\n", file.Name())
		}
	}
}

func CommitChanges() {
	var fileContentHashes []*FileInfo
	files, err := os.ReadDir(filePath)
	if err != nil {
		fmt.Printf("Error reading directory: %v", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fullFilePath := path.Join(filePath, file.Name())
		fileMeta, err := os.Stat(fullFilePath)
		if err != nil {
			fmt.Printf("Error getting file stats: %v\n", err)
			continue
		}

		f, err := os.Open(fullFilePath)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			continue
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				fmt.Printf("Error closing file %s: %v", fullFilePath, err)
			}
		}(f)

		fileHash := sha1.New()
		_, err = io.Copy(fileHash, f)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			continue
		}

		fileInfo := &FileInfo{
			FileName: file.Name(),
			FileHash: hex.EncodeToString(fileHash.Sum(nil)),
			FileSize: float64(fileMeta.Size() / 1024),
		}
		fileContentHashes = append(fileContentHashes, fileInfo)
	}

	err = SaveToMetaFile(".meta", fileContentHashes)
	if err != nil {
		fmt.Printf("Error saving metadata file: %v\n", err)
		return
	}
	fmt.Println("Changes committed.")
}

func SaveToMetaFile(filepath string, contents []*FileInfo) error {
	filepath = path.Join(filePath, filepath)
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creating metadata file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file %s: %v", filepath, err)
		}
	}(file)

	err = gob.NewEncoder(file).Encode(contents)
	if err != nil {
		return fmt.Errorf("error encoding metadata: %v", err)
	}
	return nil
}

func LoadMetaFile(filepath string, contents *[]*FileInfo) error {
	filepath = path.Join(filePath, filepath)
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("error opening metadata file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file %s: %v", filepath, err)
		}
	}(file)

	err = gob.NewDecoder(file).Decode(contents)
	if err != nil {
		return fmt.Errorf("error decoding metadata: %v", err)
	}
	return nil
}

func RemoteUploadUrlAdd() {
	if remoteUrl == "" {
		fmt.Println("No remote URL specified")
		return
	}

	pwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Error getting current working directory: %v", err))
	}
	metaFileStore, err := os.OpenFile(path.Join(pwd, "meta.json"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		fmt.Printf("Error opening metadata file: %v\n", err)
		return
	}
	defer func(metaFileStore *os.File) {
		err := metaFileStore.Close()
		if err != nil {
			fmt.Printf("Error closing file %s: %v", metaFileStore.Name(), err)
		}
	}(metaFileStore)

	var meta = struct {
		RemoteUrl string `json:"remote_url"`
		Token     string `json:"token"`
	}{
		RemoteUrl: remoteUrl,
		Token:     token,
	}
	jsonData, err := json.MarshalIndent(&meta, "", "	")
	if err != nil {
		fmt.Printf("Error marshalling metadata: %v\n", err)
		return
	}
	_, err = metaFileStore.Write(jsonData)
	if err != nil {
		fmt.Printf("Error writing metadata file: %v\n", err)
	}
}
