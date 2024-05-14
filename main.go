package main

import (
	"flag"
	"fmt"
	"net/url"

	"github.com/djangospy/django"
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

	var targetUrlString = flag.String("u", "", "Target url with protocol ex: https://target.com")
	var withVersion = flag.Bool("version", false, "Get the version of Django on the target")
	var withListPackages = flag.Bool("list-packages", false, "Get the package's list on the target")
	//var maxThread = flag.Int("max-thread", 100, "Set the max thread")
	//var delay = flag.Int("delay", 0, "Set the delay between 2 requests (default : 0)")

	flag.Parse()

	if *targetUrlString == "" {
		fmt.Printf("[%v] No target provided\n", color.RedString("err"))
		return
	}

	targetUrl, err := url.Parse(*targetUrlString)
	if err != nil || targetUrl.Host == "" || targetUrl.Scheme == "" {
		fmt.Printf("[%v] Invalid url\n", color.RedString("err"))
		return
	}

	if !*withVersion && !*withListPackages {
		fmt.Printf("[%v] No action provided\n", color.RedString("err"))
		return
	}

	if *withVersion {
		django.GetDjangoVersion(targetUrl.Scheme + "://" + targetUrl.Host)
	}

	//Lancer la récupération des packages en 2nd
}
