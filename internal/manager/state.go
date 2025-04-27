package manager

import (
	"fmt"
	"goft/pkg/fmod"
	"os"
	"path/filepath"
)

type Mod struct {
	Info *fmod.FactorioModInfo
	Dir  string
}

type Manager struct {
	Mods      []*Mod
	TargetDir string
}

func New(target string, output string, mods []string) *Manager {

	targetPath, err := filepath.Abs(target)

	if err != nil {
		panic(err)
	}

	outputPath, err := filepath.Abs(output)

	if err != nil {
		panic(err)
	}

	_, err = os.Stat(outputPath)

	if err != nil {
		panic(err)
	}

	sourceDirs := make([]string, 0, len(mods))

	if len(mods) > 0 {
		for _, mod := range mods {
			val, ok := checkPath(targetPath, mod)
			if ok {
				sourceDirs = append(sourceDirs, val)
			} else {
				fmt.Printf("error: path [%s] with [%s] was ignored\n", val, mod)
			}
		}
	} else {
		val, _ := checkPath(targetPath, "src")
		sourceDirs = append(sourceDirs, val)
	}

	fmt.Println(sourceDirs)

	list := make([]*Mod, 0, len(sourceDirs))

	for _, sourceDir := range sourceDirs {
		mod := getMod(sourceDir)

		if mod != nil {
			list = append(list, mod)
		}
	}

	return &Manager{
		Mods:      list,
		TargetDir: outputPath,
	}
}

func getMod(path string) *Mod {
	mod, err := fmod.GetInfoFrom(path)

	if err != nil {
		fmt.Printf("error: path [%s] was ignored\n", path)

		return nil
	}

	return &Mod{
		Info: mod,
		Dir:  path,
	}
}

func checkPath(path, mod string) (string, bool) {
	filepath := filepath.Join(path, mod)
	_, err := os.Stat(filepath)

	if err != nil {
		return path, false
	}

	return filepath, true
}
