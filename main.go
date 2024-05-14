package main

import (
	"flag"
	"fmt"

	"github.com/fatih/color"
)

// djangospy -u https://target.com -version True -list-packages True -max-thread 100 -delay 100
func main() {
	logo := `
░▒▓███████▓▒░       ░▒▓█▓▒░░▒▓██████▓▒░░▒▓███████▓▒░ ░▒▓██████▓▒░ ░▒▓██████▓▒░ ░▒▓███████▓▒░▒▓███████▓▒░░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░▒▓████████▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒▒▓███▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓██████▓▒░░▒▓███████▓▒░ ░▒▓██████▓▒░  
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░         ░▒▓█▓▒░     
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░         ░▒▓█▓▒░     
░▒▓███████▓▒░ ░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓██████▓▒░ ░▒▓██████▓▒░░▒▓███████▓▒░░▒▓█▓▒░         ░▒▓█▓▒░     
																														   `
	fmt.Println(logo)

	fmt.Printf("[v%v]\n", color.BlueString("0.1.0"))

	var targetUrl = flag.String("u", "", "Target url with protocol ex: https://target.com")
	var withVersion = flag.Bool("version", false, "Get the version of Django on the target")
	var withListPackages = flag.Bool("list-packages", false, "Get the package's list on the target")
	var maxThread = flag.Int("max-thread", 100, "Set the max thread")
	var delay = flag.Int("delay", 0, "Set the delay between 2 requests (default : 0)")

	flag.Parse()

	if *targetUrl == "" {
		fmt.Printf("[%v] Veuillez spécifier une target\n", color.RedString("err"))
		return
	}

	//TODO:vérifier si url OK

	//TODO:vérifier si version ou package sont ok

	//Lancer la récupération de version en 1er

	//Lancer la récupération des packages en 2nd
}
