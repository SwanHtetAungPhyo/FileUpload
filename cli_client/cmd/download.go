package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/99designs/keyring"
	"github.com/SwanHtetAungPhyo/cli_client/utils"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"os"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a file from the server",
	Long:  "Download a file from the server, unzip, decrypt, and save it locally",
	Run: func(cmd *cobra.Command, args []string) {
		Download()
	},
}

func Download() {
	client := resty.New()
	req := client.R()

	var meta = struct {
		RemoteUrl string `json:"remote_url"`
		Token     string `json:"token"`
	}{
		RemoteUrl: remoteUrl,
		Token:     token,
	}
	metaDataFile, err := os.ReadFile("./meta.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = json.Unmarshal(metaDataFile, &meta)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	latestVersion, _, err := getLatestVersion()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	resp, err := req.Get(meta.RemoteUrl + "/download/" + latestVersion)
	if err != nil {
		fmt.Printf("Error downloading file: %v\n", err)
		return
	}
	if resp.StatusCode() != 200 {
		fmt.Printf("Error downloading file: %v\n", resp.Status())
		return
	}

	zipData := resp.Body()
	ring, err := keyring.Open(keyring.Config{
		ServiceName: "key-gen",
	})
	if err != nil {
		fmt.Printf("Error opening keyring: %v\n", err)
		return
	}

	item, err := ring.Get("cipher-key")
	if err != nil {
		fmt.Printf("Error getting decryption key: %v\n", err)
		return
	}

	decryptedBytes, err := utils.DecryptBytes(zipData, item.Data)
	if err != nil {
		fmt.Printf("Error during decryption: %v\n", err)
		return
	}
	err = os.WriteFile(distFile, decryptedBytes, 0644)
	if err != nil {
		fmt.Printf("Error during unzip and decryption: %v\n", err)
		return
	}
	fmt.Println("File successfully downloaded, unzipped, and decrypted.")
}
