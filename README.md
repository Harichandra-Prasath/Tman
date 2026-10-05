# Tman

**Lightweight Manager for Tmux.**

A terminal UI for browsing and pruning your tmux server. Tman renders the whole
`Server → Session → Window → Pane` hierarchy as a navigable tree compared to 
tmux's own prefix-key interface.

![Tman TUI](assets/screenshot.png)

## What it is

Tman shells out to your local `tmux` binary and mirrors the server as a tree:

```
Tmux Server
├── session: main
│   ├── 1: bash
│   │   ├── 0. bash  ~/projects/Tman
│   │   └── 1. vim   ~/projects/Tman
│   └── 2: server
└── session: logs
    └── 1: tail
```

Nodes are colour-coded by depth — server (green), session (blue), window (aqua),
pane (purple) — and the layout is three panels:

| Panel | Behaviour |
| --- | --- |
| **Tree** | The server hierarchy. Children load the first time you expand a node. |
| **Details** | Everything known about the node under your cursor — created-at, socket path, pane command, and so on. |
| **Events** | The outcome of your last action. |
| **Guide** | The keybindings, so you never have to remember them. |

## Features

- Browse the server, session, window and pane hierarchy from a single tree.
- Lazy expansion — one `tmux` call per level, only when you ask for it.
- Delete any component: the server, a session, a window or a pane.
- Re-sync the whole tree from a live tmux server at any time.
- Per-node details, action feedback and an in-app keybinding guide.
- No config file, no daemon, no background state.

## Requirements

- **A running tmux server with at least one session.** Tman observes an existing
  server; it does not start one.
- **`tmux` on your `PATH`.** All queries go through the tmux CLI, not a socket
  library.
- **Go 1.26.3 or newer** to build from source (pinned in `go.mod`).

## Build and run

```bash
git clone https://github.com/Harichandra-Prasath/Tman.git
cd Tman
make build-binary
```

That produces `target/tman`. Start a session if you do not already have one, then
launch the manager:

```bash
tmux new -s main      # only if no server is running yet
./target/tman
```

`make build-binary` honours overrides, so you can place the binary directly:

```bash
make build-binary TARGET_DIR=~/.local/bin
```

The raw invocation, if you would rather skip the Makefile:

```bash
go build -C cmd/tman -o target/tman
```

## Keybindings

| Key | Scope | Action |
| --- | --- | --- |
| `q` / `Esc` | Global | Quit Tman |
| `r` | Global | Refresh — re-sync the tree from the tmux server |
| `Enter` | Node | Expand / collapse (fetches children on first expand) |
| `d` | Node | Delete the selected component |

Deleting the last child of a node removes that parent from the tree as well, so the
view keeps matching the server.

> **Careful:** `d` on the server node runs `tmux kill-server` — it takes down every
> session, window and pane at once. There is no confirmation prompt.

## How it works

Tman is two small layers with a clean seam between them.

**`pkg/tman` — the domain.** Every entity implements a single interface:

```go
type TmuxComponent interface {
    Name() string          // the tree label
    Details() string       // the Details panel
    GetParent() TmuxComponent
    GetChildCount() int
    SetChildCount(int)
}
```

with four implementations — `Root` (the server), `Session`, `Window` and `Pane` —
each holding a reference back up to its parent.

- `engine.go` issues the queries: `tmux list-sessions`, `tmux list-windows` and
  `tmux list-panes`, each with a `-F` format string that shapes the output into
  colon-delimited fields for the parsers.
- `events.go` maps deletion onto tmux with a type switch: `Root` → `kill-server`,
  `Session` → `kill-session`, `Window` → `kill-window`, `Pane` → `kill-pane`. Targets
  are assembled from the parent chain, so a pane is addressed as
  `session:window.pane`.

**`pkg/ui` — the presentation.** A tview `TreeView` beside three `TextView` panels,
with key handling split into a global capture and a per-node capture.

- Nodes are wrapped in a `TmuxTreeNode`, which pairs a component with its
  `*tview.TreeNode` parent. That parent link is what lets the UI prune the tree
  independently of the domain objects.
- `populate.go` fetches children on selection: expanding a session runs
  `list-windows`, expanding a window runs `list-panes`. A node that already has
  children toggles instead.
- `events.go` implements deletion, recursing upward to drop parents that have just
  lost their last child.

One `tmux` subprocess per expansion, and nothing is cached between them.

## Roadmap

Tman is young, and every rough edge below is something intended to fix. Next major features that are planned to release are  
 
- Switch to diffrent sessions using a key (similar to `tmux switchc -t $TARGET_SESSION`)
- Create Sessions, Windows, Panes under a valid Parent or Root
- General Improvements on Performance (Reducing subprocesses)


## Contributing

Fork, branch, keep the loop green, open a pull request. There is no CI, so this is
your only gate:

```bash
gofmt -l .          # expect no output
go vet ./...
go test ./...       # no tests exist yet
make build-binary
```

## License

MIT © 2026 Harichandra Prasath. See [LICENSE](LICENSE).
