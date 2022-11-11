package jobhisttab

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/pja237/slurmcommander-dev/internal/keybindings"
)

type Keys map[*key.Binding]bool

var KeyMap = Keys{
	&keybindings.DefaultKeyMap.Up:       true,
	&keybindings.DefaultKeyMap.Down:     true,
	&keybindings.DefaultKeyMap.PageUp:   true,
	&keybindings.DefaultKeyMap.PageDown: true,
	&keybindings.DefaultKeyMap.Slash:    true,
	&keybindings.DefaultKeyMap.Info:     false,
	&keybindings.DefaultKeyMap.Enter:    true,
	&keybindings.DefaultKeyMap.Stats:    true,
}

func (k *Keys) SetupKeys() {
	for k, v := range KeyMap {
		k.SetEnabled(v)
	}
}
