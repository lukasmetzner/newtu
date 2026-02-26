# Newtu

A terminal RSS news reader built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- Aggregate multiple RSS/Atom feeds into a single view
- SQLite-backed article cache for offline browsing
- Search filtering by title
- Jump to any article by row number
- Auto-refresh every 15 minutes

## Install

### From source

```sh
go install github.com/lukasmetzner/newtu@latest
```

### Pre-built binary

Download the latest binary for your platform from the [Releases](https://github.com/lukasmetzner/newtu/releases) page.

## Configuration

On first run, an empty config is created at:

- Linux: `~/.config/newtu/config.json`
- macOS: `~/Library/Application Support/newtu/config.json`

Add your feeds:

```json
{
  "rss_feeds": [
    { "source": "<source>", "url": "<url>" }
  ]
}
```

| Field    | Description                                      |
|----------|--------------------------------------------------|
| `source` | Short label shown in the Source column            |
| `url`    | Full URL to the RSS or Atom feed                  |

## Keybindings

| Key              | Action                                    |
|------------------|-------------------------------------------|
| `Up` / `Down` / `j` / `l` | Navigate articles                         |
| `Enter`          | Open selected article in browser          |
| `/`              | Enter search mode (filter by title)       |
| `Esc`            | Exit search mode and restore full list    |
| `<number>` + `Enter` | Jump to article by row number        |
| `Ctrl+C`         | Quit                                      |

## Data Storage

Articles are cached in a local SQLite database and refreshed from feeds every 15 minutes.

- Linux: `~/.cache/newtu/data.db`
- macOS: `~/Library/Caches/newtu/data.db`

Delete the database file to clear the cache.

## License

[MIT](LICENSE)
