package shop

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/brotherschimes/noblespaycash/input"
	"github.com/brotherschimes/noblespaycash/item"
	"github.com/brotherschimes/noblespaycash/selection"
)

const uniqueCustName = "Ulric"

// DoSales sets up a stock and then allows the player to sell the items
func DoSales(provider selection.AnswerProvider) {
	itemTypes, storeStock := Setup()
	TransactionLoop(provider, itemTypes, &storeStock)
}

// Setup prepares a stock of basic items for trying out sales
func Setup() (map[string]item.Type, map[string]int) {
	itemTypes := map[string]item.Type{
		"sword":    {Name: "sword", PluralName: "swords", IsWeapon: true},
		"axe":      {Name: "axe", PluralName: "axes", IsWeapon: true},
		"trailMix": {Name: "bag of trail mix", PluralName: "bags of trail mix", IsWeapon: false},
	}

	storeStock := map[string]int{
		"sword":    1,
		"axe":      1,
		"trailMix": 1,
	}

	return itemTypes, storeStock
}

// TransactionLoop loops across the customer until they are satisfied
func TransactionLoop(provider selection.AnswerProvider, itemTypes map[string]item.Type, storeStock *map[string]int) {
	says(uniqueCustName, "Hi, Bailoe!")

	customerRequests := []string{"I would like to purchase a weapon!", "I would STILL like to purchase a weapon!", "PLEASE sell me a weapon...", "A weapon, please!"}
	customerRequestIndex := 0

	isSatisfied := false

	for !isSatisfied {
		nextString := customerRequests[customerRequestIndex]
		announceItemQty(*storeStock, itemTypes)
		says(uniqueCustName, nextString)
		if customerRequestIndex < len(customerRequests)-1 {
			customerRequestIndex++
		}
		SellWeapons(provider, storeStock, itemTypes, &isSatisfied)
	}
	says(uniqueCustName, "Bye, Mr Celhai.")
	fmt.Println(uniqueCustName + " leaves.")
}

// SellWeapons provides a list of available weapons and asks which one should be sold
func SellWeapons(provider selection.AnswerProvider, stock *map[string]int, types map[string]item.Type, isSatisfied *bool) {
	if (*stock)["sword"] <= 0 && (*stock)["axe"] <= 0 {
		fmt.Println("You inform " + uniqueCustName + " that you have no weapon left for sale.")
		*isSatisfied = true
		return
	}

	SellWeaponsFoo(provider, stock, types, isSatisfied)
}

func sortedKeys(stock *map[string]int) []string {
	keys := make([]string, len(*stock))

	i := 0
	for k := range *stock {
		keys[i] = k
		i++
	}

	sort.Strings(keys)
	return keys
}

// SellWeaponsFoo ...
func SellWeaponsFoo(provider selection.AnswerProvider, stock *map[string]int, types map[string]item.Type, isSatisfied *bool) {
	fmt.Println("What would you like to sell Ulric?")

	map2 := map[string]int{}

	/* Copy Content from Map1 to Map2*/
	for key, value := range *stock {
		if value > 0 {
			map2[key] = value
		}
	}

	i := 1
	keys := sortedKeys(&map2)
	for _, key := range keys {
		fmt.Printf("[%v] %s\n", i, strings.Title(types[key].Name))
		i++
	}
	fmt.Println("[0] Nothing")

	selection := provider.GetSelection(len(keys) + 1)

	if selection == 0 {
		fmt.Println("You decide not to sell anything at this point.")
		return
	}

	SellWeapon(provider, keys[selection-1], stock, types)
}

// SellWeapon tries to sell a single weapon
func SellWeapon(provider selection.YesNoAnswerProvider, weapon string, stock *map[string]int, types map[string]item.Type) {
	weaponName := types[weapon].Name
	weaponPluralName := types[weapon].PluralName

	if (*stock)[weapon] <= 0 {
		fmt.Println("You have no " + weaponPluralName + " left for sale.")
		return
	}

	if !types[weapon].IsWeapon {
		fmt.Printf("Ulric asks how he is meant to kill goblins with a %s? ", weaponName)
		return
	}

	fmt.Println("Would you like to sell Ulric a " + weaponName + "? (y/n)")

	answer := provider.GetAnswer()

	if answer {
		does(uniqueCustName, "happily takes the "+weaponName+".")
		(*stock)[weapon]--
	} else {
		does(uniqueCustName, "is sad you did not sell him the "+weaponName+".")
	}
}

func announceItemQty(stock map[string]int, types map[string]item.Type) {
	for name, qty := range stock {
		fmt.Printf("You have %v %s.\n", qty, types[name].PluralName)
	}
}

func does(name, action string) {
	fmt.Println(name + " " + action)
}

func says(name, speech string) {
	fmt.Println(name + " says: " + "\"" + speech + "\"")
}

func main() {
	reader := input.Reader{Reader: bufio.NewReader(os.Stdin)}
	DoSales(reader)
}
