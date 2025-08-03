# Laravel CLI

A command-line tool for creating and managing Laravel projects written in Go. This tool replicates the full functionality of the official Laravel installer (`composer global require laravel/installer`).

## Features

- **Complete Laravel installer replication** - All features from the PHP installer
- **Command-line options** - Extensive flags for automated project setup
- **Starter kit support** - React, Vue, Livewire starter kits
- **Database configuration** - Support for MySQL, MariaDB, PostgreSQL, SQLite, SQL Server
- **Testing framework setup** - Pest vs PHPUnit integration
- **Git repository initialization** - Optional Git and GitHub integration  
- **NPM dependency management** - Automatic npm install and build
- **Interactive setup wizard** - User-friendly prompts for configuration
- **Cross-platform binary distribution** - Native binaries for all platforms
- **Package manager support** - APT, YUM, etc.

## Installation

### From GitHub Releases

Download the latest binary for your platform from the [GitHub Releases](https://github.com/protocol-seven/laravel-cli/releases) page.

### Linux (DEB-based distributions)

```bash
# Download and install the DEB package
wget https://github.com/protocol-seven/laravel-cli/releases/latest/download/laravel-cli_VERSION_amd64.deb
sudo dpkg -i laravel-cli_VERSION_amd64.deb
```

### Linux (RPM-based distributions)

```bash
# Download and install the RPM package
wget https://github.com/protocol-seven/laravel-cli/releases/latest/download/laravel-cli-VERSION-1.x86_64.rpm
sudo rpm -i laravel-cli-VERSION-1.x86_64.rpm
```

### From Source

```bash
git clone https://github.com/protocol-seven/laravel-cli.git
cd laravel-cli
go build -o laravel
sudo mv laravel /usr/local/bin/
```

## Usage

### Basic Project Creation

```bash
laravel new my-project
```

### Command-Line Options

```bash
# Development release
laravel new my-project --dev

# Initialize Git repository
laravel new my-project --git

# Create with specific Git branch
laravel new my-project --git --branch=develop

# Create GitHub repository
laravel new my-project --github

# Create GitHub repository in organization
laravel new my-project --github --organization=my-org

# Specify database driver
laravel new my-project --database=mysql
laravel new my-project --database=pgsql
laravel new my-project --database=sqlite

# Install starter kits
laravel new my-project --react
laravel new my-project --vue
laravel new my-project --livewire

# Custom starter kit
laravel new my-project --using=vendor/package

# Testing frameworks
laravel new my-project --pest
laravel new my-project --phpunit

# Install and build NPM dependencies
laravel new my-project --npm

# Force overwrite existing directory
laravel new my-project --force

# Quiet mode (suppress output)
laravel new my-project --quiet

# Combine multiple options
laravel new my-project --git --github --database=mysql --pest --npm
```

### Available Database Drivers

- `mysql` - MySQL
- `mariadb` - MariaDB  
- `pgsql` - PostgreSQL
- `sqlite` - SQLite (default)
- `sqlsrv` - SQL Server

### Starter Kits

- `--react` - Laravel + React starter kit
- `--vue` - Laravel + Vue starter kit  
- `--livewire` - Laravel + Livewire starter kit
- `--using=package` - Custom community starter kit

## Interactive Setup

When you create a new project, the tool will guide you through:

1. **Project Creation** - Uses Composer to create the Laravel project
2. **Environment Setup** - Copies `.env.example` to `.env`
3. **Database Configuration** - Interactive database driver selection
4. **Application Configuration** - App URL and other settings
5. **Database Migration** - Optional database migration
6. **NPM Dependencies** - Optional npm install and build

The tool automatically:
- Runs `composer install`
- Generates application key with `php artisan key:generate`
- Sets proper file permissions
- Configures environment variables

## Requirements

- **Composer** - For Laravel project creation
- **PHP** - For running Laravel and Artisan commands
- **Git** - For repository initialization (optional)
- **GitHub CLI** - For GitHub integration (optional)
- **Node.js/NPM** - For frontend dependencies (optional)

## Examples

### Simple Laravel Project
```bash
laravel new blog
```

### Full-Featured Project with All Options
```bash
laravel new my-app \
  --git \
  --github \
  --database=mysql \
  --vue \
  --pest \
  --npm
```

### Quick SQLite Project
```bash
laravel new quick-project --database=sqlite --pest
```

### React Project with GitHub
```bash
laravel new react-app --react --git --github --npm
```

## Development

### Building

```bash
go build -o laravel
```

### Running Tests

```bash
go test ./...
```

### Cross-Platform Build

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o laravel-linux-amd64

# Windows  
GOOS=windows GOARCH=amd64 go build -o laravel-windows-amd64.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o laravel-darwin-amd64
```

## Comparison with PHP Laravel Installer

This Go implementation provides **100% feature parity** with the official PHP Laravel installer:

| Feature | PHP Installer | Go CLI | Status |
|---------|---------------|--------|---------|
| Basic project creation | ✅ | ✅ | ✅ Complete |
| Command-line flags | ✅ | ✅ | ✅ Complete |
| Starter kits (React/Vue/Livewire) | ✅ | ✅ | ✅ Complete |
| Database configuration | ✅ | ✅ | ✅ Complete |
| Testing framework setup | ✅ | ✅ | ✅ Complete |
| Git repository initialization | ✅ | ✅ | ✅ Complete |
| GitHub integration | ✅ | ✅ | ✅ Complete |
| NPM dependency management | ✅ | ✅ | ✅ Complete |
| Interactive prompts | ✅ | ✅ | ✅ Complete |
| Force overwrite | ✅ | ✅ | ✅ Complete |
| Quiet mode | ✅ | ✅ | ✅ Complete |
| Custom starter kits | ✅ | ✅ | ✅ Complete |

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
