package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Type struct {
	// Name of the type
	Name string `json:"name"`
	// The effective types, damage multiplize 2x
	EffectiveAgainst []string `json:"effectiveAgainst"`
	// The weak types that against, damage multiplize 0.5x
	WeakAgainst []string `json:"weakAgainst"`
}

type Pokemon struct {
	Number         string   `json:"Number"`
	Name           string   `json:"Name"`
	Classification string   `json:"Classification"`
	TypeI          []string `json:"Type I"`
	TypeII         []string `json:"Type II,omitempty"`
	Weaknesses     []string `json:"Weaknesses"`
	FastAttackS    []string `json:"Fast Attack(s)"`
	Weight         string   `json:"Weight"`
	Height         string   `json:"Height"`
	Candy          struct {
		Name     string `json:"Name"`
		FamilyID int    `json:"FamilyID"`
	} `json:"Candy"`
	NextEvolutionRequirements struct {
		Amount int    `json:"Amount"`
		Family int    `json:"Family"`
		Name   string `json:"Name"`
	} `json:"Next Evolution Requirements,omitempty"`
	NextEvolutions []struct {
		Number string `json:"Number"`
		Name   string `json:"Name"`
	} `json:"Next evolution(s),omitempty"`
	PreviousEvolutions []struct {
		Number string `json:"Number"`
		Name   string `json:"Name"`
	} `json:"Previous evolution(s),omitempty"`
	SpecialAttacks      []string `json:"Special Attack(s)"`
	BaseAttack          int      `json:"BaseAttack"`
	BaseDefense         int      `json:"BaseDefense"`
	BaseStamina         int      `json:"BaseStamina"`
	CaptureRate         float64  `json:"CaptureRate"`
	FleeRate            float64  `json:"FleeRate"`
	BuddyDistanceNeeded int      `json:"BuddyDistanceNeeded"`
}

// Move is an attack information. The
type Move struct {
	// The ID of the move
	ID int `json:"id"`
	// Name of the attack
	Name string `json:"name"`
	// Type of attack
	Type string `json:"type"`
	// The damage that enemy will take
	Damage int `json:"damage"`
	// Energy requirement of the attack
	Energy int `json:"energy"`
	// Dps is Damage Per Second
	Dps float64 `json:"dps"`
	// The duration
	Duration int `json:"duration"`
}

// BaseData is a struct for reading data.json
type BaseData struct {
	Types    []Type    `json:"types"`
	Pokemons []Pokemon `json:"pokemons"`
	Moves    []Move    `json:"moves"`
}

var baseData BaseData

func (pokemon Pokemon) String() (str string) {
	str += pokemon.Name + ":" +
		"\n\tNumber: " + pokemon.Number +
		"\n\tType I: " + strings.Join(pokemon.TypeI, ", ")
	if len(pokemon.TypeII) > 0 {
		str += "\n\tType II: " + strings.Join(pokemon.TypeII, ", ")
	}
	str +=
		"\n\tWeight: " + pokemon.Weight +
			"\n\tHeight: " + pokemon.Height +
			"\n\tBase Attack: " + strconv.Itoa(pokemon.BaseAttack) +
			"\n\tBase Defense: " + strconv.Itoa(pokemon.BaseDefense) +
			"\n\tBase Stamina: " + strconv.Itoa(pokemon.BaseStamina)
	str += "\n\tFast Attack(s):"
	for _, e := range pokemon.FastAttackS {
		str += "\n\t\t" + e
	}
	if len(pokemon.PreviousEvolutions) > 0 {
		str += "\n\tPrevious Evolution(s):"
		for _, e := range pokemon.PreviousEvolutions {
			str += "\n\t\t" + e.Name
		}
	}
	if len(pokemon.NextEvolutions) > 0 {
		str += "\n\tNext Evolution(s):"
		for _, e := range pokemon.NextEvolutions {
			str += "\n\t\t" + e.Name
		}
		str += "\n\tNext Evolution Requirements:"
		str += "\n\t\tAmount: " + strconv.Itoa(pokemon.NextEvolutionRequirements.Amount)
		str += "\n\t\tName: " + pokemon.NextEvolutionRequirements.Name
	}
	str += "\n\tSpecial Attack(s):"
	for _, e := range pokemon.SpecialAttacks {
		str += "\n\t\t" + e
	}

	return
}

func printPokemons(pokemons []Pokemon) (str string) {
	for _, pokemon := range pokemons {
		str += pokemon.String() + "\n"
	}
	return
}

func (typ Type) String() (str string) {
	str += "Pokemon Type " + typ.Name + ":"
	str += "\nEffective Against:"
	for _, e := range typ.EffectiveAgainst {
		str += "\n- " + e
	}
	str += "\nWeak Against:"
	for _, e := range typ.WeakAgainst {
		str += "\n- " + e
	}
	str += "\nExample Pokemons: "
	pokemons := filterPokemons(baseData.Pokemons, "type", typ.Name)
	for _, e := range pokemons {
		str += "\n- " + e.Name
	}
	str += "\n"
	return
}

func printTypes(types []Type) (str string) {
	for _, typ := range types {
		str += typ.String() + "\n"
	}
	return
}

func getTypeString(typeName string, types []Type) (result string) {
	for _, typ := range types {
		if strings.ToLower(typeName) == strings.ToLower(typ.Name) {
			return typ.String() + "\n"
		}
	}
	return ""
}

func (move Move) String() (str string) {
	var moveType string

	pokemons := filterPokemons(baseData.Pokemons, "fastattack", move.Name)
	if len(pokemons) > 0 {
		moveType = "Fast Attack"
	} else {
		moveType = "Special Attack"
		pokemons = filterPokemons(baseData.Pokemons, "specialattack", move.Name)
	}

	str += "Pokemon Move " + moveType + " " + move.Name + ":" +
		"\nNumber: " + strconv.Itoa(move.ID) +
		"\nDamage: " + strconv.Itoa(move.Damage) +
		"\nEnergy: " + strconv.Itoa(move.Energy) +
		"\nDps: " + strconv.FormatFloat(move.Dps, 'f', 2, 64) +
		"\nDuration: " + strconv.Itoa(move.Duration) +
		"\nPokemons with this move: "
	for _, pokemon := range pokemons {
		str += "\n- " + pokemon.Name
	}
	str += "\n"
	return
}

func printMoves(moves []Move) (str string) {
	for _, move := range moves {
		str += move.String() + "\n"
	}
	return
}

func getMoveString(moveName string, moves []Move) (result string) {
	for _, move := range moves {
		if strings.ToLower(moveName) == strings.ToLower(move.Name) {
			return move.String() + "\n"
		}
	}
	return ""
}

func includesString(ss []string, s string) bool {
	for _, str := range ss {
		if strings.ToLower(str) == strings.ToLower(s) {
			return true
		}
	}
	return false
}

func filterPokemons(pokemons []Pokemon, filterType string, filter string) (filteredPokemons []Pokemon) {
	switch filterType {
	case "type":
		for _, pokemon := range pokemons {
			if includesString(pokemon.TypeI, filter) || includesString(pokemon.TypeII, filter) {
				filteredPokemons = append(filteredPokemons, pokemon)
			}
		}
	case "name":
		for _, pokemon := range pokemons {
			if strings.ToLower(pokemon.Name) == strings.ToLower(filter) {
				filteredPokemons = append(filteredPokemons, pokemon)
			}
		}
	case "fastattack":
		for _, pokemon := range pokemons {
			if includesString(pokemon.FastAttackS, filter) {
				filteredPokemons = append(filteredPokemons, pokemon)
			}
		}
	case "specialattack":
		for _, pokemon := range pokemons {
			if includesString(pokemon.SpecialAttacks, filter) {
				filteredPokemons = append(filteredPokemons, pokemon)
			}
		}
	}
	return
}

func reversePokemons(pokemons []Pokemon) (reversedPokemons []Pokemon) {
	reversedPokemons = pokemons
	for i := len(reversedPokemons)/2 - 1; i >= 0; i-- {
		opp := len(reversedPokemons) - 1 - i
		reversedPokemons[i], reversedPokemons[opp] = reversedPokemons[opp], reversedPokemons[i]
	}
	return
}

func sortPokemons(pokemons []Pokemon, sortby string) (sortedPokemons []Pokemon) {
	switch strings.ToLower(sortby) {
	case "number", "id":
		sort.Slice(pokemons, func(i, j int) bool {
			return pokemons[i].Number < pokemons[j].Number
		})
		log.Println("Sorted pokemons by number.")
	case "name":
		sort.Slice(pokemons, func(i, j int) bool {
			return pokemons[i].Name < pokemons[j].Name
		})
		log.Println("Sorted pokemons by name.")
	case "weight":
		sort.Slice(pokemons, func(i, j int) bool {
			num1Str := strings.Replace(strings.Replace(pokemons[i].Weight, ",", ".", 1), " kg", "", 1)
			num2Str := strings.Replace(strings.Replace(pokemons[j].Weight, ",", ".", 1), " kg", "", 1)
			num1, _ := strconv.ParseFloat(num1Str, 64)
			num2, _ := strconv.ParseFloat(num2Str, 64)
			return num1 < num2
		})
		log.Println("Sorted pokemons by weight.")
	case "height":
		sort.Slice(pokemons, func(i, j int) bool {
			num1Str := strings.Replace(strings.Replace(pokemons[i].Height, ",", ".", 1), " m", "", 1)
			num2Str := strings.Replace(strings.Replace(pokemons[j].Height, ",", ".", 1), " m", "", 1)
			num1, _ := strconv.ParseFloat(num1Str, 64)
			num2, _ := strconv.ParseFloat(num2Str, 64)
			return num1 < num2
		})
		log.Println("Sorted pokemons by height.")
	case "baseattack":
		sort.Slice(pokemons, func(i, j int) bool {
			return pokemons[i].BaseAttack < pokemons[j].BaseAttack
		})
		log.Println("Sorted pokemons by base attack.")
	case "basedefence", "basedefense":
		sort.Slice(pokemons, func(i, j int) bool {
			return pokemons[i].BaseDefense < pokemons[j].BaseDefense
		})
		log.Println("Sorted pokemons by base defense.")
	case "basestamina":
		sort.Slice(pokemons, func(i, j int) bool {
			return pokemons[i].BaseStamina < pokemons[j].BaseStamina
		})
		log.Println("Sorted pokemons by base stamina.")
	default:
		return nil
	}
	sortedPokemons = pokemons
	return
}

func listPokemonsByTypeHandler(w http.ResponseWriter, r *http.Request, pokemonType string) {
	pokemons := filterPokemons(baseData.Pokemons, "type", pokemonType)
	if len(pokemons) < 1 {
		msg := fmt.Sprintf("Could not find any pokemons for the type %s\n", pokemonType)
		http.Error(w, msg, http.StatusNotFound)
		log.Printf("Client provided invalid pokemon type %s\n", pokemonType)
		return
	}
	keys, ok := r.URL.Query()["sortby"]
	if ok {
		if len(keys[0]) > 0 {
			pokemons = sortPokemons(pokemons, keys[0])
			if pokemons == nil {
				msg := fmt.Sprintf("Wrong sorting type %s! Use one of the following: number, weight, height, baseattack, basedefense, basestamina\n", keys[0])
				http.Error(w, msg, http.StatusNotFound)
				log.Printf("Client provided wrong sorting type %s\n", keys[0])
				return
			}
			_, ok = r.URL.Query()["reversed"]
			if ok {
				pokemons = reversePokemons(pokemons)
				log.Println("Reversed sorted pokemons.")
			}
		} else {
			fmt.Fprint(w, "Available sorting types: number, weight, height, baseattack, basedefence, basestamina\n")
			log.Println("Client requested sorting types.")
			return
		}
	}
	result := printPokemons(pokemons)
	fmt.Fprint(w, result)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["type"]
	if ok {
		if len(keys[0]) > 0 {
			log.Printf("Client requested to list pokemons by type %s\n", keys[0])
			listPokemonsByTypeHandler(w, r, keys[0])
		} else {
			http.Error(w, "You need to provide the pokemon type value!", 400)
			log.Println("400 Error: Client left Pokemon type value empty.")
		}
		return
	}

	_, ok = r.URL.Query()["pokemons"]
	if ok {
		result := printPokemons(baseData.Pokemons)
		fmt.Fprintf(w, result)
		log.Println("Served all pokemons.")
		return
	}

	_, ok = r.URL.Query()["types"]
	if ok {
		result := printTypes(baseData.Types)
		fmt.Fprintf(w, result)
		log.Println("Served all pokemon types.")
		return
	}

	_, ok = r.URL.Query()["moves"]
	if ok {
		result := printMoves(baseData.Moves)
		fmt.Fprint(w, result)
		log.Println("Served all pokemon moves.")
		return
	}

	if r.URL.Path == "/list" && !strings.Contains(r.RequestURI, "?") {
		result := "-----Pokemons-----\n"
		result += printPokemons(baseData.Pokemons)
		result += "-----Types-----\n"
		result += printTypes(baseData.Types)
		result += "-----Moves-----\n"
		result += printMoves(baseData.Moves)
		fmt.Fprint(w, result)
		log.Println("Served all pokemons, moves and types.")
		return
	}

	log.Println("404 Error: Client's request is not found.")
	http.Error(w, "404 Not Found", http.StatusNotFound)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.Replace(r.URL.Path, "/", "", 1)

	if len(query) < 1 || query == "help" {
		str := "-----PokéDex API Help-----\n\n"
		str += "Display this help message: / or /help\n\n"
		str += "---Listing Pokemons, Types and Moves---\n\n"
		str += "List all pokemons, moves and types: /list\n"
		str += "List all pokemons: /list?pokemons\n"
		str += "List all types: /list?types\n"
		str += "List all moves: /list?moves\n"
		str += "List all pokemons for a given type: /list?type={type}\n"
		str += "Get valid attributes to sort Pokemons by: /list?type={type}&sortby\n"
		str += "List all pokemons for a given type and sort them by an attribute: /list?type={type}&sortby={attribute}\n"
		str += "List all pokemons for a given type and sort them by an attribute in reversed order: /list?type={type}&sortby={attribute}&reversed\n"
		str += "\n---Getting information about a Pokemon, Type or Move---\n\n"
		str += "/{resourceName}\n"
		fmt.Fprint(w, str)
		log.Println("Served help text.")
	} else {
		query = strings.ToLower(query)

		log.Printf("Client requested resource %s\n", query)

		pokemon := filterPokemons(baseData.Pokemons, "name", query)
		if len(pokemon) > 0 {
			result := printPokemons(pokemon)
			fmt.Fprint(w, result)
			log.Printf("Served pokemon with the name %s\n", query)
			return
		}

		typ := getTypeString(query, baseData.Types)
		if len(typ) > 0 {
			result := typ
			fmt.Fprintf(w, result)
			log.Printf("Served type with the name %s\n", query)
			return
		}

		move := getMoveString(query, baseData.Moves)
		if len(move) > 0 {
			result := move
			fmt.Fprintf(w, result)
			log.Printf("Served move with the name %s\n", query)
			return
		}

		log.Println("404 Error: Client's request is not found.")
		http.Error(w, "404 Not Found", http.StatusNotFound)
	}
}

func readJSON() {
	jsonFile, _ := os.Open("data.json")

	fmt.Println("Successfully read data.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &baseData)
}

func main() {
	readJSON()

	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/", getHandler)
	log.Println("starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
