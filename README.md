# whatday

whatday is a simple CLI tool that reveals historical events, notable birthdays, and interesting observances for any given day. Built with Go and Cobra, this project serves as a personal learning exercise in crafting concise, effective CLI applications.

## Features
- Display historical events, notable birthdays, and observances for any given day.
- Supports multiple locales (English and Japanese).
- Configurable via a YAML file.

## Installation

To install the CLI tool, run the following commands:

```sh
# Clone the repository
git clone https://github.com/yourusername/whatday.git

# Navigate to the project directory
cd whatday

# Build the project
make build

# Optionally, install the binary to your system
sudo install wday /usr/local/bin/
```

## Usage

To display events for today, simply run:

```sh
wday
```

To display events for a specific date, use the `--date` flag:

```sh
wday --date YYYY-MM-DD
```

To show all events found for a specific date, use the `--all` flag:

```sh
wday --date YYYY-MM-DD --all
```

## Configuration

The application can be configured via a YAML file located at `$HOME/.wday.yaml`. The default configuration is:

```yaml
locale: JaJP
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
