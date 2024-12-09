package shared

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var envVars map[string]string

func LoadEnv(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier : %v", err)
	}
	defer file.Close()

	envVars = make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Ignorer les lignes vides ou les commentaires
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Découper en clé=valeur
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("ligne mal formatée : %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		envVars[key] = value
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("erreur lors de la lecture du fichier : %v", err)
	}
	return nil
}

func GetEnv(key string) string {
	return envVars[key]
}
