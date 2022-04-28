package ini_parser

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

type Parser map[string]map[string]string

func (p Parser) readFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	section := "default"
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \n\t")
		if validateSection(line) {
			section = strings.Trim(line, "[]")
			_, ok := p[section]
			if !ok {
				p[section] = map[string]string{}
			}
		} else if validateKeyValPair(line) {
			keyValPair := strings.Split(line, "=")
			key := strings.Trim(keyValPair[0], " ")
			val := strings.Trim(keyValPair[1], " ")
			if !validateKey(key) || !validateValue(val) {
				log.Fatal(errors.New("syntax error"))
			}
			p[section][key] = val
		} else {
			log.Fatal(errors.New("syntax error"))
		}
	}
}

func (p Parser) writeToFile(filename string) {
	p.validateParser()
	p.readFile(filename)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	os.Truncate(filename, 0)
	for section, items := range p {
		_, err := file.WriteString("[" + section + "]\n")
		if err != nil {
			log.Fatal(err)
		}
		for key, val := range items {
			_, err := file.WriteString(key + " = " + val + "\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (p Parser) validateParser() {
	for section, items := range p {
		if !validateSection("[" + section + "]") {
			log.Fatal(errors.New("syntax error: section is not valid"))
		}
		for key, val := range items {
			if !validateKey(key) {
				log.Fatal(errors.New("syntax error: key is not valid"))
			}
			if !validateValue(val) {
				log.Fatal(errors.New("syntax error: value is not valid"))
			}
		}
	}

}

func validateSection(s string) bool {
	//validate section

	//check if string is empty
	if len(s) == 0 {
		return false
	}

	//check absence of square brackets
	if !strings.HasPrefix(s, "[") || !strings.HasSuffix(s, "]") {
		return false
	}

	//check if section string is []
	if s == "[]" {
		return false
	}

	//check presence of white spaces
	if hasWhiteSpaceOrSemiColon(s) {
		return false
	}

	return true
}

func validateKeyValPair(s string) bool {
	//validate key value pair line

	//check if string is empty
	if len(s) == 0 {
		return false
	}

	//check if string is like '[*]'
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		return false
	}

	//check presence of more than one '='
	numberOfEqualSigns := strings.Count(s, "=")
	if numberOfEqualSigns != 1 {
		return false
	}

	return true
}

func validateValue(s string) bool {
	//validate value

	// check if string is empty
	if len(s) == 0 {
		return false
	}

	//check presence of " " or ";"
	if hasWhiteSpaceOrSemiColon(s) {
		return false
	}

	return true
}

func validateKey(s string) bool {
	//validate key

	// check if string is empty
	if len(s) == 0 {
		return false
	}

	//check presence of " " or ";"
	if hasWhiteSpaceOrSemiColon(s) {
		return false
	}

	return true
}

func hasWhiteSpaceOrSemiColon(s string) bool {
	return strings.Contains(s, " ") || strings.Contains(s, ";")
}
