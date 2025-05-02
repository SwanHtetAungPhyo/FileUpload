package utils

//
//import (
//	"archive/zip"
//	"bytes"
//	"crypto/aes"
//	"crypto/cipher"
//	"fmt"
//	"io"
//	"os"
//	"path/filepath"
//)
//
//// UnzipAndDecrypt handles the unzipping, decrypting, and writing of files
//func UnzipAndDecrypt(zipData []byte, outputDir string) error {
//	// Unzip the data into memory
//	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
//	if err != nil {
//		return fmt.Errorf("failed to read zip data: %w", err)
//	}
//
//	// Create the output directory if it doesn't exist
//	err = os.MkdirAll(outputDir, os.ModePerm)
//	if err != nil {
//		return fmt.Errorf("failed to create output directory: %w", err)
//	}
//
//	// Iterate through each file in the zip
//	for _, zipFile := range zipReader.File {
//		err = processZipFile(zipFile, outputDir)
//		if err != nil {
//			return fmt.Errorf("error processing file %s: %w", zipFile.Name, err)
//		}
//	}
//
//	return nil
//}
//
//// processZipFile processes each file in the zip (decrypts and writes it)
//func processZipFile(zipFile *zip.File, outputDir string) error {
//	// Open the zip file
//	fileReader, err := zipFile.Open()
//	if err != nil {
//		return err
//	}
//	defer fileReader.Close()
//
//	// Decrypt the file
//	decryptedData, err := decryptFile(fileReader)
//	if err != nil {
//		return err
//	}
//
//	// Create the file on disk
//	outputPath := filepath.Join(outputDir, zipFile.Name)
//	outputFile, err := os.Create(outputPath)
//	if err != nil {
//		return err
//	}
//	defer func(outputFile *os.File) {
//		err := outputFile.Close()
//		if err != nil {
//			fmt.Println(err.Error())
//		}
//	}(outputFile)
//
//	// Write the decrypted data to the file
//	_, err = outputFile.Write(decryptedData)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// decryptFile decrypts the given reader data using AES
//func decryptFile(input io.Reader) ([]byte, error) {
//	// Read the encrypted data from the input
//	encryptedData, err := io.ReadAll(input)
//	if err != nil {
//		return nil, fmt.Errorf("failed to read encrypted data: %w", err)
//	}
//
//	// Decrypt the data using AES (you can reuse the DecryptFile logic here)
//	key := generateKey()
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
//	}
//	aesGcm, err := cipher.NewGCM(block)
//	if err != nil {
//		return nil, fmt.Errorf("failed to create GCM cipher: %w", err)
//	}
//
//	// Get the nonce and cipher text
//	nonceSize := aesGcm.NonceSize()
//	if len(encryptedData) < nonceSize {
//		return nil, fmt.Errorf("encrypted data is too short")
//	}
//	nonce, cipherText := encryptedData[:nonceSize], encryptedData[nonceSize:]
//
//	// Decrypt the data
//	plainText, err := aesGcm.Open(nil, nonce, cipherText, nil)
//	if err != nil {
//		return nil, fmt.Errorf("failed to decrypt data: %w", err)
//	}
//
//	return plainText, nil
//}
