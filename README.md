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
sudo make install
```

## Usage

To display events for today, simply run:

```sh
wday
```

To display events for a specific date, use the `--date` flag:

```sh
wday --date MM-DD
```

**example**

```sh
wday --date 01-01

元日
```

To show all events found for a specific date, use the `--all` flag:

```sh
wday --date MM-DD --all
```

**example**

```sh
wday --date 01-01 --all

元日
太陽暦施行の日
神戸港記念日
中華民国設立記念日
少年法施行の日
鉄腕アトムの日
日本初の点字新聞「あけぼの」創刊記念日
オールインワンゲルの日
スカルプDの発毛DAY
肉汁水餃子の日
資格チャレンジの日
釜飯の日
あずきの日
省エネルギーの日
Myハミガキの日
もったいないフルーツの日
```

## Available Locales

To list all supported locales, use the following command:

```sh
wday locale list
```

To set a specific locale, use the `set` command followed by the locale code:

```sh
wday locale set [locale code]
```

Supported locales include 'EnUS' for English (US) and 'JaJP' for Japanese.

## Configuration

The application can be configured via a YAML file located at `$HOME/.wday.yaml`. The default configuration is:

```yaml
locale: JaJP
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
