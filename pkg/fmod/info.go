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

func (fmi *FactorioModInfo) ToZip(from, to string) (original uint64, compressed uint64, err error) {
	fpath := filepath.Join(to, fmi.ToZipName())

	archive, err := os.Create(fpath)

	if err != nil {
		return 0, 0, err
	}

	defer archive.Close()

	w := zip.NewWriter(archive)
	defer w.Close()

	total, err := addFilesToZip(w, from, fmi.Name)

	if err != nil {
		return 0, 0, err
	}

	if err := w.Close(); err != nil {
		return 0, 0, errors.New("Warning: closing zipfile writer failed: " + err.Error())
	}

	info, err := archive.Stat()

	if err != nil {
		return 0, 0, err
	}

	return total, uint64(info.Size()), nil
}

func addFilesToZip(w *zip.Writer, basePath, baseInZip string) (uint64, error) {
	total := uint64(0)

	files, err := os.ReadDir(basePath)

	if err != nil {
		return 0, err
	}

	for _, file := range files {
		fullfilepath := filepath.Join(basePath, file.Name())

		fullfilepath = filepath.ToSlash(fullfilepath)

		if _, err := os.Stat(fullfilepath); os.IsNotExist(err) {
			// ensure the file exists. For example a symlink pointing to a non-existing location might be listed but not actually exist
			continue
		}

		if file.Type()&os.ModeSymlink != 0 {
			// ignore symlinks alltogether
			continue
		}

		if file.IsDir() {
			dirTotal, err := addFilesToZip(w, fullfilepath, filepath.Join(baseInZip, file.Name()))

			if err := err; err != nil {
				return 0, err
			}

			total += dirTotal
		} else if file.Type().IsRegular() {
			dat, err := os.ReadFile(fullfilepath)

			if err != nil {
				return 0, err
			}

			fp := filepath.ToSlash(filepath.Join(baseInZip, file.Name()))

			f, err := w.Create(fp)

			if err != nil {
				return 0, err
			}

			file_total, err := f.Write(dat)

			if err != nil {
				return 0, err
			}

			total += uint64(file_total)
		} else {
			// we ignore non-regular files because they are scary
		}
	}

	return total, nil
}
