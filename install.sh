#!/bin/bash

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${YELLOW}=== Installing the AmneziaWG CLI Manager (acm) ===${NC}"

if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}Error: Please run the installer as root.${NC}"
  echo -e "If you are using the one-liner, make sure to pipe to 'sudo bash'."
  exit 1
fi

echo -e "\n${GREEN}[1/3] Adding the amnezia repository...${NC}"
add-apt-repository -y ppa:amnezia/ppa

echo -e "\n${GREEN}[2/3] Updating packages and installing dependencies...${NC}"
apt update

apt install -y amneziawg-tools curl

echo -e "\n${GREEN}[3/3] Downloading and installing the acm utility...${NC}"

TARGET_PATH="/usr/local/bin/acm"
DOWNLOAD_URL="https://github.com/gsrlabs/acm/releases/download/v1.0.0/acm"

echo "Downloading from GitHub Releases..."
curl -sSL "$DOWNLOAD_URL" -o "$TARGET_PATH"

if [ ! -s "$TARGET_PATH" ]; then
    echo -e "${RED}Error: Failed to download the 'acm' binary! Please check the release URL.${NC}"
    rm -f "$TARGET_PATH"
    exit 1
fi

chmod +x "$TARGET_PATH"

echo -e "\n${GREEN}The installation has been completed successfully! 🎉${NC}"
echo -e "Now you can launch the manager from anywhere by simply entering the command:"
echo -e "${YELLOW}acm${NC}\n"