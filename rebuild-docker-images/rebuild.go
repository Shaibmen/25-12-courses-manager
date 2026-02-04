package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Конфигурация: папки и соответствующие имена контейнеров
var config = map[string]string{
	"../document-service":     "metapohuism/document-service",
	"../auth-service":         "shaibmen/auth-service",
	"../online-courses":       "shaibmen/online-courses",
	"../online-courses-front": "shaibmen/online-courses-front:local",
}

func main() {
	// Получаем текущую директорию
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Ошибка получения текущей директории: %v", err)
	}

	// Проверяем наличие Docker
	if !isDockerAvailable() {
		log.Fatal("Docker не найден. Убедитесь, что Docker установлен и доступен в PATH")
	}

	// Собираем контейнеры для каждой папки
	for folder, containerName := range config {
		buildDockerImage(baseDir, folder, containerName)
	}

	fmt.Println("Все контейнеры успешно собраны!")
}

// buildDockerImage собирает Docker-образ из указанной папки
func buildDockerImage(baseDir, folder, containerName string) {
	// Формируем полный путь к папке
	folderPath := filepath.Join(baseDir, folder)

	// Проверяем существование папки
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		log.Printf("Папка %s не найдена, пропускаем...", folderPath)
		return
	}

	// Проверяем наличие Dockerfile в папке
	dockerfilePath := filepath.Join(folderPath, "Dockerfile")
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		log.Printf("Dockerfile не найден в %s, пропускаем...", folderPath)
		return
	}

	fmt.Printf("Сборка контейнера '%s' из папки '%s'...\n", containerName, folder)

	// Выполняем команду docker build
	cmd := exec.Command("docker", "build", "-t", containerName, folderPath)

	// Настраиваем вывод команд
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Выполняем команду
	err := cmd.Run()
	if err != nil {
		log.Printf("Ошибка при сборке контейнера '%s': %v", containerName, err)
		return
	}

	fmt.Printf("Контейнер '%s' успешно собран!\n", containerName)
}

// isDockerAvailable проверяет доступность Docker
func isDockerAvailable() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}
