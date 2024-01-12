package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

// Création d'une structure pour le joueur
type Player struct {
	name  string
	score int
	Lives int
}

// Création d'une structure pour le score
type ScoreData struct {
	PlayerName string `json:"player_name"`
	Score      int    `json:"score"`
}

// Constantes pour les couleurs + la vie du joueur
const (
	Reset     = "\033[0m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Purple    = "\033[35m"
	Cyan      = "\033[36m"
	Gray      = "\033[90m"
	White     = "\033[97m"
	Black     = "\033[30m"
	Orange    = "\033[38;5;208m"
	Pink      = "\033[38;5;200m"
	Gold      = "\033[38;5;214m"
	Bold      = "\033[1m"
	Underline = "\033[4m"
	Reverse   = "\033[7m"
	lives     = 10
)

// Fonction qui cherche si un fichier de score existe déjà
func loadScore(player *Player) {
	file, err := os.Open("score.json")
	if err != nil {
		// Si le fichier n'existe pas encore, ne faites rien
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	scoreData := ScoreData{}
	err = decoder.Decode(&scoreData)
	if err != nil {
		log.Fatal(err)
	}

	player.name = scoreData.PlayerName
	player.score = scoreData.Score
}

// Fonction qui sauvegarde le score dans un fichier json
func saveScore(player Player) {
	scoreData := ScoreData{
		PlayerName: player.name,
		Score:      player.score,
	}

	file, err := os.Create("score.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(scoreData)
	if err != nil {
		log.Fatal(err)
	}
}

// On créé le joueur avec son score et son nombre de vies
func createPlayer(name string) Player {
	player := Player{
		name:  name,
		score: 0,
		Lives: 10,
	}
	return player
}

// Fait rentrer dans le jeu
func getUserInput() string {
	var userInput string
	fmt.Println("")
	_, err := fmt.Scanf("%s \n", &userInput)
	if err != nil {
		return "Error"
	}
	return userInput
}

// Permet de lancer le fichier .txt + de choisir un mot aléatoire
func findRandom() (string, string) {
	f, err := os.Open("words.txt") // Ouvre le fichier .txt
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close() // Fermer le fichier .txt

	scanner := bufio.NewScanner(f) // Scanner to read the .txt file

	rand.Seed(time.Now().UnixNano()) // Génération d'un nombre aléatoire

	n := 0

	var ligneAleatoire string // Variable qui contient le mot aléatoire

	// On choisit le mot aléatoire et on stocke dans la variable
	for scanner.Scan() {
		ligne := scanner.Text()
		n++
		if rand.Intn(n) == 0 {
			ligneAleatoire = ligne
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	randomIndex := rand.Intn(len(ligneAleatoire))
	hintLetter := string(ligneAleatoire[randomIndex])

	return ligneAleatoire, hintLetter
}

// Fonction pour créer un nouveau joueur avec un score de 0 + un fichier json qui sauvegarde le score
func loadOrCreatePlayer() Player {
	player := Player{
		name:  "joueur",
		score: 0,
		Lives: 10,
	}

	file, err := os.Open("score.json")
	if err != nil {
		// Si le fichier n'existe pas encore, retourne le nouveau joueur
		return player
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	scoreData := ScoreData{}
	err = decoder.Decode(&scoreData)
	if err != nil {
		log.Fatal(err)
	}

	player.name = scoreData.PlayerName
	player.score = scoreData.Score

	return player
}

// Fonction pour trouver un élément dans une liste
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// Fonction pour trouver les lettres dans le mot mystère
func getCharaterPositions(char string, word string) []int {
	pos := []int{}
	for idx, ch := range word {
		if char == fmt.Sprintf("%c", ch) {
			pos = append(pos, idx)
		}
	}
	return pos
}

// Fonction pour afficher une lettre du mot mystère avant de commencer le jeu
func displayPartialWordWithHint(word string) (string, string) {

	// Choissisez une lettre aléatoire
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(word))
	displayedWord := strings.Builder{}

	// Afficher le mot mystère avec des _ sauf une lettre indice
	for i := 0; i < len(word); i++ {
		if i == randomIndex {
			displayedWord.WriteRune(rune(word[i]))
		} else {
			displayedWord.WriteString(" _")
		}
	}

	hintLetter := string(word[randomIndex])
	return fmt.Sprintf("[%s ]", displayedWord.String()), hintLetter
}

// Fonction qui permet de remplacer les _ par les lettres trouvées
func guessedResult(slice []string, word string) []string {
	strlen := utf8.RuneCountInString(word) // Longueur du mot mystère
	guessed := make([]string, strlen)      // Tableau qui contient les lettres trouvées

	// Remplir le tableau avec des _
	for idx := range guessed {
		guessed[idx] = "_"
	}

	//Remplacez les _ par les lettres trouvées
	for i := 0; i < len(slice); i++ {
		postions := getCharaterPositions(slice[i], word)
		for j := 0; j < len(postions); j++ {
			guessed[postions[j]] = slice[i]
		}
	}
	return guessed
}

// Fonction qui permet de compter le nombre de lettres uniques dans le mot mystère
func countUnique(word string) int {
	characters := []string{}
	for _, ch := range word {
		s := fmt.Sprintf("%c", ch)
		k, _ := Find(characters, s)
		if k == -1 {
			characters = append(characters, s)
		} else {
			continue
		}
	}
	return len(characters)
}

// Fonction qui permet de vérifier si l'input est une lettre
func checkChar(usedChars string) bool {
	if usedChars >= "a" && usedChars <= "z" {
		return true
	}
	return false
}

// Fonction qui ouvre le hangman.txt et permet de choisir la position du pendu selon le nombre de vies
func manPositions(filename string, startPos, endPos int) ([]string, error) {
	f, err := os.Open("hangman.txt") // Ouverture du fichier .txt
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f) // Scanner pour lire le fichier .txt

	var lines []string

	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		if lineNumber >= startPos && lineNumber <= endPos {
			lines = append(lines, scanner.Text())
		}
		if lineNumber > endPos {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if startPos > lineNumber || endPos > lineNumber {
		return nil, fmt.Errorf("Invalid start or end line number")
	}

	return lines, nil
}

// Fonction pour afficher le pendu selon la vie
func txtPosition(lives int) {

	if lives == 10 { // Pour 10 vies, ne rien afficher
		filename := "hangman.txt"
		startPos := 0
		endPos := 0

		lineNumber, err := manPositions(filename, startPos, endPos) // Appel de la fonction manPositions
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, lignes := range lineNumber {
			fmt.Println(lignes)
		}
		return
	}

	// On commence à afficher le pendu
	if lives == 9 {
		filename := "hangman.txt"
		startPos := 1
		endPos := 8

		lineNumber, err := manPositions(filename, startPos, endPos)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, lignes := range lineNumber {
			fmt.Println(lignes)
		}
		return
	}

	if lives == 8 {
		filename := "hangman.txt"
		startPos := 8
		endPos := 15

		lineNumber, err := manPositions(filename, startPos, endPos)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, lignes := range lineNumber {
			fmt.Println(lignes)
		}
		return
	}

	if lives == 7 {
		filename := "hangman.txt"
		startPos := 16
		endPos := 24

		lineNumber, err := manPositions(filename, startPos, endPos)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, lignes := range lineNumber {
			fmt.Println(lignes)
		}
		return
	}

	if lives == 6 {
		filename := "hangman.txt"
		startPos := 25
		endPos := 33

		lineNumber, err := manPositions(filename, startPos, endPos)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, lignes := range lineNumber {
			fmt.Println(lignes)
		}
		return
	}

	if lives == 5 {
		filename := "hangman.txt"
		startPos := 33
		endPos := 41

		lineNumber, err := manPositions(filename, startPos, endPos)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, lignes := range lineNumber {
			fmt.Println(lignes)
		}
		return
	}

	if lives == 4 {
		filename := "hangman.txt"
		startPos := 41
		endPos := 49

		lineNumber, err := manPositions(filename, startPos, endPos)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, lignes := range lineNumber {
			fmt.Println(lignes)
		}
		return
	}

	if lives == 3 {
		filename := "hangman.txt"
		startPos := 49
		endPos := 57

		lineNumber, err := manPositions(filename, startPos, endPos)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, lignes := range lineNumber {
			fmt.Println(lignes)
		}
		return
	}

	if lives == 2 {
		filename := "hangman.txt"
		startPos := 57
		endPos := 65

		lineNumber, err := manPositions(filename, startPos, endPos)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, lignes := range lineNumber {
			fmt.Println(lignes)
		}
		return
	}

	if lives == 1 {
		filename := "hangman.txt"
		startPos := 65
		endPos := 73

		lineNumber, err := manPositions(filename, startPos, endPos)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, lignes := range lineNumber {
			fmt.Println(lignes)
		}
		return
	}

	if lives == 0 {
		filename := "hangman.txt"
		startPos := 73
		endPos := 81

		lineNumber, err := manPositions(filename, startPos, endPos)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, lignes := range lineNumber {
			fmt.Println(lignes)
		}
		return
	}
}

func initHangman(p1 Player) {
	p1.Lives = 10
	usedChars := []string{}      // Tableau qui contient les lettres utilisées
	correctGuess := []string{}   // Tableau qui contient les lettres trouvées
	incorrectGuess := []string{} // Tableau qui contient les lettres incorrectes

	clear()
	fmt.Println("Welcome to the Hangman game! Press 's' to start.")
	data := getUserInput()
	for data != "s" {
		clear()
		fmt.Println("To start, press 's'.")
		data = getUserInput()
	}
	if data == "s" {
		clear()
		fmt.Println("Let the game begin!")
	}

	guess, hintLetter := findRandom()
	uc := countUnique(guess)

	// Afficher une lettre du mot mystère avant de commencer le jeu
	displayedPartialWord, hintLetter := displayPartialWordWithHint(guess)
	fmt.Printf("Word Hint: %s (Hint Letter: %s)\n", displayedPartialWord, hintLetter)

	// Début du jeu
	fmt.Println("Find the mystery word!")
	for {
		txtPosition(lives - len(incorrectGuess))
		fmt.Println(guessedResult(correctGuess, guess))
		fmt.Println("Choose a letter: ")
		value := getUserInput()
		if checkChar(value) == false {
			fmt.Println("Please enter a valid letter.")
			continue
		}
		k, _ := Find(usedChars, value)
		if k == -1 {
			fmt.Println("")
			usedChars = append(usedChars, value)
			fmt.Println("")
		} else {
			fmt.Println("You've already used this letter. Try another one.")
			continue
		}
		if strings.Contains(guess, value) {
			correctGuess = append(correctGuess, value)
			fmt.Printf(" - Yes! '%s' is part of the word.\n", value)
		} else {
			incorrectGuess = append(incorrectGuess, value)
			fmt.Printf(" - No, '%s' is not in the word.\n", value)
		}
		fmt.Println("")
		fmt.Println("")
		fmt.Printf("Remaining Lives: %d\n", p1.Lives-len(incorrectGuess))
		fmt.Print("Incorrect Letters: ")
		fmt.Println(incorrectGuess)
		fmt.Print("Correct Letters: ")
		fmt.Println(correctGuess)
		if len(incorrectGuess) == p1.Lives {
			break
		}
		if len(correctGuess) == uc {
			break
		}
	}
	if len(incorrectGuess) == p1.Lives {
		txtPosition(lives - len(incorrectGuess))
		fmt.Println("")
		fmt.Println("----------------Game Over----------------")
		fmt.Printf("The word was '%s'\n", guess)
		fmt.Println("")
		fmt.Printf("Your score is %d", p1.score)
	}
	if len(correctGuess) == uc {
		fmt.Println("")
		fmt.Println("----------------Congratulations!----------------")
		fmt.Printf("The word was '%s'\n", guess)
		fmt.Println("")
		fmt.Printf("Your score was: %d\n", p1.score)
		p1.score += 10
		fmt.Printf("Your new score is: %d", p1.score)
	}

	fmt.Println("\nDo you want to play again ? Type 'yes' to replay or 'no' to quit.")
	dataRe := getUserInput()
	for dataRe != "yes" && dataRe != "no" {
		clear()
		fmt.Println("Please enter 'yes' to replay or 'no' to quit.")
		dataRe = getUserInput()
	}

	if dataRe == "yes" {
		saveScore(p1)
		initHangman(p1)
	} else {
		saveScore(p1)
		return
	}
}

func main() {
	p1 := loadOrCreatePlayer()
	initHangman(p1)
}
