package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/99designs/keyring"
	"github.com/SwanHtetAungPhyo/cli_client/utils"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file to a server",
	Long:  "Upload a file to a server with encryption",
	Run: func(cmd *cobra.Command, args []string) {
		if err := Upload(); err != nil {
			log.Fatalf("Upload failed: %v", err)
		}
	},
}

func Upload() error {
	// Validate input file
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("error getting absolute path: %w", err)
	}

	fmt.Printf("Uploading file...\n%s\n", absPath)

	// Calculate hash of original file
	originalHash, err := calculateFileHash(absPath)
	if err != nil {
		return fmt.Errorf("error calculating file hash: %w", err)
	}
	originalHashStr := hex.EncodeToString(originalHash)

	// Check if file has changed
	if shouldSkipUpload(originalHashStr) {
		fmt.Println("File content hasn't changed - skipping upload")
		return nil
	}

	zipBytes, _, err := encryptFile(absPath)
	if err != nil {
		return err
	}

	// Read server metadata
	meta, err := readMetadata()
	if err != nil {
		return err
	}

	// Upload the file
	response, err := uploadToServer(meta.RemoteUrl, filepath.Base(absPath), zipBytes)
	if err != nil {
		return err
	}

	// Save version history
	if err := saveVersionHistory(response.Filename, originalHashStr); err != nil {
		return fmt.Errorf("error saving version history: %w", err)
	}

	fmt.Println("Successfully uploaded file to server")
	return nil
}

func calculateFileHash(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

func shouldSkipUpload(currentHash string) bool {
	_, latestHash, err := getLatestVersion()
	if err != nil && !strings.Contains(err.Error(), "no versions found") {
		log.Printf("Warning: error checking versions: %v", err)
		return false
	}
	return err == nil && strings.EqualFold(latestHash, currentHash)
}

func encryptFile(path string) ([]byte, []byte, error) {
	ring, err := keyring.Open(keyring.Config{ServiceName: "key-gen"})
	if err != nil {
		return nil, nil, fmt.Errorf("error opening keyring: %w", err)
	}

	item, err := ring.Get("cipher-key")
	if err != nil {
		return nil, nil, fmt.Errorf("error getting encryption key: %w", err)
	}

	return utils.EncryptFile(path, item.Data)
}

func readMetadata() (struct {
	RemoteUrl string `json:"remote_url"`
	Token     string `json:"token"`
}, error) {
	var meta struct {
		RemoteUrl string `json:"remote_url"`
		Token     string `json:"token"`
	}

	metaFile, err := os.ReadFile("meta.json")
	if err != nil {
		return meta, fmt.Errorf("error reading meta.json: %w", err)
	}

	if err := json.Unmarshal(metaFile, &meta); err != nil {
		return meta, fmt.Errorf("error parsing meta.json: %w", err)
	}

	if meta.RemoteUrl == "" {
		return meta, fmt.Errorf("no remote URL specified in meta.json")
	}

	return meta, nil
}

func uploadToServer(remoteUrl, baseName string, data []byte) (struct {
	Success  bool   `json:"success"`
	Filename string `json:"filename"`
}, error) {
	var response struct {
		Success  bool   `json:"success"`
		Filename string `json:"filename"`
	}

	client := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second)

	fileName := baseName + time.Now().Format("20060102150405") + ".enc"

	resp, err := client.R().
		SetFileReader("file", fileName, bytes.NewReader(data)).
		Post(remoteUrl + "/upload")

	if err != nil {
		return response, fmt.Errorf("error uploading file: %w", err)
	}

	if resp.StatusCode() != 200 {
		return response, fmt.Errorf("server returned error: %s\nResponse body: %s", resp.Status(), resp.Body())
	}

	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return response, fmt.Errorf("error parsing response: %w", err)
	}

	if !response.Success {
		return response, fmt.Errorf("upload failed according to server response")
	}

	return response, nil
}

func saveVersionHistory(filename, hash string) error {
	// Ensure versions directory exists
	if err := os.MkdirAll(".filehash", 0755); err != nil {
		return fmt.Errorf("error creating versions directory: %w", err)
	}

	file, err := os.OpenFile(filepath.Join(".filehash", "versions.csv"),
		os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return fmt.Errorf("error opening versions file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.Write([]string{filename, hash})
}

func getLatestVersion() (string, string, error) {
	file, err := os.Open(filepath.Join(".filehash", "versions.csv"))
	if err != nil {
		return "", "", fmt.Errorf("error opening versions file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return "", "", fmt.Errorf("error reading versions file: %w", err)
	}

	if len(records) == 0 {
		return "", "", fmt.Errorf("no versions found")
	}

	return records[len(records)-1][0], records[len(records)-1][1], nil
}
