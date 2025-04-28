package fmod

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FactorioModInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (fmi *FactorioModInfo) ToString() string {
	if fmi == nil {
		return "<nil>"
	}

	return fmt.Sprintf("mod: %s [version: %s]", fmi.Name, fmi.Version)
}

func (fmi *FactorioModInfo) ToZipName() string {
	if fmi == nil {
		return "nil.zip"
	}

	return fmt.Sprintf("%s_%s.zip", fmi.Name, fmi.Version)
}

func GetInfoFrom(path string) (*FactorioModInfo, error) {
	filepath := filepath.Join(path, "info.json")

	_, err := os.Stat(filepath)

	if errors.Is(err, os.ErrNotExist) {
		return nil, errors.Join(ErrInvalidMod, err)
	}

	jsonFile, err := os.Open(filepath)

	if err != nil {
		return nil, errors.Join(ErrInvalidMod, err)
	}

	defer jsonFile.Close()

	bdata, err := io.ReadAll(jsonFile)

	if err != nil {
		return nil, errors.Join(ErrInvalidMod, err)
	}

	info := new(FactorioModInfo)

	err = json.Unmarshal(bdata, &info)

	if err != nil {
		return nil, errors.Join(ErrInvalidMod, err)
	}

	return info, nil
}

func (fmi *FactorioModInfo) ToZip(from, to string) error {
	fpath := filepath.Join(to, fmi.ToZipName())

	archive, err := os.Create(fpath)

	if err != nil {
		return err
	}

	defer archive.Close()

	w := zip.NewWriter(archive)
	defer w.Close()

	if err := addFilesToZip(w, from, fmi.Name); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return errors.New("Warning: closing zipfile writer failed: " + err.Error())
	}

	return nil
}

func addFilesToZip(w *zip.Writer, basePath, baseInZip string) error {
	files, err := os.ReadDir(basePath)

	if err != nil {
		return err
	}

	for _, file := range files {
		fullfilepath := filepath.Join(basePath, file.Name())

		if _, err := os.Stat(fullfilepath); os.IsNotExist(err) {
			// ensure the file exists. For example a symlink pointing to a non-existing location might be listed but not actually exist
			continue
		}

		if file.Type()&os.ModeSymlink != 0 {
			// ignore symlinks alltogether
			continue
		}

		if file.IsDir() {
			if err := addFilesToZip(w, fullfilepath, filepath.Join(baseInZip, file.Name())); err != nil {
				return err
			}
		} else if file.Type().IsRegular() {
			dat, err := os.ReadFile(fullfilepath)

			if err != nil {
				return err
			}

			f, err := w.Create(filepath.Join(baseInZip, file.Name()))

			if err != nil {
				return err
			}

			_, err = f.Write(dat)

			if err != nil {
				return err
			}
		} else {
			// we ignore non-regular files because they are scary
		}
	}
	return nil
}
