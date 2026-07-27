package commands

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"codedock.run/codedock/internal/config"
)

func runBackup() {
	dataDir := filepath.Clean(config.Get().Server.DataDir)

	parentDir := filepath.Dir(dataDir)
	baseName := filepath.Base(dataDir)

	backupsDir := filepath.Join(dataDir, "backups")
	if err := os.MkdirAll(backupsDir, 0755); err != nil {
		exitError("Failed to create backups directory: %v", err)
	}

	timestamp := time.Now().Format("20060102-150405")
	backupFile := filepath.Join(backupsDir, fmt.Sprintf("codedock-backup-%s.tar.gz", timestamp))

	fmt.Printf("📦 Creating backup of %s...\n", dataDir)
	cmd := exec.Command("tar", "--exclude=backups", "-czf", backupFile, "-C", parentDir, baseName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		exitError("Backup failed: %v", err)
	}

	fmt.Printf("✅ Backup created successfully at: %s\n", backupFile)
}

func runRestore(args []string) {
	if len(args) < 1 {
		exitError("Usage: codedockd restore <backup-file>")
	}

	backupFile := args[0]
	if _, err := os.Stat(backupFile); err != nil {
		exitError("Backup file not found: %s", backupFile)
	}

	dataDir := filepath.Clean(config.Get().Server.DataDir)
	baseName := filepath.Base(dataDir)

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", config.Get().Server.Port), 500*time.Millisecond)
	if err == nil {
		conn.Close()
		exitError("Codedock daemon is currently running. Please stop the daemon before running restore.")
	}

	fmt.Println("🔍 Validating backup archive...")

	cmdList := exec.Command("tar", "-tzf", backupFile)
	out, err := cmdList.Output()
	if err != nil {
		exitError("Failed to list archive contents: %v", err)
	}

	lines := strings.SplitSeq(string(out), "\n")
	for line := range lines {
		if line == "" {
			continue
		}
		clean := filepath.Clean(line)
		if strings.Contains(line, "..") || strings.HasPrefix(line, "/") {
			exitError("Archive validation failed: path traversal detected (%s)", line)
		}
		if clean != baseName && !strings.HasPrefix(clean, baseName+"/") {
			exitError("Archive validation failed: entry outside data directory (%s)", line)
		}
	}

	cmdTypes := exec.Command("tar", "-tvf", backupFile)
	outTypes, err := cmdTypes.Output()
	if err != nil {
		exitError("Failed to check archive types: %v", err)
	}

	linesTypes := strings.SplitSeq(string(outTypes), "\n")
	for line := range linesTypes {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "l") || strings.HasPrefix(line, "h") {
			exitError("Archive validation failed: symlinks or hardlinks detected")
		}
	}

	stagingDir, err := os.MkdirTemp(filepath.Dir(dataDir), "restore-staging-*")
	if err != nil {
		exitError("Failed to create staging directory: %v", err)
	}
	defer os.RemoveAll(stagingDir)

	fmt.Println("⏳ Extracting backup to staging area...")
	cmdExtract := exec.Command("tar", "--no-same-owner", "-xzf", backupFile, "-C", stagingDir)
	cmdExtract.Stdout = os.Stdout
	cmdExtract.Stderr = os.Stderr
	if err := cmdExtract.Run(); err != nil {
		exitError("Restore extraction failed: %v", err)
	}

	stagedDataDir := filepath.Join(stagingDir, baseName)
	stagedInfo, err := os.Lstat(stagedDataDir)
	if err != nil || !stagedInfo.IsDir() || stagedInfo.Mode()&os.ModeSymlink != 0 {
		exitError("Restore failed: %s in archive must be a real directory", baseName)
	}

	absDataDir, _ := filepath.Abs(dataDir)
	absBackupFile, _ := filepath.Abs(backupFile)
	isInsideDataDir := strings.HasPrefix(absBackupFile, absDataDir+string(os.PathSeparator))

	backupOldDir := absDataDir + ".bak-" + time.Now().Format("20060102150405")
	if err := os.Rename(dataDir, backupOldDir); err != nil && !os.IsNotExist(err) {
		exitError("Failed to backup existing data directory: %v", err)
	}

	if err := os.Rename(stagedDataDir, dataDir); err != nil {
		_ = os.Rename(backupOldDir, dataDir)
		exitError("Failed to replace data directory: %v", err)
	}

	if isInsideDataDir {
		rel, err := filepath.Rel(absDataDir, absBackupFile)
		if err == nil {
			oldBackupPath := filepath.Join(backupOldDir, rel)
			newBackupPath := filepath.Join(dataDir, rel)
			_ = os.MkdirAll(filepath.Dir(newBackupPath), 0755)
			_ = os.Rename(oldBackupPath, newBackupPath)
		}
	}

	_ = os.RemoveAll(backupOldDir)

	fmt.Println("✅ Restore completed successfully!")
	fmt.Println("🔄 Please restart Codedock to apply changes.")
}

func runDiagnostics() {
	fmt.Println("📊 System Diagnostics")
	fmt.Println("-------------------")

	fmt.Printf("OS/Arch:      %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("CPUs:         %d\n", runtime.NumCPU())

	cmdDocker := exec.Command("docker", "--version")
	outDocker, err := cmdDocker.Output()
	if err == nil {
		fmt.Printf("Docker:       %s", string(outDocker))
	} else {
		fmt.Printf("Docker:       Not found or not accessible\n")
	}

	cmdFree := exec.Command("free", "-m")
	outFree, err := cmdFree.Output()
	if err == nil {
		lines := strings.Split(string(outFree), "\n")
		if len(lines) > 1 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 3 {
				fmt.Printf("Memory:       %s MB total, %s MB used\n", fields[1], fields[2])
			}
		}
	}

	fmt.Println("\nDisk Usage (/):")
	cmdDf := exec.Command("df", "-h", "/")
	cmdDf.Stdout = os.Stdout
	cmdDf.Run()
}
