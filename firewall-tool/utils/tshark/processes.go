package tshark

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var sseChan = make(chan string)

type TsharkData struct {
	No          int    `json:"no"`
	Time        string `json:"time"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Protocol    string `json:"protocol"`
	Length      int    `json:"length"`
	Info        string `json:"info"`
}

func parseTsharkOutput(line string) (TsharkData, error) {
	re := regexp.MustCompile(`(\d+)\s+([\d.]+)\s+([\d.:a-fA-F]+|[\w:]+)\s+→\s+([\d.:a-fA-F]+|[\w:]+)\s+(\S+)\s+(\d+)\s+(.*)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 8 {
		return TsharkData{}, fmt.Errorf("failed to parse line: %s", line)
	}

	no, err := strconv.Atoi(matches[1])
	if err != nil {
		return TsharkData{}, fmt.Errorf("failed to parse No: %v", err)
	}
	time := matches[2]
	source := matches[3]
	destination := matches[4]
	protocol := matches[5]
	length, err := strconv.Atoi(matches[6])
	if err != nil {
		return TsharkData{}, fmt.Errorf("failed to parse Length: %v", err)
	}
	info := matches[7]

	return TsharkData{
		No:          no,
		Time:        time,
		Source:      source,
		Destination: destination,
		Protocol:    protocol,
		Length:      length,
		Info:        info,
	}, nil
}

func IsTsharkInstalled() bool {
	cmd := exec.Command("which", "tshark")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func InstallTshark() error {
	pkgManager := getPkgManager()
	cmd := exec.Command("sh", "-c", pkgManager+" install -y tshark")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getPkgManager() string {
	if _, err := exec.LookPath("apt-get"); err == nil {
		return "sudo apt-get"
	}
	if _, err := exec.LookPath("yum"); err == nil {
		return "sudo yum"
	}
	return "sudo apt-get"
}

func IsTsharkRunning() bool {
	cmd := exec.Command("pgrep", "tshark")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) != ""
}

// Helper function to upload the zip file asynchronously and delete it after upload
func uploadPcapZipFile(zipFilePath string, url string) error {
	file, err := os.Open(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", zipFilePath, err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("zip", filepath.Base(zipFilePath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("failed to create new request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %v", err)
	}
	fmt.Print("Response: ", resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Delete the zip file after successful upload
		if err := os.Remove(zipFilePath); err != nil {
			return fmt.Errorf("failed to delete zip file %s: %v", zipFilePath, err)
		}
		//delete capture.old.pcap
		if err := os.Remove("./capture.old.pcap"); err != nil {
			return fmt.Errorf("failed to delete file %s: %v", "capture.old.pcap", err)
		}
		//delete folder capture.old
		if err := os.RemoveAll("./capture.old"); err != nil {
			return fmt.Errorf("failed to delete folder %s: %v", "capture.old", err)
		}
		return fmt.Errorf("failed to upload file, status: %v", resp.Status)
	}

	fmt.Printf("File %s successfully uploaded\n", zipFilePath)

	// Delete the zip file after successful upload
	if err := os.Remove(zipFilePath); err != nil {
		return fmt.Errorf("failed to delete zip file %s: %v", zipFilePath, err)
	}
	//delete capture.old.pcap
	if err := os.Remove("./capture.old.pcap"); err != nil {
		return fmt.Errorf("failed to delete file %s: %v", "capture.old.pcap", err)
	}
	//delete folder capture.old
	if err := os.RemoveAll("./capture.old"); err != nil {
		return fmt.Errorf("failed to delete folder %s: %v", "capture.old", err)
	}
	fmt.Printf("Files deleted after upload\n")
	return nil
}

// Helper function to rename the file
func renameFile(oldPath string) (string, error) {
	// Extract the directory, file name and extension
	dir := filepath.Dir(oldPath)
	base := filepath.Base(oldPath)
	ext := filepath.Ext(base) // This extracts the extension (.pcap)

	// Get the file name without extension
	name := strings.TrimSuffix(base, ext)

	// Create the new name by appending .old before the extension
	newName := fmt.Sprintf("%s.old%s", name, ext)

	// Create the new file path
	newPath := filepath.Join(dir, newName)

	// Rename the file
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return "", fmt.Errorf("failed to rename file: %v", err)
	}

	fmt.Printf("File renamed from %s to %s\n", oldPath, newPath)
	return newPath, nil
}

// Zips a folder using `sudo` command and returns the zip file path
func zipFolderWithSudo(folderPath string) (string, error) {
	zipFileName := folderPath + ".zip"
	cmd := exec.Command("sudo", "zip", "-r", zipFileName, folderPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to zip folder %s: %v", folderPath, err)
	}

	fmt.Printf("Folder %s successfully zipped into %s\n", folderPath, zipFileName)
	return zipFileName, nil
}

// Split the PCAP file by session and zip the resulting folder
func processAndZipPcapFile(pcapFilePath string) (string, error) {
	splitCmd := exec.Command("../model_pipeline_files/SplitCap.exe", "-r", pcapFilePath, "-s", "session")
	err := splitCmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to split PCAP file: %v", err)
	}

	// After splitting, a folder with the same name as the PCAP file without the extension will be created
	folderName := strings.TrimSuffix(pcapFilePath, filepath.Ext(pcapFilePath))

	// Zip the folder
	zipFilePath, err := zipFolderWithSudo(folderName)
	if err != nil {
		return "", fmt.Errorf("failed to zip folder %s: %v", folderName, err)
	}

	return zipFilePath, nil
}

func StartTshark(outputFile string) error {
	// If file exists, rename and send to API
	if _, err := os.Stat(outputFile); err == nil {
		oldFile, _ := renameFile(outputFile)

		fmt.Printf("Renamed existing file to %s\n", oldFile)

		// Process and zip the file after renaming
		go func() {
			apiURL := "https://6c97-35-245-14-76.ngrok-free.app/predict" // Replace with your actual backend API endpoint
			zipFilePath, err := processAndZipPcapFile(oldFile)
			if err != nil {
				fmt.Printf("Failed to process and zip PCAP file: %v\n", err)
				return
			}

			if err := uploadPcapZipFile(zipFilePath, apiURL); err != nil {
				fmt.Printf("Failed to upload zip file: %v\n", err)
			}
		}()
	}

	file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", outputFile, err)
	}
	file.Close()

	cmd := exec.Command("sudo", "tshark", "-P", "-w", outputFile, "-F", "pcap")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %v", err)
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start tshark: %v", err)
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			data, err := parseTsharkOutput(line)
			if err != nil {
				fmt.Printf("error parsing tshark output: %v\n", err)
				continue
			}
			jsonData, err := json.Marshal(data)
			if err != nil {
				fmt.Printf("error marshalling json: %v\n", err)
				continue
			}
			sseChan <- string(jsonData)

			// Insert data into the database
			dbFile := "firewall.db"
			if err := InsertTsharkData(dbFile, data); err != nil {
				fmt.Printf("Failed to insert data into the database: %v\n", err)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("error reading from stdout: %v\n", err)
		}
		close(sseChan)
	}()

	return nil
}

func SseHandler(c *gin.Context) {
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.String(http.StatusInternalServerError, "Streaming unsupported!")
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	for msg := range sseChan {
		fmt.Fprintf(c.Writer, "data: %s\n\n", msg)
		flusher.Flush()
	}
}

func InsertTsharkData(dbFile string, data TsharkData) error {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	createTableSQL := `CREATE TABLE IF NOT EXISTS requests (
        no INTEGER PRIMARY KEY,
        time TEXT,
        source TEXT,
        destination TEXT,
        protocol TEXT,
        length INTEGER,
        info TEXT
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	insertSQL := `INSERT INTO requests (no, time, source, destination, protocol, length, info) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(insertSQL, data.No, data.Time, data.Source, data.Destination, data.Protocol, data.Length, data.Info)
	if err != nil {
		return fmt.Errorf("failed to insert data: %v", err)
	}

	return nil
}
