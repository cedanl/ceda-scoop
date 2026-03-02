# r-scoop-template

**R Project Template with Scoop Installation**

A template repository for R projects that includes a [Scoop](https://scoop.sh)-based PowerShell install script for easy local setup on Windows.

## 📋 Table of Contents

- [Quick Start](#-quick-start)
- [Project Structure](#-project-structure)
- [Prerequisites](#-prerequisites)
- [Installation](#-installation)
- [Usage](#-usage)
- [Contributing](#-contributing)
- [License](#-license)

## 🚀 Quick Start

1. Click **"Use this template"** at the top of this page to create your own repository.
2. Clone your new repository:
   ```powershell
   git clone https://github.com/yourusername/your-repo-name.git
   cd your-repo-name
   ```
3. Run the installation script (PowerShell):
   ```powershell
   .\install.ps1
   ```
4. Open `r-scoop-template.Rproj` in RStudio and start coding.

## 📁 Project Structure

```
r-scoop-template/
├── R/                     # R functions and scripts
│   └── hello.R            # Example function
├── .Rprofile              # R profile with renv settings
├── .gitignore             # R-appropriate gitignore
├── DESCRIPTION            # R package metadata
├── install.ps1            # Scoop-based Windows installation script
├── renv.lock              # Reproducible package lockfile (renv)
└── r-scoop-template.Rproj # RStudio project file
```

## ✅ Prerequisites

- Windows 10 or later
- PowerShell 5.1 or PowerShell 7+
- Internet connection (for downloading Scoop, R, and packages)

## 📦 Installation

The `install.ps1` script automates the full setup:

| Step | What it does |
|------|-------------|
| 1 | Installs [Scoop](https://scoop.sh) if not present |
| 2 | Adds the `main` and `extras` Scoop buckets |
| 3 | Installs **R** via `scoop install r` |
| 4 | Installs **Rtools** via `scoop install rtools` (needed to build packages from source) |
| 5 | Restores R package dependencies with `renv::restore()` |

Run it from a PowerShell terminal in the project root:

```powershell
.\install.ps1
```

> **Note:** The script may prompt you to set the execution policy for the current user (`RemoteSigned`) the first time you run it.

## 💻 Usage

After installation, open `r-scoop-template.Rproj` in RStudio.

Customize this template for your own project:

1. Edit `DESCRIPTION` to update package name, title, and author details.
2. Add your R functions under `R/`.
3. Add packages to `renv.lock` by running `renv::snapshot()` after installing them.

## 🤝 Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## 📄 License

This project is licensed under the [Apache License 2.0](LICENSE).
