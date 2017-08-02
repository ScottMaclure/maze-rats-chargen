package main

// TODO Test suite
// TODO Render to browser - see https://golang.org/doc/articles/wiki/#tmp_6

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/inflection"
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
	Magic struct {
		PhysicalEffects  []string
		PhysicalElements []string
		PhysicalForms    []string
		EtherialEffects  []string
		EtherialElements []string
		EtherialForms    []string
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

// globals
var data MazeData
var char MazeChar
var tpl *template.Template

// Used for generating MazeData.Magic
var spellTable [6][2][2]string

// why, golang?
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

var vowels = []string{"a", "e", "i", "o", "u"}

/**
 * Given a string input, return it formatted for a sentence.
 * Doesn't care about capitalization, that's handled by css.
 * E.g. "Battle Scars" -> "Battle Scars" (plural)
 * "Birthmark" -> "a Birthmark" (a)
 * "Animal Scent" -> "an Animal Scent" (an)
 * As more cases as discovered, reveal here.
 */
func naturalLanguageSplice(input string) string {

	input = strings.TrimSpace(input)

	//  plural test, "Battle Scars" means NO a OR an.
	if isPlural(input) {
		return input
	}

	// vowel test for "an apple" etc.
	first := strings.ToLower(input[0:1])
	if contains(vowels, first) {
		return "an " + input
	}

	return "a " + input
}

func isPlural(input string) bool {

	words := strings.Split(input, " ")
	lastWord := strings.ToLower(words[len(words)-1])
	singular := inflection.Singular(lastWord)
	plural := inflection.Plural(lastWord)

	// log.Printf("lastWord=%s, singular=%s, plural=%s", lastWord, singular, plural)

	if singular != plural && lastWord == plural {
		return true
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

func initSpellTable() {

	spellTable[0][0][0] = "PhysicalEffects"
	spellTable[0][0][1] = "PhysicalForms"
	spellTable[0][1][0] = "EtherialElements"
	spellTable[0][1][1] = "PhysicalForms"

	spellTable[1][0][0] = "PhysicalEffects"
	spellTable[1][0][1] = "EtherialForms"
	spellTable[1][1][0] = "EtherialElements"
	spellTable[1][1][1] = "EtherialForms"

	spellTable[2][0][0] = "EtherialEffects"
	spellTable[2][0][1] = "PhysicalForms"
	spellTable[2][1][0] = "PhysicalEffects"
	spellTable[2][1][1] = "PhysicalElements"

	spellTable[3][0][0] = "EtherialEffects"
	spellTable[3][0][1] = "EtherialForms"
	spellTable[3][1][0] = "PhysicalEffects"
	spellTable[3][1][1] = "EtherialElements"

	spellTable[4][0][0] = "PhysicalElements"
	spellTable[4][0][1] = "PhysicalForms"
	spellTable[4][1][0] = "EtherialEffects"
	spellTable[4][1][1] = "PhysicalElements"

	spellTable[5][0][0] = "PhysicalElements"
	spellTable[5][0][1] = "EtherialForms"
	spellTable[5][1][0] = "EtherialEffects"
	spellTable[5][1][1] = "EtherialElements"

}

func initInflections() {
	inflection.AddUncountable("cautious")
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

	log.Println("Generated " + char.FirstName + " " + char.Surname + ", " + char.Background)

	return char
}

// Use the spellTable and json Magic data to generate a random combination spellname.
func generateSpellName() string {

	row := randInt(0, 5)
	col := randInt(0, 1)

	spellType1 := spellTable[row][col][0]
	spellType2 := spellTable[row][col][1]

	//	log.Println("row:", row+1, "col:", col+1, "spellType1:", spellType1, "spellType2:", spellType2)

	v := reflect.ValueOf(data.Magic)
	spellValue1 := v.FieldByName(spellType1).Interface().([]string)
	spellValue2 := v.FieldByName(spellType2).Interface().([]string)

	spellName := getRandom(spellValue1) + " " + getRandom(spellValue2)
	log.Println("spellName (" + spellType1 + " + " + spellType2 + "): " + spellName)

	return spellName
}

func renderCharacterAsHtml(w http.ResponseWriter, char MazeChar) {
	err := tpl.Execute(w, char)
	if err != nil {
		log.Println(err)
	}
}

func writeJsonHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
}

func renderCharacterAsJson(w http.ResponseWriter, char MazeChar) {
	writeJsonHeaders(w)
	fmt.Fprintf(w, prettyPrint(char))
}

func handleJson(w http.ResponseWriter, r *http.Request) {
	char = generateCharacter(data)
	renderCharacterAsJson(w, char)
}

func handleHtml(w http.ResponseWriter, r *http.Request) {
	char = generateCharacter(data)
	renderCharacterAsHtml(w, char)
}

func handleJsonSpell(w http.ResponseWriter, r *http.Request) {
	writeJsonHeaders(w)
	fmt.Fprintf(w, prettyPrint(generateSpellName()))
}

func initData() {
	var err error
	data, err = loadData("./data.json")
	if err != nil {
		log.Fatal(err)
	}
	//	log.Println("data:", prettyPrint(data))
}

func main() {

	log.Println("maze-rats-chargen server starting up.")

	rand.Seed(time.Now().Unix()) // pseudorandom

	// Load shared resources once on start, not every request.

	var err error

	initData()

	file, err := ioutil.ReadFile("./chargen.html")
	if err != nil {
		log.Fatal(err)
	}

	funcs := template.FuncMap{
		"join": strings.Join,
		"nls":  naturalLanguageSplice,
	}
	tpl, err = template.New("chargen").Funcs(funcs).Parse(string(file))
	if err != nil {
		log.Fatal(err)
	}

	initSpellTable()
	initInflections()

	http.HandleFunc("/", handleHtml)
	http.HandleFunc("/json", handleJson) // TODO What about using accept header?
	http.HandleFunc("/spell", handleJsonSpell)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.ListenAndServe(":5000", nil)

}
