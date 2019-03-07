package text

import "fmt"

type Generator func() string

func LoadingText() Generator {
	counter := 0
	list := []string{"-", "\\", "|", "/"}
	return func() string {
		counter = (counter + 1) % len(list)
		return fmt.Sprintf("loading %s", list[counter])
	}
}
