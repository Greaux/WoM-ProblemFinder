package main

import (
	"fmt"
	"net"
	"runtime"
	"time"

	"github.com/StackExchange/wmi"
	"golang.org/x/sys/windows/registry"
)

type Win32Product struct {
	Name    string
	Version string
}

type Win32_OperatingSystem struct {
	Caption     string
	Version     string
	BuildNumber string
}

func CheckDotNet() {

	fmt.Println("================================================================")
	fmt.Println("Получение установленных версий .Net")

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\NET Framework Setup\NDP`, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		fmt.Println("Ошибка при открытии реестра:", err)
		return
	}
	defer key.Close()

	subKeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		fmt.Println("Ошибка при чтении подключений реестра:", err)
		return
	}

	for _, subKey := range subKeys {
		if subKey == "Servicing" {
			continue
		}

		versionKey, err := registry.OpenKey(key, subKey, registry.QUERY_VALUE)
		if err != nil {
			fmt.Println("Ошибка при открытии подключения реестра:", err)
			continue
		}

		version, _, err := versionKey.GetStringValue("Version")
		if err != nil {
			fmt.Println("Ошибка при чтении значения реестра:", err)
			versionKey.Close()
			continue
		}

		fmt.Println("Версия .NET:", version)

		versionKey.Close()
	}
}

func ping(ipAddress string) error {
	conn, err := net.DialTimeout("ip4:icmp", ipAddress, time.Second*5)
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}

func CheckPing() {
	fmt.Println("================================================================")
	fmt.Println("Проверка соединения")
	fmt.Print("Проверка доступа к серверу игры: ")
	ipAddress := "51.255.1.114" // Замените на нужный IP-адрес
	err := ping(ipAddress)

	if err != nil {
		fmt.Println("Ошибка при выполнении ping:", err)
	} else {
		fmt.Println("Ping успешно выполнен")
	}

	fmt.Print("Проверка доступа к серверу лаунчера: ")
	ipAddress = "gisgames.ru" // Замените на нужный IP-адрес
	err = ping(ipAddress)

	if err != nil {
		fmt.Println("Ошибка при выполнении ping:", err)
	} else {
		fmt.Println("Ping успешно выполнен")
	}
}

func CheckCppRuntime() {
	fmt.Println("================================================================")
	fmt.Println("Получение установленных версий C++ Runtime")

	var products []Win32Product
	query := "SELECT Name, Version FROM Win32_Product WHERE Name LIKE 'Microsoft Visual C++%'"
	err := wmi.Query(query, &products)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return
	}

	for _, product := range products {
		fmt.Printf("Имя: %s, Версия: %s\n", product.Name, product.Version)
	}
}

func CheckDirectX() {
	fmt.Println("================================================================")
	fmt.Println("Получение установленной версий DirectX")
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\DirectX`, registry.QUERY_VALUE)
	if err != nil {
		fmt.Println("Ошибка при открытии ключа реестра:", err)
		return
	}
	defer k.Close()

	version, _, err := k.GetStringValue("Version")
	if err != nil {
		fmt.Println("Ошибка при получении значения из реестра:", err)
		return
	}

	fmt.Println("Установленная версия DirectX:", version)
}

func GetWindowsVersion() {
	var dst []Win32_OperatingSystem
	query := "SELECT Caption, Version, BuildNumber FROM Win32_OperatingSystem"
	err := wmi.Query(query, &dst)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса WMI:", err)
		return
	}

	if len(dst) > 0 {
		fmt.Println("Операционная система:", dst[0].Caption)
		fmt.Println("Версия:", dst[0].Version)
		fmt.Println("Номер сборки:", dst[0].BuildNumber)
	} else {
		fmt.Println("Не удалось получить информацию о версии Windows")
	}
}
func main() {
	if runtime.GOOS == "windows" {
		GetWindowsVersion()
	} else {
		fmt.Println("Программа расчитана на Windows. Запускайте игру из под Wine")
		return
	}

	CheckPing() // Пинг до серверов

	CheckDotNet() // вывод версий DotNet

	CheckCppRuntime() // вывод версий C++ Runtime

	CheckDirectX() // Вывод версий DirectX

	fmt.Println("Нажмите Enter для выхода")
	fmt.Scanln() // Ожидание нажатия Enter
}
