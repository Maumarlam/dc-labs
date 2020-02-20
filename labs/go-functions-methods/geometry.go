// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 156.

// Package geometry defines simple types for plane geometry.
//!+point
package main

import "math"
import "fmt"
import "math/rand"
import "os"
import "strconv"


//MAIN FUNCTION
func main() {
	sides := os.Args[1]
	points, err := strconv.Atoi(sides)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("- Generating a [%d] sides figure \n", points)
	aristas := make(Path, points)

	fmt.Println("- Figure's vertices")
	for i:=0; i<points; i++ {
		point := Point{rand.Float64() * 100, rand.Float64() * 100}
		fmt.Printf("    - (%f, %f) \n", point.X(), point.Y())
		aristas[i] = point
	}
	fmt.Println("- Figure's perimeter")
	fmt.Printf("    -")
	distance := aristas.Distance()
	fmt.Printf(" = %f \n" , distance)
}

// same thing, but as a method of the Point type
func (p Point) Distance(q Point) float64 {
	return math.Hypot(q.X() - p.X(), q.Y() - p.Y())
}

type Point struct{ x, y float64 }

//!-point
func (p Point) X() float64 {
	return p.x
}
//!+path
func (p Point) Y() float64 {
	return p.y
}
// A Path is a journey connecting the points with straight lines.

type Path []Point

// Distance returns the distance traveled along the path.
func (path Path) Distance() float64 {
	sum := 0.0
	for i := range path {
		if i > 0 {
			temp := path[i-1].Distance(path[i])
			if i< len(path)-1 {
				fmt.Printf(" %f +", temp)
			}else{
				fmt.Printf(" %f ", temp)
			}
			sum += temp
		}
	}
	return sum
}

//!-path
