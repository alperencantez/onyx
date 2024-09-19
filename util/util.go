package util

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"

	"net/http"

	"onyx/symlink"
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
	url := fmt.Sprintf("%s/%s/%s", remoteRegistry, packageName, version)
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
		fmt.Printf("error unmarshalling JSON: %v\n Skipping installation", err)
		return "", "", nil
	}

	binPath, ok := packageData["bin"].(map[string]interface{})
	if ok {
		err := os.MkdirAll("./node_modules/.bin", 0755)
		if err != nil {
			fmt.Println("Couldn't create .bin directory skipping")
		}

		for k, v := range binPath {
			pwd, _ := os.Getwd()

			actualBinaryPath := fmt.Sprintf(pwd+"/node_modules/%v/%v", k, v)
			symlink.Create(actualBinaryPath, pwd+"/node_modules/.bin/"+k)
		}
	}

	dist, ok := packageData["dist"].(map[string]interface{})
	if !ok {
		fmt.Printf("dist field is missing or not a map\n Skipping installation")
		return "", "", nil
	}

	tarballURL, ok := dist["tarball"].(string)
	if !ok {
		fmt.Printf("tarball field is missing or not a string\n Skipping installation")
		return "", "", nil
	}

	resolvedVersion, ok := packageData["version"].(string)
	if !ok {
		fmt.Printf("version field is missing or not a string\n Skipping installation")
		return "", "", nil
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

func RemovePackageFromNodeModules(packageName string) error {
	packagePath := filepath.Join("node_modules", packageName)
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		return fmt.Errorf("package %s is not installed in node_modules", packageName)
	}

	err := os.RemoveAll(packagePath)
	if err != nil {
		return fmt.Errorf("failed to remove package from node_modules: %v", err)
	}

	fmt.Printf("Removed %s from node_modules\n", packageName)
	return nil
}

func RemovePackageFromPackageJSON(packageName string) error {
	if _, err := os.Stat("package.json"); os.IsNotExist(err) {
		return fmt.Errorf("package.json not found")
	}

	file, err := os.ReadFile("package.json")
	if err != nil {
		return fmt.Errorf("error reading package.json: %v", err)
	}

	var packageJSON types.PackageJSON
	err = json.Unmarshal(file, &packageJSON)
	if err != nil {
		return fmt.Errorf("error parsing package.json: %v", err)
	}

	removed := false

	if _, exists := packageJSON.Dependencies[packageName]; exists {
		delete(packageJSON.Dependencies, packageName)
		removed = true
	}

	if _, exists := packageJSON.DevDependencies[packageName]; exists {
		delete(packageJSON.DevDependencies, packageName)
		removed = true
	}

	if !removed {
		return fmt.Errorf("package %s not found in dependencies or devDependencies", packageName)
	}

	updatedJSON, err := json.MarshalIndent(packageJSON, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling updated package.json: %v", err)
	}

	err = os.WriteFile("package.json", updatedJSON, 0644)
	if err != nil {
		return fmt.Errorf("error writing updated package.json: %v", err)
	}

	fmt.Printf("Removed %s from package.json\n", packageName)
	return nil
}

func ReadPackageJSON() (*types.PackageJSON, error) {
	file, err := os.ReadFile("package.json")
	if err != nil {
		return nil, fmt.Errorf("error reading package.json: %v", err)
	}

	var packageJSON types.PackageJSON
	err = json.Unmarshal(file, &packageJSON)
	if err != nil {
		return nil, fmt.Errorf("error parsing package.json: %v", err)
	}

	return &packageJSON, nil
}

func RunCustomScript(script string) error {
	cmd := exec.Command("sh", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running script: %v", err)
	}

	return nil
}
