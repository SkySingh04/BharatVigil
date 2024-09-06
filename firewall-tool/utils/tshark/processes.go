package tshark

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
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

func StartTshark(outputFile string) error {
	if _, err := os.Stat(outputFile); err == nil {
		if err := os.Remove(outputFile); err != nil {
			return fmt.Errorf("failed to remove existing file %s: %v", outputFile, err)
		}
		fmt.Printf("removed existing file %s\n", outputFile)
	}

	file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", outputFile, err)
	}
	file.Close()

	cmd := exec.Command("sudo", "tshark", "-P", "-w", outputFile)
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
