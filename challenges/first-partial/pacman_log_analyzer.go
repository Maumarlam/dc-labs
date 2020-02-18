package main

import (
	"fmt"
	"os"
	"bufio"
	"log"
	"strings"
)

type Package struct {		//struct where all my package information will be at
	name string
	installDate string
	lastUpdate string
	updateCount int
	removalDate string
}


func check(e error){		//Function to verify if there is a nil and if not to continue
	if(e!=nil){
		fmt.Println("Error")
		os.Exit(-1)
	}
}

func write(toWrite string, file *os.File){	//to write into a txt file, got this from internet

	_,err:=file.WriteString(toWrite)
	check(err)

}


func main() {
	fmt.Println("Pacman Log Analyzer")

	if len(os.Args) < 2 {
		fmt.Println("You must send at least one pacman log file to analize")
		fmt.Println("usage: ./pacman_log_analizer <logfile>")
		os.Exit(1)
	}

	// Your fun starts here.

	//The variables for the first output
	var instalations, removals, upgrades, actualInstalls int

	m := make( map[string]*Package)		//Create a map to store the packages

	fileTxt, err := os.Create("packages_report.txt")
	check(err)

	file, err := os.Open("pacman.txt")	//Check if the file is usable
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)	//To read the file 
	for scanner.Scan() { 	//This reads line by line the file

		line := scanner.text()
		info := strings.Split(line, " ")//Gives me an array of the words in that line separated by a space

		var name, date string

		//An if to verify if its a package operation
		if(len(info) > 4) {
			name = info[4]
			date = line[0][1:len(line[0])] +" "+ line[1][:len(line[1])-1]
		}

		operation := info[3]	//What the package went through

		if(operation == "installed"){
			instalations++
			//Not going to lie, a friend helped me access the variables from the map
			m[name] = &Package{name, date, "-", 0, "-"}
		}

		if(operation == "upgraded") {
			m[name].lastUpdate = date
			m[name].updateCount++
			upgrades++
		}

		if(operation == "reinstalled"){
			m[name].installDate = date
			m[name].lastUpdate = "-"
			m[name].updateCount = 0
			m[name].removalDate = "-"
			removals--

		if(operation == "removed") {
			m[name].removalDate = date
			removals++
			instalations--
			upgrades -= m[name].updateCount
		}

	}

	check(scanner.Err())

	actualInstalls = instalations - removals

	//template of the Txt file
	writeToTxt("Pacman Packages Report\n",fileTxt)
	writeToTxt("----------------------\n",fileTxt)
	writeToTxt("- Installed packages : "+strconv.Itoa(installed)+"\n",fileTxt)
	writeToTxt("- Removed packages   : "+strconv.Itoa(removed)+"\n",fileTxt)
	writeToTxt("- Upgraded packages  : "+strconv.Itoa(upgraded)+"\n",fileTxt)
	writeToTxt("- Current installed  : "+strconv.Itoa(current)+"\n\n",fileTxt)

	for _,p := range m{
		writeToTxt("- Package Name        : "+p.name+"\n",fileTxt)
		writeToTxt("  - Install date      : "+p.installDate+"\n",fileTxt)
		writeToTxt("  - Last update date  : "+p.lastUpdate+"\n",fileTxt)
		writeToTxt("  - How many updates  : "+strconv.Itoa(p.updateCount)+"\n",fileTxt)
		writeToTxt("  - Removal date      : "+p.removalDate+"\n",fileTxt)

	}

	check(fileTxt.Close())

	fmt.Println("Program finished")

	return

 
}
