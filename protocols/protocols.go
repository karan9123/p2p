/*
package protocols

import (

	"fmt"
	"os"
	"strings"

)

	type protocol struct {
		Name    string
		Code    int
		Version string
	}

	func GetProtocols(filename string) ([]protocol, error) {
		// Read the file contents into a byte slice
		fileContent, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Could not get contents from file due to %s \n", err.Error())
			return nil, err
		}

		// Convert the byte slice to a string and split it into lines
		lines := strings.Split(string(fileContent), "\n")

		// Create an empty list of protocols
		protocols := []protocol{}

		// Iterate over each line in the file
		for _, line := range lines {
			// Skip empty lines
			if len(line) == 0 {
				continue
			}

			// Split the line into comma-separated values
			fields := strings.Split(line, ",")

			// Create a new protocol struct and populate its fields
			p := protocol{
				Name:    fields[0],
				Code:    parseInt(fields[1]),
				Version: fields[2],
			}

			// Add the protocol to the list
			protocols = append(protocols, p)
		}

		// Print the list of protocols
		return protocols, nil
	}

// parseInt converts a string to an integer
// unsafe function, error not handled

	func parseInt(s string) int {
		var result int
		_, err := fmt.Sscanf(s, "%d", &result)
		if err != nil {
			fmt.Printf("error converting string to integer\n")
			return 0
		}
		return result
	}
*/
package protocol
