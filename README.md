# Cortile [WIP]
Tiling manager with hot corner support for Xfce, OpenBox and other [EWMH Compliant Window Managers](https://en.m.wikipedia.org/wiki/Extended_Window_Manager_Hints).

## Features
- Workspace based tiling.
- Keyboard and hot corner events.
- Vertical and horizontal tiling.
- Resize of master/slave area.
- Auto detection of panels.
- Multi monitor support.
- Customizable layouts.

## Install
### Requirements
Install [go](https://go.dev/) with `apt`:
```bash
sudo apt install golang
```

Install [go](https://go.dev/) with `pacman`:
```bash
sudo pacman -S go
```

### Use remote source
Install [cortile](https://github.com/leukipp/cortile) from GitHub:
```bash
go get -u github.com/leukipp/cortile
go install github.com/leukipp/cortile
```

### Use local source
Clone [cortile](https://github.com/leukipp/cortile) from GitHub:
```bash
git clone https://github.com/leukipp/cortile.git
cd cortile
```

Once you have made local changes run:
```bash
go build
go install
```

## Run
Start in verbose mode:
```bash
cortile -v
```

Window resizing can be done with <kbd>Alt</kbd>+<kbd>Right-Click</kbd> on Xfce.

## Config
The config file is located at `~/.config/cortile/config.toml`.

| Corner events                      | Description                       |
| ---------------------------------- | --------------------------------- |
| <kbd>top</kbd>-<kbd>left</kbd>     | Cycle through layouts             |
| <kbd>top</kbd>-<kbd>right</kbd>    | Make the active window master     |
| <kbd>bottom</kbd>-<kbd>right</kbd> | Increase number of master windows |
| <kbd>bottom</kbd>-<kbd>left</kbd>  | Decrease number of master windows |

| Key events                                    | Description                       |
| --------------------------------------------- | --------------------------------- |
| <kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>t</kbd> | Tile current workspace            |
| <kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>u</kbd> | Untile current workspace          |
| <kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>m</kbd> | Make the active window master     |
| <kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>i</kbd> | Increase number of master windows |
| <kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>d</kbd> | Decrease number of master windows |
| <kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>s</kbd> | Cycle through layouts             |
| <kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>n</kbd> | Goto next window                  |
| <kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>p</kbd> | Goto previous window              |
| <kbd>Ctrl</kbd>+<kbd>]</kbd>                  | Increase size of master windows   |
| <kbd>Ctrl</kbd>+<kbd>[</kbd>                  | Decrease size of master windows   |

## WIP
- Configurable hot corners (current: hardcoded corner events).
- Configurable LTR/RTL support (current: master is on the right side).
- Proper dual monitor support (current: only biggest monitor is tiled).
- Resizable windows (current: only master/slave proportion can be changed).

## Credits
Based on [zentile](https://github.com/blrsn/zentile) from [Berin Larson](https://github.com/blrsn/).

## License
[MIT](/LICENSE)
