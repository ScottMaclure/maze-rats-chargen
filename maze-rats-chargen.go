package main

// TODO Test suite
// TODO Render to browser - see https://golang.org/doc/articles/wiki/#tmp_6

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type MazeData struct {
	Features        []string
	Paths           []string // TODO Should this be a struct with the description?
	Items           []string
	CombatGear      []string // TODO What about list of actual items inside categories? random(random?
	Appearances     []string
	PhysicalDetails []string
	Backgrounds     []string
	Clothing        []string
	Personalities   []string
	Mannerisms      []string
	Characters      struct {
		MaleNames          []string
		FemaleNames        []string
		UpperClassSurnames []string
		LowerClassSurnames []string
	}
}

type MazeChar struct {
	Strength       int
	Dexterity      int
	Will           int
	Health         int
	ArmorRating    int
	AttackBonus    int
	SpellSlotsMax  int      // number of spell slots TODO needed when Feature=SS recorded? starting characters only
	SpellSlots     []string // actual generated spells
	Feature        string   // what feature you got, for obviousness
	Path           string   // can be empty
	Items          []string
	CombatGear     []string
	Appearance     string
	PhysicalDetail string
	Background     string
	Clothing       string
	Personality    string
	Mannerism      string
	FirstName      string
	Surname        string
	Level          int
	XP             int
}

// why, golang?
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func randInt(min, max int) int {
	return (min + rand.Intn((max+1)-min))
}

func prettyPrint(data interface{}) string {
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return ""
	}
	return string(b[:]) // fuck sake, golang
}

func loadData(filePath string) (MazeData, error) {
	var data MazeData
	dataFile, err := os.Open(filePath)
	defer dataFile.Close()
	if err != nil {
		return data, err
	}
	jsonParser := json.NewDecoder(dataFile)
	err = jsonParser.Decode(&data)
	return data, err
}

func randAbility() int {
	switch randInt(1, 6) {
	case 1, 2:
		return 0
	case 3, 4, 5:
		return 1
	case 6:
		return 2
	default:
		log.Fatal("Unrecognised ability score.")
	}
	return 0
}

// TODO Pass array of exclusions to avoid dupes. See randomItems() and contains()
func getRandom(options []string) string {
	return options[rand.Intn(len(options))]
}

func randomItems(options []string, max int) []string {

	var items = make([]string, 0)
	var current = 0
	var item string

	// TODO limit iterations by reducing options instead of checking for dupes.
	for current < max {
		item = getRandom(options)
		if contains(items, item) == false {
			items = append(items, item)
			current = current + 1
		}
	}

	return items
}

func generateStartingFeature(data MazeData, char *MazeChar) {
	char.Feature = getRandom(data.Features)
	switch char.Feature {
	case "AB":
		char.AttackBonus = char.AttackBonus + 1
	case "SS":
		char.SpellSlotsMax = char.SpellSlotsMax + 1
	case "PA":
		char.Path = getRandom(data.Paths)
	default:
		log.Fatal("Unknown feature:", char.Feature)

	}
}

func generateCombatGear(data MazeData, char *MazeChar) {

	// Everyone gets light armor and a shield
	char.CombatGear = append(char.CombatGear,
		"Light Armor (+1 armor)",
		"Shield (+1 armor, 1 hand)",
		getRandom(data.CombatGear),
		getRandom(data.CombatGear),
	)

}

func generateCharacter(data MazeData) MazeChar {
	var char MazeChar

	char.Level = 1
	char.XP = 0
	char.Strength = randAbility()
	char.Dexterity = randAbility()
	char.Will = randAbility()
	char.Health = 4
	char.ArmorRating = 6 // base human, may increase with armor
	char.AttackBonus = 0
	char.SpellSlotsMax = 0

	generateStartingFeature(data, &char)

	char.Items = randomItems(data.Items, 2) // TODO 6, when more data

	generateCombatGear(data, &char)

	char.Appearance = getRandom(data.Appearances)
	char.PhysicalDetail = getRandom(data.PhysicalDetails)
	char.Background = getRandom(data.Backgrounds)
	char.Clothing = getRandom(data.Clothing)
	char.Personality = getRandom(data.Personalities)
	char.Mannerism = getRandom(data.Mannerisms)

	// TODO Should we determine gender?
	char.FirstName = getRandom(append(data.Characters.MaleNames, data.Characters.FemaleNames...))
	char.Surname = getRandom(append(data.Characters.UpperClassSurnames, data.Characters.LowerClassSurnames...))

	return char
}

func renderCharacterAsHtml(w http.ResponseWriter, char MazeChar) {

	// TODO do this once at start when rendering web pages
	t, err := template.ParseFiles("./chargen.html")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(os.Stdout, char)
	if err != nil {
		log.Fatal(err)
	}
}

func renderCharacterAsJson(w http.ResponseWriter, char MazeChar) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, prettyPrint(char))
}

func handler(w http.ResponseWriter, r *http.Request) {

	var data MazeData
	var char MazeChar
	var err error

	data, err = loadData("./data.json")
	if err != nil {
		log.Fatal(err)
	}
	//	log.Println("data:", prettyPrint(data))

	char = generateCharacter(data)
	log.Println("Generated " + char.FirstName + " " + char.Surname + ", " + char.Background)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	renderCharacterAsJson(w, char)
}

func main() {

	rand.Seed(time.Now().Unix()) // pseudorandom

	//	renderCharacterAsHtml(char)

	http.HandleFunc("/", handler)
	http.ListenAndServe(":5000", nil)
}
