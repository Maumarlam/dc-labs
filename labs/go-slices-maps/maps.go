package main
import (
	"strings"
	"golang.org/x/tour/wc"
)
func WordCount(s string) map[string]int {
	m := make(map[string]int) //Creo mi mapa que luego lleno
	x := strings.Fields(s) //Me regresa las palabras separadas por espacios 
	for _, n := range x { //Solo me importa el valor, no el indice, para contar
		m[n]++
	}
	return m
}
func main() {
	wc.Test(WordCount)
}
