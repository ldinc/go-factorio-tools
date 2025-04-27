package manager

import "fmt"

func (m *Mod) GetZipString() string {
	return m.Info.ToZipName()
}

func (m *Manager) BuildAll() {
	for _, mod := range m.Mods {
		zipName := mod.GetZipString()

		fmt.Println("make archive:", zipName)

		err := mod.Info.ToZip(mod.Dir, m.TargetDir)

		if err != nil {
			panic(err)
		}

		fmt.Printf("make archive: moved to -> [%s]\n", m.TargetDir)
	}
}
