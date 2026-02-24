# notion-tui

A terminal UI for browsing your Notion workspace, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- Browse all pages and databases in your workspace
- Fuzzy-filter the list by typing
- View page content as rendered markdown
- Search across your workspace with `/`
- In-memory cache for instant repeated views

## Setup

### 1. Create a Notion Integration

1. Go to https://www.notion.so/my-integrations
2. Create a new integration
3. Copy the "Internal Integration Secret"
4. Share your pages/databases with the integration (in Notion, open a page → ··· → Connections → Add your integration)

### 2. Configure the token

Either set an environment variable:

```bash
export NOTION_API_TOKEN=ntn_xxxx...
```

Or create a config file at `~/.notion-tui/config.yaml`:

```yaml
token: ntn_xxxx...
```

### 3. Build and run

```bash
make build
./notion-tui
```

Or install to your Go bin:

```bash
make install
notion-tui
```

## Key Bindings

| Key | Action |
|---|---|
| Type | Fuzzy-filter list |
| Enter | Open selected page |
| Esc | Go back |
| / | Search mode |
| j/k | Scroll in viewer |
| q / Ctrl+C | Quit |

## Requirements

- Go 1.21+
- A Notion integration token with access to your workspace
