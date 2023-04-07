// Copyright 2023 Krisna Pranav, Sankar-2006. All rights reserved.
// Use of this source code is governed by a Apache-2.0 License
// license that can be found in the LICENSE file

// updater
package updater

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/axllent/semver"
	"github.com/krishpranav/Mailtrix/utils/logger"
)

var (
	AllowPrereleases = false

	tempDir string
)

type Releases []struct {
	Name       string `json:"name"`
	Tag        string `json:"tag_name"`
	Prerelease bool   `json:"prerelease"`
	Assets     []struct {
		BrowserDownloadURL string `json:"browser_download_url"`
		ID                 int64  `json:"id"`
		Name               string `json:"name"`
		Size               int64  `json:"size"`
	} `json:"assets"`
}

type Release struct {
	Name string
	Tag  string
	URL  string
	Size int64
}

func GithubLatest(repo, name string) (string, string, string, error) {
	releaseURL := fmt.Sprintf("https://api.github.com/repos/%s/releases", repo)

	resp, err := http.Get(releaseURL)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", "", "", err
	}

	linkOS := runtime.GOOS
	linkArch := runtime.GOARCH
	linkExt := ".tar.gz"
	if linkOS == "windows" {
		linkExt = ".zip"
	}

	var allReleases = []Release{}

	var releases Releases

	if err := json.Unmarshal(body, &releases); err != nil {
		return "", "", "", err
	}

	archiveName := fmt.Sprintf("%s-%s-%s%s", name, linkOS, linkArch, linkExt)

	for _, r := range releases {
		if !semver.IsValid(r.Tag) {
			continue
		}

		if !AllowPrereleases && (semver.Prerelease(r.Tag) != "" || r.Prerelease) {
			continue
		}

		for _, a := range r.Assets {
			if a.Name == archiveName {
				thisRelease := Release{a.Name, r.Tag, a.BrowserDownloadURL, a.Size}
				allReleases = append(allReleases, thisRelease)
				break
			}
		}
	}

	if len(allReleases) == 0 {
		return "", "", "", fmt.Errorf("No binary releases found")
	}

	var latestRelease = Release{}

	for _, r := range allReleases {
		if semver.Compare(r.Tag, latestRelease.Tag) == 1 {
			latestRelease = r
		}
	}

	return latestRelease.Tag, latestRelease.Name, latestRelease.URL, nil
}

func GreaterThan(toVer, fromVer string) bool {
	return semver.Compare(toVer, fromVer) == 1
}

func GithubUpdate(repo, appName, currentVersion string) (string, error) {
	ver, filename, downloadURL, err := GithubLatest(repo, appName)

	if err != nil {
		return "", err
	}

	if ver == currentVersion {
		return "", fmt.Errorf("No new release found")
	}

	if semver.Compare(ver, currentVersion) < 1 {
		return "", fmt.Errorf("No newer releases found (latest %s)", ver)
	}

	tmpDir := getTempDir()

	outFile := filepath.Join(tmpDir, filename)

	if err := downloadToFile(downloadURL, outFile); err != nil {
		return "", err
	}

	newExec := filepath.Join(tmpDir, "Mailtrix")

	if runtime.GOOS == "windows" {
		if _, err := Unzip(outFile, tmpDir); err != nil {
			return "", err
		}
		newExec = filepath.Join(tmpDir, "Mailtrix.exe")
	} else {
		if err := TarGZExtract(outFile, tmpDir); err != nil {
			return "", err
		}
	}

	if runtime.GOOS != "windows" {
		if err := os.Chmod(newExec, 0755); err != nil {
			return "", err
		}
	}

	cmd := exec.Command(newExec, "-h")
	if err := cmd.Run(); err != nil {
		return "", err
	}

	oldExec, err := os.Executable()
	if err != nil {
		return "", err
	}

	if err = replaceFile(oldExec, newExec); err != nil {
		return "", err
	}

	return ver, nil
}

func downloadToFile(url, fileName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath.Clean(fileName))
	if err != nil {
		return err
	}

	defer func() {
		if err := out.Close(); err != nil {
			logger.Log().Errorf("Error closing file: %s\n", err)
		}
	}()

	_, err = io.Copy(out, resp.Body)

	return err
}

func replaceFile(dst, src string) error {
	source, err := os.Open(filepath.Clean(src))
	if err != nil {
		return err
	}

	dstDir := filepath.Dir(dst)
	binaryFilename := filepath.Base(dst)
	dstOld := fmt.Sprintf("%s.old", binaryFilename)
	dstNew := fmt.Sprintf("%s.new", binaryFilename)
	newTmpAbs := filepath.Join(dstDir, dstNew)
	oldTmpAbs := filepath.Join(dstDir, dstOld)

	fi, _ := os.Stat(dst)
	srcPerms := fi.Mode().Perm()

	tmpNew, err := os.OpenFile(filepath.Clean(newTmpAbs), os.O_CREATE|os.O_RDWR, srcPerms)
	if err != nil {
		return err
	}

	if _, err := io.Copy(tmpNew, source); err != nil {
		return err
	}

	if err := tmpNew.Close(); err != nil {
		return err
	}

	if err := source.Close(); err != nil {
		return err
	}

	if err := os.Rename(dst, oldTmpAbs); err != nil {
		return err
	}

	if err := os.Rename(newTmpAbs, dst); err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		tmpDir := os.TempDir()
		delFile := filepath.Join(tmpDir, filepath.Base(oldTmpAbs))
		if err := os.Rename(oldTmpAbs, delFile); err != nil {
			return err
		}
	} else {
		if err := os.Remove(oldTmpAbs); err != nil {
			return err
		}
	}

	if err := os.Remove(src); err != nil {
		return err
	}

	return nil
}

func getTempDir() string {
	if tempDir == "" {
		randBytes := make([]byte, 6)
		if _, err := rand.Read(randBytes); err != nil {
			panic(err)
		}
		tempDir = filepath.Join(os.TempDir(), "updater-"+hex.EncodeToString(randBytes))
	}
	if err := mkDirIfNotExists(tempDir); err != nil {
		logger.Log().Errorf("Error: %v", err)
		os.Exit(2)
	}

	return tempDir
}

func mkDirIfNotExists(path string) error {
	if !isDir(path) {
		return os.MkdirAll(path, os.ModePerm)
	}

	return nil
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) || !info.Mode().IsRegular() {
		return false
	}

	return true
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) || !info.IsDir() {
		return false
	}

	return true
}
