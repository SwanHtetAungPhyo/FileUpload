
---

# 🔐 File Upload CLI with Version Control

This CLI tool enables secure, encrypted file uploads to a server with version control and deduplication, using a locally stored encryption key and hash-based versioning.

## ✨ Features

* 🔑 **Key Management**: CLI securely stores the encryption key using the OS keychain.
* 📦 **Init + Hashing**: Initialize a project and calculate SHA256 hashes of files/folders.
* 🔒 **Encrypted Upload**: Upload files to the server with AES encryption.
* 📁 **Destination Handling**: Server saves encrypted files in a designated folder.
* ⚠️ **No Duplicates**: Only changed files (based on hash comparison) can be uploaded again.
* ⬇️ **Secure Download**: Encrypted files can be downloaded and decrypted using the stored key.
* 📂 **Restore to Path**: Decrypted files are restored to the original specified location.

---

## 🛠️ Commands Overview

| Command      | Flags                           | Description                                    |
| ------------ | ------------------------------- | ---------------------------------------------- |
| `keygen`     | `-u <username>` `-p <password>` | Generates and stores an encryption key.        |
| `init`       | `-f <file path>`                | Initializes a new project and calculates hash. |
| `status`     | `-f <file path>`                | Checks file change status by comparing hash.   |
| `commit`     | `-f <file path>`                | Commits file for upload if it has changed.     |
| `upload`     | `-f <file path>`                | Uploads encrypted file to server.              |
| `download`   | `-d <destination path>`         | Downloads and decrypts file from server.       |
| `remote-url` | `-u <url>` `-t <token>`         | Sets or updates remote server URL and token.   |

---

## 🧠 How It Works

1. **Initialization**

    * Calculates SHA256 hash of the target file/folder.
    * Saves hash to `.versions.csv`.

2. **Encryption & Upload**

    * File is encrypted with key from keychain.
    * If hash already exists in `.versions.csv`, file won't be uploaded.
    * Upload includes encrypted file + filename + metadata.

3. **Server Handling**

    * Receives encrypted file.
    * Stores in destination folder with unique filename.

4. **Download & Decryption**

    * Downloads encrypted file using metadata.
    * Uses stored encryption key to decrypt it.
    * Writes decrypted content to specified destination.

---

## 🧪 Example Usage

```bash
cli-client keygen -u alice -p secret123
cli-client init -f ./data/report.pdf
cli-client upload -f ./data/report.pdf
cli-client download -d ./restored/report.pdf
```

---

