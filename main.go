package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "laravel",
	Short: "Laravel CLI tool for project management",
	Long: `Laravel CLI is a command-line tool for creating and managing Laravel projects.

Available Commands:
  new     Create a new Laravel project

Usage:
  laravel new <project-name>    Create a new Laravel project in the specified directory

Examples:
  laravel new my-project        Creates a new Laravel project called 'my-project'

For more information about a command, run:
  laravel <command> --help`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new Laravel project",
	Long:  "Create a new Laravel project by cloning the Laravel repository and running the setup process.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		createNewProject(projectName)
	},
}

func main() {
	rootCmd.AddCommand(newCmd)
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createNewProject(projectName string) {
	// Check if directory already exists
	if _, err := os.Stat(projectName); !os.IsNotExist(err) {
		fmt.Printf("Error: Directory '%s' already exists\n", projectName)
		os.Exit(1)
	}

	// Setup signal handling for Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nGoodbye!")
		os.Exit(0)
	}()

	fmt.Printf("Creating new Laravel project: %s\n", projectName)
	
	// Clone Laravel repository
	fmt.Println("Cloning Laravel repository...")
	cmd := exec.Command("git", "clone", "https://github.com/laravel/laravel.git", projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error cloning Laravel repository: %v\n", err)
		os.Exit(1)
	}

	// Change to project directory
	projectDir := filepath.Join(".", projectName)
	
	// Run setup program
	runSetup(projectDir)
}

func runSetup(projectDir string) {
	fmt.Println("\nRunning Laravel project setup...")
	
	// Copy .env.example to .env
	envExamplePath := filepath.Join(projectDir, ".env.example")
	envPath := filepath.Join(projectDir, ".env")
	
	if err := copyFile(envExamplePath, envPath); err != nil {
		fmt.Printf("Error copying .env.example to .env: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Copied .env.example to .env")
	
	// Ask for configuration values
	reader := bufio.NewReader(os.Stdin)
	
	// Ask for App Port
	appPort := askForPort(reader, "App Port", 80)
	updateEnvFile(envPath, "APP_PORT", strconv.Itoa(appPort))
	
	// Ask for Vite Port
	vitePort := askForPort(reader, "Vite Port", 5173)
	updateEnvFile(envPath, "VITE_PORT", strconv.Itoa(vitePort))
	
	// Ask for App Name
	appName := askForString(reader, "App Name", "My Laravel Application")
	updateEnvFile(envPath, "APP_NAME", fmt.Sprintf("\"%s\"", appName))
	
	fmt.Printf("\nLaravel project '%s' has been created successfully!\n", filepath.Base(projectDir))
	fmt.Printf("Next steps:\n")
	fmt.Printf("  cd %s\n", filepath.Base(projectDir))
	fmt.Printf("  composer install\n")
	fmt.Printf("  php artisan key:generate\n")
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	
	return os.WriteFile(dst, input, 0644)
}

func askForPort(reader *bufio.Reader, name string, defaultValue int) int {
	for {
		fmt.Printf("%s (default: %d): ", name, defaultValue)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}
		
		input = strings.TrimSpace(input)
		if input == "" {
			// Check if default port is available
			if isPortAvailable(defaultValue) {
				return defaultValue
			} else {
				fmt.Printf("Default port %d is not available. Please choose another port.\n", defaultValue)
				continue
			}
		}
		
		port, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Please enter a valid number.")
			continue
		}
		
		if port < 1 || port > 65535 {
			fmt.Println("Port must be between 1 and 65535.")
			continue
		}
		
		if !isPortAvailable(port) {
			fmt.Printf("Port %d is not available. Please choose another port.\n", port)
			continue
		}
		
		return port
	}
}

func askForString(reader *bufio.Reader, name, defaultValue string) string {
	fmt.Printf("%s (default: %s): ", name, defaultValue)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return defaultValue
	}
	
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	
	return input
}

func isPortAvailable(port int) bool {
	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}

func updateEnvFile(envPath, key, value string) error {
	content, err := os.ReadFile(envPath)
	if err != nil {
		return err
	}
	
	lines := strings.Split(string(content), "\n")
	updated := false
	
	// Look for existing key and update it
	keyPattern := regexp.MustCompile("^" + regexp.QuoteMeta(key) + "=.*")
	for i, line := range lines {
		if keyPattern.MatchString(line) {
			lines[i] = fmt.Sprintf("%s=%s", key, value)
			updated = true
			break
		}
	}
	
	// If key doesn't exist, add it
	if !updated {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}
	
	newContent := strings.Join(lines, "\n")
	return os.WriteFile(envPath, []byte(newContent), 0644)
}
