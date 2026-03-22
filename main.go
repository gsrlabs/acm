package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

const jsonFileName = "awg_configs.json"

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
)

func main() {
	configs := loadConfigs()

	if len(configs) == 0 {
		fmt.Println(ColorYellow + "The list of configurations is empty. Let's add the first one." + ColorReset)
		addConfig(&configs)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		
		if len(configs) > 0 {
			showAwgStatus()
		}

		fmt.Println(ColorCyan + "\n=== AmneziaWG CLI Manager ===" + ColorReset)
		for i, cfg := range configs {
			fmt.Printf("%s%d%s - %s (%s)\n", ColorGreen, i+1, ColorReset, cfg.Name, cfg.Path)
		}
		fmt.Println(ColorCyan + "-------------------------" + ColorReset)
		fmt.Println(ColorYellow + "[+] - Add a new config" + ColorReset)
		fmt.Println(ColorRed + "[-] - Disable connections (down)" + ColorReset)
		fmt.Println("[0] - Exit")
		fmt.Print(ColorCyan + "\nSelect an action: " + ColorReset)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "0":
			fmt.Println("Exiting the program...")
			return
		case "+":
			addConfig(&configs)
		case "-":
			disconnectAll(configs)
			pause(reader)
		default:
			index, err := strconv.Atoi(input)
			if err == nil && index > 0 && index <= len(configs) {
				connect(configs[index-1])
				pause(reader)
			} else {
				fmt.Println(ColorRed + "Incorrect input, please try again." + ColorReset)
				pause(reader)
			}
		}
	}
}

func getJsonPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return jsonFileName
	}
	
	configDir := filepath.Join(home, ".config", "acm")
	
	os.MkdirAll(configDir, 0755)
	
	return filepath.Join(configDir, jsonFileName)
}

func loadConfigs() []Config {
	var configs []Config
	data, err := os.ReadFile(getJsonPath())
	if err != nil {
		return configs
	}
	json.Unmarshal(data, &configs)
	return configs
}

func saveConfigs(configs []Config) {
	data, _ := json.MarshalIndent(configs, "", "  ")
	os.WriteFile(getJsonPath(), data, 0644)
}

func addConfig(configs *[]Config) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(ColorCyan + "Enter the absolute path to the configuration (example, /home/user/amnezia_for_awg.conf): " + ColorReset)
		path, _ := reader.ReadString('\n')
		path = strings.TrimSpace(path)

		fmt.Print(ColorCyan + "Enter the absolute path to the configuration (example, Germany): " + ColorReset)
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		if path != "" && name != "" {
			*configs = append(*configs, Config{Name: name, Path: path})
			saveConfigs(*configs)
			fmt.Println(ColorGreen + "The config has been successfully added!" + ColorReset)
		} else {
			fmt.Println(ColorRed + "The path and name cannot be empty." + ColorReset)
		}

		fmt.Print(ColorYellow + "Do you want to add another config? Enter '+' (or 'Enter' to return to the menu): " + ColorReset)
		choice, _ := reader.ReadString('\n')
		if strings.TrimSpace(choice) != "+" {
			break
		}
	}
}

func connect(cfg Config) {
	fmt.Printf(ColorBlue+"\nConnection to %s...\n"+ColorReset, cfg.Name)
	cmd := exec.Command("sudo", "awg-quick", "up", cfg.Path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(ColorRed+"Connection error:", err, ColorReset)
	} else {
		fmt.Println(ColorGreen + "Successfully connected!" + ColorReset)
	}
	
	fmt.Println(ColorYellow + "\nCurrent status of interfaces:" + ColorReset)
	showAwgStatus()
}

func disconnectAll(configs []Config) {
	fmt.Println(ColorBlue + "\nDisabling all configurations..." + ColorReset)
	for _, cfg := range configs {
		fmt.Printf("Disconnection %s (%s)...\n", cfg.Name, cfg.Path)
		cmd := exec.Command("sudo", "awg-quick", "down", cfg.Path)
		err := cmd.Run()
		if err == nil {
			fmt.Printf(ColorGreen+"%s disconnected.\n"+ColorReset, cfg.Name)
		}
	}
	fmt.Println(ColorGreen + "The shutdown process is complete." + ColorReset)
}

func showAwgStatus() {
	cmd := exec.Command("sudo", "awg", "show")
	output, err := cmd.CombinedOutput()
	if err != nil || len(string(output)) == 0 {
		fmt.Println(ColorRed + "No Active AmneziaWG interfaces found." + ColorReset)
	} else {
		fmt.Println(ColorBlue + string(output) + ColorReset)
	}
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func pause(reader *bufio.Reader) {
	fmt.Print(ColorYellow + "\nPress Enter to continue..." + ColorReset)
	reader.ReadString('\n')
}