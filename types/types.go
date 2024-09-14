package types

type PackageJSON struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Main        string `json:"main"`
	Scripts     struct {
		Test string `json:"test"`
	} `json:"scripts"`
	Author          string            `json:"author"`
	License         string            `json:"license"`
	Dependencies    map[string]string `json:"dependencies,omitempty"`
	DevDependencies map[string]string `json:"devDependencies,omitempty"`
}
