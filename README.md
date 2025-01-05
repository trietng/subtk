# subtk

subtk is a command-line toolkit for downloading and manipulating subtitles for movies and TV shows.

## Use-cases

At the moment, subtk supports the following use-cases:

- Download subtitles for movies and TV shows.
- Search for subtitles for movies and TV shows.
- Set preferred language for subtitles.

## Specification

### Configuration

The program reserved the folder `~/.subtk` for storing configuration files and resources.

Configuration file is stored at `~/.subtk/config.json`.Example configuration file:

```json
{
    "default_language": "en",
    "api_keys": {
        "subdl": "<API_KEY>",
    }
}
```

## Supported subtitle services

- [subdl](https://subdl.com/)

## Installation

### Standalone executable

Download the executable and place it in your `%PATH%` (Windows-only).

### Go installation

You can also download the source code release and install it using standard Go command `go install`.

## License

Copyright (c) 2025 trietng

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for details.