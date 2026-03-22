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

// Структура для хранения конфигурации
type Config struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

const jsonFileName = "awg_configs.json"

// Цвета для красивого вывода в терминал
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

	// Если список пуст, сразу предлагаем добавить конфиг
	if len(configs) == 0 {
		fmt.Println(ColorYellow + "Список конфигураций пуст. Давайте добавим первую." + ColorReset)
		addConfig(&configs)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		
		// Показываем текущий статус, если есть конфиги
		if len(configs) > 0 {
			showAwgStatus()
		}

		fmt.Println(ColorCyan + "\n=== AmneziaWG Manager ===" + ColorReset)
		for i, cfg := range configs {
			fmt.Printf("%s%d%s - %s (%s)\n", ColorGreen, i+1, ColorReset, cfg.Name, cfg.Path)
		}
		fmt.Println(ColorCyan + "-------------------------" + ColorReset)
		fmt.Println(ColorYellow + "[+] - Добавить новый конфиг" + ColorReset)
		fmt.Println(ColorRed + "[-] - Отключить ВСЕ соединения (down)" + ColorReset)
		fmt.Println("[0] - Выход")
		fmt.Print(ColorCyan + "\nВыберите действие: " + ColorReset)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "0":
			fmt.Println("Выход из программы...")
			return
		case "+":
			addConfig(&configs)
		case "-":
			disconnectAll(configs)
			pause(reader)
		default:
			// Проверяем, ввел ли пользователь число
			index, err := strconv.Atoi(input)
			if err == nil && index > 0 && index <= len(configs) {
				connect(configs[index-1])
				pause(reader)
			} else {
				fmt.Println(ColorRed + "Неверный ввод, попробуйте снова." + ColorReset)
				pause(reader)
			}
		}
	}
}

// ... (импорты те же, добавь "path/filepath")

func getJsonPath() string {
	// Получаем путь к домашней директории пользователя (например, /home/gsr)
	home, err := os.UserHomeDir()
	if err != nil {
		return jsonFileName // fallback на текущую директорию, если что-то пошло не так
	}
	
	// Создаем путь ~/.config/awm/
	configDir := filepath.Join(home, ".config", "awm")
	
	// Создаем папку, если её нет
	os.MkdirAll(configDir, 0755)
	
	// Полный путь: /home/user/.config/awm/awg_configs.json
	return filepath.Join(configDir, jsonFileName)
}

// Обновленная загрузка
func loadConfigs() []Config {
	var configs []Config
	data, err := os.ReadFile(getJsonPath()) // Используем фиксированный путь
	if err != nil {
		return configs
	}
	json.Unmarshal(data, &configs)
	return configs
}

// Обновленное сохранение
func saveConfigs(configs []Config) {
	data, _ := json.MarshalIndent(configs, "", "  ")
	os.WriteFile(getJsonPath(), data, 0644) // Используем фиксированный путь
}

// Функция добавления нового конфига
func addConfig(configs *[]Config) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(ColorCyan + "Введите абсолютный путь к конфигу (например, /home/user/awg0.conf): " + ColorReset)
		path, _ := reader.ReadString('\n')
		path = strings.TrimSpace(path)

		fmt.Print(ColorCyan + "Введите название конфига (например, Германия): " + ColorReset)
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		if path != "" && name != "" {
			*configs = append(*configs, Config{Name: name, Path: path})
			saveConfigs(*configs)
			fmt.Println(ColorGreen + "Конфиг успешно добавлен!" + ColorReset)
		} else {
			fmt.Println(ColorRed + "Путь и имя не могут быть пустыми." + ColorReset)
		}

		fmt.Print(ColorYellow + "Хотите добавить еще один конфиг? Введите '+' (или 'Enter' чтобы вернуться в меню): " + ColorReset)
		choice, _ := reader.ReadString('\n')
		if strings.TrimSpace(choice) != "+" {
			break
		}
	}
}

// Запуск awg-quick up
func connect(cfg Config) {
	fmt.Printf(ColorBlue+"\nПодключение к %s...\n"+ColorReset, cfg.Name)
	cmd := exec.Command("sudo", "awg-quick", "up", cfg.Path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(ColorRed+"Ошибка при подключении:", err, ColorReset)
	} else {
		fmt.Println(ColorGreen + "Успешно подключено!" + ColorReset)
	}
	
	// После подключения показываем статус
	fmt.Println(ColorYellow + "\nТекущий статус интерфейсов:" + ColorReset)
	showAwgStatus()
}

// Отключение всех конфигов циклом
func disconnectAll(configs []Config) {
	fmt.Println(ColorBlue + "\nОтключение всех конфигураций..." + ColorReset)
	for _, cfg := range configs {
		fmt.Printf("Отключение %s (%s)...\n", cfg.Name, cfg.Path)
		cmd := exec.Command("sudo", "awg-quick", "down", cfg.Path)
		// Мы не перенаправляем stderr, чтобы не спамить ошибками, если интерфейс и так опущен
		err := cmd.Run()
		if err == nil {
			fmt.Printf(ColorGreen+"%s отключен.\n"+ColorReset, cfg.Name)
		}
	}
	fmt.Println(ColorGreen + "Процесс отключения завершен." + ColorReset)
}

// Выполнение sudo awg show
func showAwgStatus() {
	cmd := exec.Command("sudo", "awg", "show")
	output, err := cmd.CombinedOutput()
	if err != nil || len(string(output)) == 0 {
		fmt.Println(ColorRed + "Активных AmneziaWG интерфейсов не найдено." + ColorReset)
	} else {
		fmt.Println(ColorBlue + string(output) + ColorReset)
	}
}

// Очистка экрана (поддерживает Linux/macOS)
func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// Пауза перед очисткой экрана
func pause(reader *bufio.Reader) {
	fmt.Print(ColorYellow + "\nНажмите Enter для продолжения..." + ColorReset)
	reader.ReadString('\n')
}