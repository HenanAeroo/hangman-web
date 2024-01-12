package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// Fonction pour effacer le terminal
func clear() {
	// Détermine le système d'exploitation en cours d'utilisation
	osType := runtime.GOOS

	// Commande de nettoyage du terminal en fonction du système d'exploitation
	var cmd *exec.Cmd

	switch osType {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	case "linux", "darwin":
		cmd = exec.Command("clear")
	default:
		// Si le système d'exploitation n'est pas pris en charge, nous ne faisons rien
		fmt.Println("-----------------------------------------")
		fmt.Println("Il semble que le système d'exploitation n'est pas pris en charge ! Nous recommandons vivement" +
			"l'utilisation de Linux / Windows pour avoir un nettoyage du terminal régulier, et ainsi une meilleure lisibilité")
		fmt.Println("-----------------------------------------")
		return
	}

	// Exécute la commande pour effacer le terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Erreur lors du nettoyage du terminal : %v\n", err)
	}
}
