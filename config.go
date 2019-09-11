package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
)

var Config cfg

type cfg struct {
	Keybindings map[string]string
	Gap         int
	Proportion  float64
	HideDecor   bool `toml:"remove_decorations"`
}

func init() {
	writeDefaultConfig()
	toml.DecodeFile(configFilePath(), &Config)
}

func writeDefaultConfig() {
	if _, err := os.Stat(configFolderPath()); os.IsNotExist(err) {
		os.MkdirAll(configFolderPath(), 0700)
	}

	if _, err := os.Stat(configFilePath()); os.IsNotExist(err) {
		ioutil.WriteFile(configFilePath(), []byte(defaultConfig), 0644)
	}
}

func configFolderPath() string {
	var configFolder string
	switch runtime.GOOS {
	case "linux":
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome != "" {
			configFolder = xdgConfigHome
		} else {
			configFolder, _ = homedir.Expand("~/.config/zentile/")
		}
	default:
		configFolder, _ = homedir.Expand("~/.zentile/")
	}

	return configFolder
}

func configFilePath() string {
	return filepath.Join(configFolderPath(), "config.toml")
}

var defaultConfig = `# Window decorations will be removed if set to true
remove_decorations = false

# Adds spacing between windows
gap = 5

# How much to increment the master area size.
proportion = 0.1

[keybindings]

# Tile the current workspace.
tile = "Control-Shift-t"

# Untile the current workspace.
untile = "Control-Shift-u"

# Make the active window as master.
make_active_window_master = "Control-Shift-m"

# Increase number of masters.
increase_master = "Control-Shift-i"

# Decrease the number of masters.
decrease_master = "Control-Shift-d"

# Cycles through the available layouts.
switch_layout = "Control-Shift-s"

# Moves focus to the next window.
next_window = "Control-Shift-n"

# Moves focus to the previous window.
previous_window = "Control-Shift-p"

# Increases the size of the master windows.
increment_master = "Control-bracketright"

# Decreases the size of the master windows.
decrement_master = "Control-bracketleft"
`
