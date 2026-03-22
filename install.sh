#!/bin/bash

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${YELLOW}=== Installing the AmneziaWG CLI Manager (acm) ===${NC}"

if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}Error: Please run the installer as root.${NC}"
  echo -e "Use the command: ${YELLOW}sudo ./install.sh${NC}"
  exit 1
fi

echo -e "\n${GREEN}[1/3] Adding a repository amnezia...${NC}"

add-apt-repository -y ppa:amnezia/ppa

echo -e "\n${GREEN}[2/3] Package updates and installation of amneziawg-tools...${NC}"
apt update
apt install -y amneziawg-tools

echo -e "\n${GREEN}[3/3] Installing the utility in the system...${NC}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY_PATH="$SCRIPT_DIR/acm"
TARGET_PATH="/usr/local/bin/acm"

if [ ! -f "$BINARY_PATH" ]; then
    echo -e "${RED}Error: Binary file 'acm' was not found in the folder $SCRIPT_DIR!${NC}"
    exit 1
fi

cp "$BINARY_PATH" "$TARGET_PATH"
chmod +x "$TARGET_PATH"

echo -e "\n${GREEN}The installation has been completed successfully! 🎉${NC}"
echo -e "Now you can launch the manager from anywhere by simply entering the command:"
echo -e "${YELLOW}acm${NC}\n"