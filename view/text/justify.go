package text

import "strings"

/*
build Markov Chain data structures;
while( more words to generate ) {
   generate next word;
   if( word is short enough to fit on current output line )
      add word and trailing space(s) to the line;
         // Two spaces if it is the end of a sentence.  See below.
   else {
      add spaces to justify the line; // details in phase 3
      print the line;
      clear the linked list;
      add the word and trailing spaces to the line;
   }
}
if( output line is not emtpy )
   print output line;
*/
func Justify(str string, width int, center bool) []string {
	var lines []string
	var current []string
	currentLength := 0

	for _, word := range strings.Fields(str) {
		if len(word)+currentLength >= width && len(current) != 0 { // guard against super long word
			currentLine := strings.Join(current, " ")
			if center {
				prefix := strings.Repeat(" ", (width-len(currentLine))/2)
				suffix := strings.Repeat(" ", width-len(currentLine)-len(prefix))
				currentLine = prefix + currentLine + suffix
			}
			lines = append(lines, currentLine)
			current = make([]string, 0)
			currentLength = 0
		}
		current = append(current, word)
		currentLength += len(word) + 1
	}
	if len(current) != 0 {
		lines = append(lines, strings.Join(current, " "))
	}
	return lines
}
