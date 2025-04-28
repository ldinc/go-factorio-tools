package manager

import "fmt"

func (m *Mod) GetZipString() string {
	return m.Info.ToZipName()
}

func (m *Manager) BuildAll() {
	for _, mod := range m.Mods {
		zipName := mod.GetZipString()

		fmt.Printf("|> build release [%s]\n\t --- archive: %s\n", mod.Info.Name, zipName)

		err := mod.Info.ToZip(mod.Dir, m.TargetDir)

		if err != nil {
			panic(err)
		}

		fmt.Printf("\t --- location: [%s]\n", m.TargetDir)
	}
}
