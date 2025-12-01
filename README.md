# Helm

A minimalistic TUI Pomodoro-like timer designed for pure focus.

## Installation

```bash
curl -sSL https://raw.githubusercontent.com/0xjuanma/helm/main/scripts/install.sh | sh
```

Or with Go:

```bash
go install github.com/0xjuanma/helm@latest
```

Or build from source:

```bash
git clone https://github.com/0xjuanma/helm.git
cd helm
go build
./helm
```

## Usage

Run `helm` to launch the timer interface.

### Controls

| Key | Action |
|-----|--------|
| `j/k` | Navigate |
| `enter` | Select |
| `space` | Start/Pause |
| `r` | Reset |
| `n` | Skip to next step |
| `c` | Customize workflows |
| `esc` | Back |
| `q` | Quit |

### Workflows

- **Pomodoro** - Classic 25/5 minute work/break cycle
- **Design Interview** - Structured interview practice (customizable)
- **Custom** - Create your own workflow

Settings are stored in `~/.helm/settings.json`.

**Author:** [@0xjuanma](https://github.com/0xjuanma)
