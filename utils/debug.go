package utils

import (
	"fmt"
	"github.com/goccy/go-json"
	"os"
)

// DumpNodeWithExit prints a formatted JSON representation of the provided interface and exits the program.
// This function is primarily used for debugging purposes.
func DumpNodeWithExit(whatever any) {
	j, _ := json.MarshalIndent(whatever, "", "\t")
	fmt.Println(string(j))
	os.Exit(1)
}

// DumpNodeNoExit prints a formatted JSON representation of the provided interface without exiting the program.
// This function is primarily used for debugging purposes.
func DumpNodeNoExit(whatever any) {
	j, _ := json.MarshalIndent(whatever, "", "\t")
	fmt.Println(string(j))
}
