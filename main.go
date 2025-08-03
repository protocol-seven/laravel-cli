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

const (
	VERSION = "1.0.0"
)

var (
	// Command flags
	dev                     bool
	git                     bool
	branch                  string
	github                  string
	organization            string
	database                string
	react                   bool
	vue                     bool
	livewire                bool
	livewireClassComponents bool
	workos                  bool
	pest                    bool
	phpunit                 bool
	npm                     bool
	using                   string
	force                   bool
	quiet                   bool
)

var databaseDrivers = []string{"mysql", "mariadb", "pgsql", "sqlite", "sqlsrv"}

var rootCmd = &cobra.Command{
	Use:     "laravel",
	Version: VERSION,
	Short:   "Laravel CLI tool for project management",
	Long: `  _                               _
  | |                             | |
  | |     __ _ _ __ __ ___   _____| |
  | |    / _` + "`" + ` |  __/ _` + "`" + ` \ \ / / _ \ |
  | |___| (_| | | | (_| |\ V /  __/ |
  |______\__,_|_|  \__,_| \_/ \___|_|

Laravel CLI is a command-line tool for creating and managing Laravel projects.
This tool replicates the functionality of the official Laravel installer.

Available Commands:
  new     Create a new Laravel application

Usage:
  laravel new <project-name>    Create a new Laravel project in the specified directory

Examples:
  laravel new my-project                        Create a Laravel project
  laravel new my-project --git                  Initialize with Git
  laravel new my-project --database=mysql       Use MySQL database
  laravel new my-project --pest                 Use Pest testing framework
  laravel new my-project --vue                  Install Vue starter kit

For more information about a command, run:
  laravel <command> --help`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new Laravel application",
	Long:  "Create a new Laravel application with interactive setup and optional features.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		createNewProject(projectName)
	},
}

func main() {
	// Add flags to the new command
	newCmd.Flags().BoolVar(&dev, "dev", false, "Install the latest \"development\" release")
	newCmd.Flags().BoolVar(&git, "git", false, "Initialize a Git repository")
	newCmd.Flags().StringVar(&branch, "branch", "", "The branch that should be created for a new repository")
	newCmd.Flags().StringVar(&github, "github", "", "Create a new repository on GitHub")
	newCmd.Flags().StringVar(&organization, "organization", "", "The GitHub organization to create the new repository for")
	newCmd.Flags().StringVar(&database, "database", "", fmt.Sprintf("The database driver your application will use. Possible values are: %s", strings.Join(databaseDrivers, ", ")))
	newCmd.Flags().BoolVar(&react, "react", false, "Install the React Starter Kit")
	newCmd.Flags().BoolVar(&vue, "vue", false, "Install the Vue Starter Kit")
	newCmd.Flags().BoolVar(&livewire, "livewire", false, "Install the Livewire Starter Kit")
	newCmd.Flags().BoolVar(&livewireClassComponents, "livewire-class-components", false, "Generate stand-alone Livewire class components")
	newCmd.Flags().BoolVar(&workos, "workos", false, "Use WorkOS for authentication")
	newCmd.Flags().BoolVar(&pest, "pest", false, "Install the Pest testing framework")
	newCmd.Flags().BoolVar(&phpunit, "phpunit", false, "Install the PHPUnit testing framework")
	newCmd.Flags().BoolVar(&npm, "npm", false, "Install and build NPM dependencies")
	newCmd.Flags().StringVar(&using, "using", "", "Install a custom starter kit from a community maintained package")
	newCmd.Flags().BoolVarP(&force, "force", "f", false, "Forces install even if the directory already exists")
	newCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress output")

	rootCmd.AddCommand(newCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createNewProject(projectName string) {
	// Validate project name
	if err := validateProjectName(projectName); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Check if directory already exists
	if !force {
		if _, err := os.Stat(projectName); !os.IsNotExist(err) {
			fmt.Printf("Error: Directory '%s' already exists. Use --force to override.\n", projectName)
			os.Exit(1)
		}
	}

	// Setup signal handling for Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nGoodbye!")
		os.Exit(0)
	}()

	// Validate database option if provided
	if database != "" && !contains(databaseDrivers, database) {
		fmt.Printf("Error: Invalid database driver [%s]. Possible values are: %s\n",
			database, strings.Join(databaseDrivers, ", "))
		os.Exit(1)
	}

	if !quiet {
		printLaravelLogo()
		fmt.Printf("Creating new Laravel project: %s\n", projectName)
	}

	// Ensure required tools are available
	ensureRequiredTools()

	// Create project directory if force is used
	if force {
		if err := os.RemoveAll(projectName); err != nil && !os.IsNotExist(err) {
			fmt.Printf("Error removing existing directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Determine which installation method to use
	starterKit := getStarterKit()
	version := getVersion()

	// Create Laravel project using composer
	createLaravelProject(projectName, starterKit, version)

	// Change to project directory
	projectDir := filepath.Join(".", projectName)

	// Run post-installation setup
	runPostInstallation(projectDir)

	// Interactive setup
	runInteractiveSetup(projectDir)

	// Git setup if requested
	if git || github != "" {
		initializeGitRepository(projectDir)
	}

	// Install testing framework
	if pest {
		installPest(projectDir)
	}

	// GitHub setup if requested
	if github != "" {
		createGitHubRepository(projectName, projectDir)
	}

	// NPM setup if requested
	if npm {
		runNpmCommands(projectDir)
	}

	// Final instructions
	printCompletionMessage(projectName)
}

func validateProjectName(name string) error {
	// Project name validation (similar to PHP installer)
	if matched, _ := regexp.MatchString(`[^\pL\pN\-_.]`, name); matched {
		return fmt.Errorf("the name may only contain letters, numbers, dashes, underscores, and periods")
	}
	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func printLaravelLogo() {
	fmt.Print(`
  ` + "\033[31m" + ` _                               _
  | |                             | |
  | |     __ _ _ __ __ ___   _____| |
  | |    / _` + "`" + ` |  __/ _` + "`" + ` \ \ / / _ \ |
  | |___| (_| | | | (_| |\ V /  __/ |
  |______\__,_|_|  \__,_| \_/ \___|_|` + "\033[0m" + `

`)
}

func ensureRequiredTools() {
	// Check if composer is available
	if _, err := exec.LookPath("composer"); err != nil {
		fmt.Println("Error: Composer is required but not found in PATH")
		fmt.Println("Please install Composer: https://getcomposer.org/")
		os.Exit(1)
	}

	// Check if PHP is available
	if _, err := exec.LookPath("php"); err != nil {
		fmt.Println("Error: PHP is required but not found in PATH")
		os.Exit(1)
	}
}

func getStarterKit() string {
	if react {
		return "laravel/react-starter-kit"
	}
	if vue {
		return "laravel/vue-starter-kit"
	}
	if livewire {
		return "laravel/livewire-starter-kit"
	}
	if using != "" {
		return using
	}
	return ""
}

func getVersion() string {
	if dev {
		return "dev-master"
	}
	return ""
}

func createLaravelProject(projectName, starterKit, version string) {
	var cmd *exec.Cmd
	var args []string

	if starterKit != "" {
		// Use starter kit
		args = []string{"create-project", starterKit, projectName, "--stability=dev"}

		// Handle Laravel starter kit variations
		if isLaravelStarterKit(starterKit) {
			if livewireClassComponents {
				starterKit = starterKit + ":dev-components"
			}
			if workos {
				starterKit = starterKit + ":dev-workos"
			}
			args[1] = starterKit
		}
	} else {
		// Standard Laravel installation
		args = []string{"create-project", "laravel/laravel", projectName}
		if version != "" {
			args = append(args, version)
		}
		args = append(args, "--remove-vcs", "--prefer-dist", "--no-scripts")
	}

	if !quiet {
		fmt.Println("Installing Laravel...")
	}

	cmd = exec.Command("composer", args...)
	if !quiet {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error creating Laravel project: %v\n", err)
		os.Exit(1)
	}
}

func isLaravelStarterKit(starterKit string) bool {
	return strings.HasPrefix(starterKit, "laravel/")
}

func runPostInstallation(projectDir string) {
	commands := [][]string{
		{"composer", "run", "post-root-package-install", "-d", projectDir},
		{"php", filepath.Join(projectDir, "artisan"), "key:generate", "--ansi"},
	}

	// Make artisan executable on Unix systems
	if os.PathSeparator == '/' {
		artisanPath := filepath.Join(projectDir, "artisan")
		if err := os.Chmod(artisanPath, 0755); err != nil {
			fmt.Printf("Warning: Could not make artisan executable: %v\n", err)
		}
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		if !quiet {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

		if err := cmd.Run(); err != nil {
			fmt.Printf("Error running %s: %v\n", strings.Join(cmdArgs, " "), err)
		}
	}
}

func runInteractiveSetup(projectDir string) {
	if !quiet {
		fmt.Println("\nRunning Laravel project setup...")
	}

	// Copy .env.example to .env if it doesn't exist
	envExamplePath := filepath.Join(projectDir, ".env.example")
	envPath := filepath.Join(projectDir, ".env")

	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		if err := copyFile(envExamplePath, envPath); err != nil {
			fmt.Printf("Error copying .env.example to .env: %v\n", err)
			os.Exit(1)
		}
		if !quiet {
			fmt.Println("Copied .env.example to .env")
		}
	}

	// Interactive prompts for configuration
	reader := bufio.NewReader(os.Stdin)

	// Configure database if not specified via flag
	if database == "" {
		database = promptForDatabase(reader)
	}

	// Configure database connection
	configureDatabaseConnection(projectDir, database, filepath.Base(projectDir))

	// Ask for App URL configuration
	appURL := askForString(reader, "App URL", "http://localhost:8000")
	updateEnvFile(envPath, "APP_URL", appURL)

	// Database migration prompt
	if database != "" && database != "sqlite" {
		if askForConfirmation(reader, "Would you like to run the default database migrations?") {
			runMigrations(projectDir)
		}
	} else if database == "sqlite" {
		// Create SQLite database file
		dbPath := filepath.Join(projectDir, "database", "database.sqlite")
		if _, err := os.Create(dbPath); err != nil {
			fmt.Printf("Warning: Could not create SQLite database file: %v\n", err)
		} else {
			if askForConfirmation(reader, "Would you like to run the default database migrations?") {
				runMigrations(projectDir)
			}
		}
	}

	// NPM prompt if not specified via flag
	if !npm {
		npm = askForConfirmation(reader, "Would you like to run npm install and npm run build?")
	}
}

func promptForDatabase(reader *bufio.Reader) string {
	availableDatabases := getAvailableDatabases()

	fmt.Println("\nWhich database will your application use?")
	for i, db := range availableDatabases {
		fmt.Printf("%d) %s\n", i+1, db)
	}

	for {
		fmt.Print("Please select (1-" + strconv.Itoa(len(availableDatabases)) + ") [1]: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			return "sqlite" // Default to SQLite
		}

		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(availableDatabases) {
			fmt.Println("Please enter a valid number.")
			continue
		}

		selected := availableDatabases[choice-1]
		// Extract database driver from display string
		parts := strings.Fields(selected)
		if len(parts) > 0 {
			switch parts[0] {
			case "SQLite":
				return "sqlite"
			case "MySQL":
				return "mysql"
			case "MariaDB":
				return "mariadb"
			case "PostgreSQL":
				return "pgsql"
			case "SQL":
				return "sqlsrv"
			}
		}
		return "sqlite"
	}
}

func getAvailableDatabases() []string {
	return []string{
		"SQLite",
		"MySQL",
		"MariaDB",
		"PostgreSQL",
		"SQL Server",
	}
}

func configureDatabaseConnection(projectDir, dbDriver, projectName string) {
	envPath := filepath.Join(projectDir, ".env")
	envExamplePath := filepath.Join(projectDir, ".env.example")

	// Update DB_CONNECTION
	updateEnvFile(envPath, "DB_CONNECTION", dbDriver)
	updateEnvFile(envExamplePath, "DB_CONNECTION", dbDriver)

	if dbDriver == "sqlite" {
		// Comment out database configuration for SQLite
		commentDatabaseConfigForSQLite(envPath)
		commentDatabaseConfigForSQLite(envExamplePath)
		return
	}

	// Uncomment database configuration for non-SQLite databases
	uncommentDatabaseConfig(envPath)
	uncommentDatabaseConfig(envExamplePath)

	// Set default ports for specific databases
	defaultPorts := map[string]string{
		"pgsql":  "5432",
		"sqlsrv": "1433",
	}

	if port, exists := defaultPorts[dbDriver]; exists {
		updateEnvFile(envPath, "DB_PORT", port)
		updateEnvFile(envExamplePath, "DB_PORT", port)
	}

	// Set database name based on project name
	dbName := strings.ReplaceAll(strings.ToLower(projectName), "-", "_")
	updateEnvFile(envPath, "DB_DATABASE", dbName)
	updateEnvFile(envExamplePath, "DB_DATABASE", dbName)
}

func commentDatabaseConfigForSQLite(envPath string) {
	configs := []string{
		"DB_HOST=127.0.0.1",
		"DB_PORT=3306",
		"DB_DATABASE=laravel",
		"DB_USERNAME=root",
		"DB_PASSWORD=",
	}

	for _, config := range configs {
		replaceInFile(envPath, config, "# "+config)
	}
}

func uncommentDatabaseConfig(envPath string) {
	configs := []string{
		"# DB_HOST=127.0.0.1",
		"# DB_PORT=3306",
		"# DB_DATABASE=laravel",
		"# DB_USERNAME=root",
		"# DB_PASSWORD=",
	}

	for _, config := range configs {
		replaceInFile(envPath, config, strings.TrimPrefix(config, "# "))
	}
}

func runMigrations(projectDir string) {
	fmt.Println("Running database migrations...")
	cmd := exec.Command("php", "artisan", "migrate", "--no-interaction")
	cmd.Dir = projectDir
	if !quiet {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: Database migration failed: %v\n", err)
	}
}

func initializeGitRepository(projectDir string) {
	if !quiet {
		fmt.Println("Initializing Git repository...")
	}

	branchName := branch
	if branchName == "" {
		branchName = getDefaultGitBranch()
	}

	commands := [][]string{
		{"git", "init", "-q"},
		{"git", "add", "."},
		{"git", "commit", "-q", "-m", "Set up a fresh Laravel app"},
		{"git", "branch", "-M", branchName},
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = projectDir
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: Git command failed: %v\n", err)
		}
	}
}

func getDefaultGitBranch() string {
	cmd := exec.Command("git", "config", "--global", "init.defaultBranch")
	output, err := cmd.Output()
	if err != nil || len(output) == 0 {
		return "main"
	}
	return strings.TrimSpace(string(output))
}

func installPest(projectDir string) {
	if !quiet {
		fmt.Println("Installing Pest testing framework...")
	}

	commands := [][]string{
		{"composer", "remove", "phpunit/phpunit", "--dev", "--no-update"},
		{"composer", "require", "pestphp/pest", "pestphp/pest-plugin-laravel", "--no-update", "--dev"},
		{"composer", "update"},
		{"php", "./vendor/bin/pest", "--init"},
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = projectDir
		cmd.Env = append(os.Environ(), "PEST_NO_SUPPORT=true")
		if !quiet {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: Pest installation step failed: %v\n", err)
		}
	}
}

func createGitHubRepository(projectName, projectDir string) {
	// Check if GitHub CLI is available and authenticated
	cmd := exec.Command("gh", "auth", "status")
	if err := cmd.Run(); err != nil {
		fmt.Println("Warning: GitHub CLI not available or not authenticated. Skipping GitHub repository creation.")
		return
	}

	if !quiet {
		fmt.Println("Creating GitHub repository...")
	}

	repoName := projectName
	if organization != "" {
		repoName = organization + "/" + projectName
	}

	flags := "--private"
	if github != "" && github != "true" {
		flags = github
	}

	cmd = exec.Command("gh", "repo", "create", repoName, "--source=.", "--push", flags)
	cmd.Dir = projectDir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")

	if !quiet {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: GitHub repository creation failed: %v\n", err)
	}
}

func runNpmCommands(projectDir string) {
	if !quiet {
		fmt.Println("Installing and building NPM dependencies...")
	}

	commands := [][]string{
		{"npm", "install"},
		{"npm", "run", "build"},
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = projectDir
		if !quiet {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: NPM command failed: %v\n", err)
		}
	}
}

func printCompletionMessage(projectName string) {
	fmt.Printf("\n\033[44;37m INFO \033[0m Application ready in \033[1m[%s]\033[0m. You can start your local development using:\n\n", projectName)
	fmt.Printf("\033[90m➜\033[0m \033[1mcd %s\033[0m\n", projectName)

	if !npm {
		fmt.Printf("\033[90m➜\033[0m \033[1mnpm install && npm run build\033[0m\n")
	}

	fmt.Printf("\033[90m➜\033[0m \033[1mcomposer run dev\033[0m\n")
	fmt.Println()
	fmt.Printf("  New to Laravel? Check out our \033]8;;https://laravel.com/docs/installation#next-steps\033\\documentation\033]8;;\033\\. \033[1mBuild something amazing!\033[0m\n")
	fmt.Println()
}

// Helper function to replace string in file
func replaceInFile(filePath, search, replace string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	newContent := strings.ReplaceAll(string(content), search, replace)
	return os.WriteFile(filePath, []byte(newContent), 0644)
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0644)
}

func askForConfirmation(reader *bufio.Reader, message string) bool {
	fmt.Printf("%s (y/N): ", message)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes"
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
