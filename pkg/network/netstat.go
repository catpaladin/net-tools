package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

const (
	listeningState = "0A"
)

// Netstat retrieves and prints TCP connections
func Netstat() {
	switch runtime.GOOS {
	case "linux":
		printLinuxTCPConnections()
	case "darwin":
		printDarwinTCPConnections()
	default:
		color.Red("Unsupported platform: %s\n", runtime.GOOS)
	}
}

func printLinuxTCPConnections() {
	files := []string{"/proc/net/tcp", "/proc/net/tcp6"}
	color.Green("%-45s %-15s %s\n", "Local Address", "Port", "PID/Program")
	color.Green(strings.Repeat("-", 80))
	for _, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open file %s: %v\n", filePath, err)
			continue
		}
		defer file.Close()

		rs := RealFileSystem{}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "  sl") {
				continue
			}
			localAddr, localPort, _, _, status, inode := parseLinuxTCPLine(line)
			if status != listeningState {
				continue
			}
			pid, program := getProcessByInode(rs, inode)
			if len(program) > 12 {
				program = program[:12]
			}
			color.Green("%-45s %-15s %d/%s\n",
				localAddr, localPort, pid, program)
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filePath, err)
		}
	}
}

// FileSystem defines an interface for filesystem operations.
type FileSystem interface {
	Glob(pattern string) ([]string, error)
	Readlink(name string) (string, error)
	ReadFile(name string) ([]byte, error)
}

// RealFileSystem is a concrete implementation of FileSystem using the os and filepath packages.
type RealFileSystem struct{}

func (fs RealFileSystem) Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

func (fs RealFileSystem) Readlink(name string) (string, error) {
	return os.Readlink(name)
}

func (fs RealFileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

// getProcessByInode retrieves the PID and command line of the process associated with the given inode.
func getProcessByInode(fs FileSystem, inode string) (int, string) {
	procDirs, err := fs.Glob("/proc/[0-9]*/fd/[0-9]*")
	if err != nil {
		return 0, ""
	}

	for _, fdPath := range procDirs {
		link, err := fs.Readlink(fdPath)
		if err != nil || !strings.Contains(link, "socket:["+inode+"]") {
			continue
		}

		parts := strings.Split(fdPath, "/")
		if len(parts) < 3 {
			continue
		}

		pid, err := strconv.Atoi(parts[2])
		if err != nil {
			continue
		}

		cmdline, err := fs.ReadFile(filepath.Join("/proc", parts[2], "cmdline"))
		if err != nil {
			continue
		}

		return pid, strings.ReplaceAll(string(cmdline), "\x00", " ")
	}

	return 0, ""
}

func parseLinuxTCPLine(line string) (localAddr, localPort, remoteAddr, remotePort, status, inode string) {
	fields := strings.Fields(line)
	if len(fields) < 12 {
		return "", "", "", "", "", ""
	}

	local := fields[1]
	remote := fields[2]
	status = fields[3]
	inode = fields[9]

	localAddr, localPort = parseHexIPPort(local)
	remoteAddr, remotePort = parseHexIPPort(remote)

	return localAddr, localPort, remoteAddr, remotePort, status, inode
}

func parseHexIPPort(hex string) (ip, port string) {
	parts := strings.Split(hex, ":")
	if len(parts) != 2 {
		return "", ""
	}

	ip = parseHexIP(parts[0])
	port = parseHexPort(parts[1])
	return ip, port
}

func parseHexIP(hex string) string {
	if len(hex) == 32 { // IPv6
		return parseHexIPv6(hex)
	}
	if len(hex) == 8 { // IPv4
		return parseHexIPv4(hex)
	}
	return "0.0.0.0" // Invalid length, return default invalid IP
}

func parseHexIPv4(hex string) string {
	var ip string
	for i := len(hex); i > 0; i -= 2 {
		ip += fmt.Sprintf("%d.", parseHexToInt(hex[i-2:i]))
	}
	return strings.TrimSuffix(ip, ".")
}

func parseHexIPv6(hex string) string {
	var segments []string
	for i := 0; i < len(hex); i += 4 {
		segment := hex[i : i+4]
		segments = append(segments, segment)
	}
	ip := strings.Join(segments, ":")
	parsedIP := net.ParseIP(ip)
	if parsedIP != nil {
		return parsedIP.String()
	}
	return ip
}

func parseHexPort(hex string) string {
	port := parseHexToInt(hex)
	return fmt.Sprintf("%d", port)
}

func parseHexToInt(hex string) int {
	var result int
	fmt.Sscanf(hex, "%x", &result)
	return result
}

// uses netstat on the mac, because syscall is too complicated
func printDarwinTCPConnections() {
	cmd := exec.Command("netstat", "-an", "|", "grep", "LISTEN")
	stdout, err := cmd.Output()
	if err != nil {
		color.Red("Failed to run netstat command: %v\n", err)
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "tcp") || strings.HasPrefix(line, "tcp4") || strings.HasPrefix(line, "tcp6") {
			color.Green(line)
		}
	}

	if err := scanner.Err(); err != nil {
		color.Red("Error reading netstat output: %v\n", err)
	}
}
