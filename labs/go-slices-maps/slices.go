package main


import (
	"golang.org/x/tour/pic"
)

func Pic(dx, dy int) [][]uint8 {
	pixel := make([][]uint8, dy)
	data := make([]uint8, dx)

	for i := range pixel {		//Me va a dar el index en dy
		for j := range data {	//Indice en dx
			data[j] = uint8(i*j)	//Funcion x*y
		}
		pixel[i] = data			//Se la agrego a cada pixel en el rango
	}
	return pixel
}
func main() {
	pic.Show(Pic)
}

