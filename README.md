# Laravel CLI

A command-line tool for creating and managing Laravel projects written in Go.

## Features

- Create new Laravel projects with interactive setup
- Configurable app and Vite ports with availability checking
- Interactive setup wizard for common Laravel configuration
- Cross-platform binary distribution
- Package manager support (APT, YUM, etc.)

## Installation

### From GitHub Releases

Download the latest binary for your platform from the [GitHub Releases](https://github.com/your-username/laravel-cli/releases) page.

### Linux (DEB-based distributions)

```bash
# Download and install the DEB package
wget https://github.com/your-username/laravel-cli/releases/latest/download/laravel-cli_VERSION_amd64.deb
sudo dpkg -i laravel-cli_VERSION_amd64.deb
```

### Linux (RPM-based distributions)

```bash
# Download and install the RPM package
wget https://github.com/your-username/laravel-cli/releases/latest/download/laravel-cli-VERSION-1.x86_64.rpm
sudo rpm -i laravel-cli-VERSION-1.x86_64.rpm
```

### From Source

```bash
git clone https://github.com/your-username/laravel-cli.git
cd laravel-cli
go build -o laravel
sudo mv laravel /usr/local/bin/
```

## Usage

### Create a New Laravel Project

```bash
laravel new my-project
```

This will:
1. Clone the Laravel repository into a new directory called `my-project`
2. Copy `.env.example` to `.env`
3. Run an interactive setup wizard that asks for:
   - App Port (default: 80) - with port availability checking
   - Vite Port (default: 5173) - with port availability checking
   - App Name (default: "My Laravel Application")

### Get Help

```bash
laravel --help
```

## Interactive Setup

The setup wizard will ask you to configure:

- **App Port**: The port your Laravel application will run on (default: 80)
- **Vite Port**: The port for Vite development server (default: 5173)
- **App Name**: The name of your application (default: "My Laravel Application")

The tool automatically checks port availability and will prompt you to choose a different port if the specified one is already in use.

You can press `Ctrl+C` at any time during setup to exit with a "Goodbye!" message.

## Requirements

- Git (for cloning Laravel repository)
- Composer (for Laravel dependencies - run after project creation)
- PHP (for running Laravel)

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

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
