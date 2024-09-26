<h1 align="center">⭕️ Onyx</h1>

<div style="display:flex; justify-content:center">

[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/dwyl/esta/issues)
![Version](https://img.shields.io/npm/v/%40cantez%2Fonyx)
[![Open Source](https://badges.frapsoft.com/os/v1/open-source.svg?v=103)](https://opensource.org/)

</div>

> [!WARNING]  
> Onyx is still in its early development stage and not a production-ready type of software. It might behave in some unexpected ways.

Onyx is a Node.js package manager built with Go (yes, ironically) and the Cobra CLI framework. It allows you to manage your Node.js dependencies without relying on traditional package managers like `npm`, `yarn` or even `pnpm`. The tool provides features for installing, removing, and running scripts from `package.json`, with additional support for handling `devDependencies`.

## Features

- **Install a Single Package**: Install individual Node.js packages directly from the npm registry.
- **Install All Dependencies**: Install all dependencies listed in `package.json`.
- **Install as DevDependencies**: Install a package as a devDependency using a flag.
- **Remove a Package**: Remove a package from `node_modules` and `package.json`.
- **Run Custom Scripts**: Run custom npm scripts defined in the `scripts` section of `package.json`.
- **Install Global Packages**: Install packages globally with permissions management.
- **Graceful Error Handling**: Handles missing or incomplete package metadata gracefully and skips problematic packages.
- **Custom Lockfile**: After a package installation Onyx has a special `.onyxlock.yaml` file in order to track dependency tree in project directory.

## Prerequisites

- [Go](https://golang.org/doc/install) (>=1.20)
- [Git](https://git-scm.com/downloads)
- [Node](https://nodejs.org) (for testing and running Onyx)

## Installation

### Linux / Darwin

1. Clone the repository:

   ```bash
   git clone https://github.com/alperencantez/onyx.git
   ```

2. Change to the project directory:

   ```bash
   cd onyx
   ```

3. Build the CLI tool:

   ```bash
   go build -o .
   ```

4. Move the binary to your `$PATH` for global usage:

   ```bash
   sudo mv onyx /usr/local/bin/
   ```

### Linux

Onyx is distributed via `npm`.
You can run

```bash
npm i @cantez/onyx -g
```

> [!NOTE]  
> Due to a third-party dependency it can only be obtained through `npm`. Do not attempt to get it with a `yarn add`.

Now you can use `onyx` from any directory!

## Usage

### 1. Install a Single Package

To install a package:

```bash
onyx get next 14.1.3
```

To install a package as a devDependency:

```bash
onyx get lodash --dev
```

### 2. Install All Dependencies from `package.json`

If you already have a `package.json` file, you can install all dependencies listed:

```bash
onyx deps
```

This command installs all dependencies and devDependencies listed in the `package.json` file.

### 3. Remove a Package

To remove a package from `node_modules` and `package.json`:

```bash
onyx remove lodash
```

This will delete the `lodash` package from your `node_modules` and remove its entry from the `package.json` file.

### 4. Run a Custom Script

You can run custom scripts defined in the `scripts` section of your `package.json`:

```bash
onyx r build
```

This will execute the `build` script defined in your `package.json`.

### 5. Install Global Packages

To install a package globally:

```bash
sudo onyx get lodash -g
```

This command installs `lodash` globally in your system.

### Error Handling

The tool provides graceful error handling for cases such as:

- Missing `package.json`.
- Missing `node_modules`.
- Packages with incomplete metadata in the npm registry.
- Invalid or malformed versions (e.g., `^` and `~` symbols).

### Logs and Warnings

In cases where a package cannot be installed due to missing metadata or other issues, warnings will be logged, but the process will continue for other packages:

```bash
Warning: 'dist' field is missing for package 'gsap'. Skipping installation.
```

## Development

### Running Locally

1. Clone the repository:

   ```bash
   git clone https://github.com/alperencantez/onyx.git
   ```

2. Make your changes.

3. Run the project locally:

   ```bash
   go run main.go
   ```

### Contributing

If you’d like to contribute, please fork the repository and use a feature branch. Pull requests are welcome!

1. Fork the repo and create your branch:

   ```bash
   git checkout -b feature/my-feature
   ```

2. Make changes and commit them:

   ```bash
   git commit -m "Add new feature"
   ```

3. Push to the branch:

   ```bash
   git push origin feature/my-feature
   ```

4. Create a PR.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
