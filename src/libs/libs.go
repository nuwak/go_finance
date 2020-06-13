package libs

import "fmt"

func Contains(arr []interface{}, str interface{}) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func Print(name string, value float64) {
	fmt.Printf("%-10s: %g\n", name, value)
}
