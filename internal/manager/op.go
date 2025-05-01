package manager

import (
	"fmt"

	"github.com/dustin/go-humanize"
)

func (m *Mod) GetZipString() string {
	return m.Info.ToZipName()
}

func (m *Manager) BuildAll() {
	for _, mod := range m.Mods {
		zipName := mod.GetZipString()

		fmt.Printf("|> build release [%s]\n\t --- archive: %s\n", mod.Info.Name, zipName)

		total, err := mod.Info.ToZip(mod.Dir, m.TargetDir)

		if err != nil {
			panic(err)
		}

		fmt.Printf("\t --- size: %s\n", humanize.Bytes(total))
		fmt.Printf("\t --- location: [%s]\n", m.TargetDir)
	}
}
