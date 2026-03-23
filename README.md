# 🛡️ acm — AmneziaWG CLI Manager / [RU](assets/README-RU.md)

[![GitHub Release](https://img.shields.io/github/v/release/gsrlabs/acm?style=flat-square&color=blue)](https://github.com/gsrlabs/acm/releases)
[![Language](https://img.shields.io/badge/language-Go-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![OS](https://img.shields.io/badge/OS-Linux-orange?style=flat-square&logo=linux)](https://ubuntu.com/)

**acm** is a lightweight and user-friendly CLI utility written in Go for managing your AmneziaWG (WireGuard) configurations directly in the terminal. Forget about manually typing long paths to config files — manage your connections through an interactive menu.

[How to properly use AmneziaWG 2.0 via console](assets/manual-awg-quick.md)

---

## ✨ Key Features

  * 🚀 **Interactive Menu:** Select the desired config by pressing a single digit.
  * 📂 **Smart Storage:** All config paths are saved in a JSON database (`~/.config/acm/`).
  * 🛠️ **Quick Management:** Add new configs with a single click (`+`) and disconnect all interfaces at once (`-`).
  * 📊 **Monitoring:** Automatic display of `awg show` status directly within the interface.
  * 🎨 **Clean Interface:** Color-coded terminal output for better readability.

---

## 🛠️ System Requirements

  * **OS:** Ubuntu / Debian / Pop\!\_OS (or any distribution with `apt` support).
  * **Dependencies:** The utility will automatically install `amneziawg-tools` if they are not already present in the system.

---

## 🚀 Quick Installation

The easiest way to install `acm` is by running a command in your terminal. The script will automatically add the PPA, update packages, and download the latest binary:

```bash
curl -sSL https://raw.githubusercontent.com/gsrlabs/acm/main/install.sh | sudo bash
```

Once the installation is complete, simply type:

```bash
acm
```

---

## 📖 How to Use

### First Launch

If the configuration list is empty, the program will prompt you to add your first server:

1.  Enter the **absolute path** to the `.conf` file (e.g., `/home/user/vpn/germany.conf`).
2.  Enter a **friendly name** (e.g., `Germany`).

### Main Menu

  * **Digits (1, 2, ...)** — Connect to the selected server.
  * **`+`** — Add another config to the list.
  * **`-`** — Run `down` for all configs in your list (complete VPN disconnection).
  * **`0`** — Exit the program.

---

## ⚙️ Technical Details

  * **Configuration:** Your settings are stored in `~/.config/acm/awg_configs.json`. You can edit it manually if you need to quickly change paths.
  * **Automation:** The utility uses `sudo` to execute `awg-quick` and `awg show` commands. The system may prompt for your password during the first action.

---

## 👨‍💻 Development

If you want to build the project from source:

1.  Clone the repository:
    ```bash
    git clone https://github.com/gsrlabs/acm.git
    ```
2.  Build the binary:
    ```bash
    go build -o acm main.go
    ```

