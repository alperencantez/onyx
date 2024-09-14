package util

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"net/http"

	"onyx/types"

	"os"
	"path/filepath"
	"strings"
)

func Prompt(reader *bufio.Reader, question, defaultValue string) string {
	fmt.Print(question + " (" + defaultValue + ") ")
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)
	if answer == "" {
		return defaultValue
	}
	return answer
}

func GetDefaultPackageName() string {
	dir, err := os.Getwd()
	if err != nil {
		return "@onyx/starter"
	}
	return filepath.Base(dir)
}

func WritePackageJSON(pkg types.PackageJSON) error {
	file, err := os.Create("package.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(pkg)
}

func GetPackageMetadata(packageName, version string, remoteRegistry string) (string, string, error) {
	url := fmt.Sprintf("%s%s/%s", remoteRegistry, packageName, version)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var packageData map[string]interface{}
	if err := json.Unmarshal(body, &packageData); err != nil {
		return "", "", fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	dist, ok := packageData["dist"].(map[string]interface{})
	if !ok {
		return "", "", fmt.Errorf("dist field is missing or not a map")
	}

	tarballURL, ok := dist["tarball"].(string)
	if !ok {
		return "", "", fmt.Errorf("tarball field is missing or not a string")
	}

	resolvedVersion, ok := packageData["version"].(string)
	if !ok {
		return "", "", fmt.Errorf("version field is missing or not a string")
	}

	return tarballURL, resolvedVersion, nil

}

func DownloadAndExtract(url, packageName string, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	nodeModulesPath := filepath.Join(path, packageName)
	os.MkdirAll(nodeModulesPath, 0755)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		relativePath := header.Name
		if len(relativePath) > 8 && relativePath[:8] == "package/" {
			relativePath = relativePath[8:]
		}
		if relativePath == "" {
			continue
		}

		filePath := filepath.Join(nodeModulesPath, relativePath)

		if header.Typeflag == tar.TypeDir {
			os.MkdirAll(filePath, 0755)
		} else {

			if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
				return err
			}
			file, err := os.Create(filePath)
			if err != nil {
				return err
			}
			defer file.Close()
			io.Copy(file, tarReader)
		}
	}

	return nil
}

func UpdatePackageJSON(packageName, version string, isDev bool) error {
	file, err := os.ReadFile("package.json")
	if err != nil {
		return err
	}

	var packageJSON types.PackageJSON

	err = json.Unmarshal(file, &packageJSON)
	if err != nil {
		return err
	}

	if isDev {
		if packageJSON.DevDependencies == nil {
			packageJSON.DevDependencies = make(map[string]string)
		}

		packageJSON.DevDependencies[packageName] = version
	} else {
		if packageJSON.Dependencies == nil {
			packageJSON.Dependencies = make(map[string]string)
		}

		packageJSON.Dependencies[packageName] = version
	}

	updatedJSON, err := json.MarshalIndent(packageJSON, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile("package.json", updatedJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

func InstallGlobally(packageName, version string, remoteRegistry string) {
	var globalDir string
	var globalPathUnix string = "/usr/local/lib/node_modules"
	var globalPathWindows string = "C:\\Program Files\\nodejs\\node_modules"

	if isWindows() {
		globalDir = globalPathWindows
	} else {
		globalDir = globalPathUnix
	}

	fmt.Printf("Installing %s@%s globally...\n", packageName, version)
	tarballURL, resolvedVersion, err := GetPackageMetadata(packageName, version, remoteRegistry)
	if err != nil {
		log.Fatalf("Error fetching metadata for %s: %v", packageName, err)
	}

	err = DownloadAndExtract(tarballURL, packageName, globalDir)
	if err != nil {
		log.Fatalf("Error installing %s: %v", packageName, err)
	}

	fmt.Printf("%s@%s installed globally successfully.\n", packageName, resolvedVersion)

}

func isWindows() bool {
	return strings.HasPrefix(strings.ToLower(os.Getenv("OS")), "windows")
}
