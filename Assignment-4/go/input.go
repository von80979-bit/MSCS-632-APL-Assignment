// The start menu, the typed input, the ranking parser, the slot maximum variable, and the sample employees.

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

var errInvalidRanking = errors.New("Use 1 to 3 different letters from M, A, and E, or - for a day off.")

var shiftLetters = map[rune]Shift{'M': Morning, 'A': Afternoon, 'E': Evening}

// console asks questions on out and reads one line per answer from in.
type console struct {
	lines *bufio.Scanner
	out   io.Writer
}

// ask prints the prompt and returns the answer without outer spaces. It returns false at the end of input.
func (console console) ask(prompt string) (string, bool) {
	fmt.Fprint(console.out, prompt)
	if !console.lines.Scan() {
		// The answer did not end the prompt line, so end it here.
		fmt.Fprintln(console.out)
		return "", false
	}
	return strings.TrimSpace(console.lines.Text()), true
}

// readEmployees runs the start menu and the typed input, and returns the employees in entry order.
func readEmployees(in io.Reader, out io.Writer) []Employee {
	console := console{lines: bufio.NewScanner(in), out: out}
	fmt.Fprint(out, "Employee Scheduler\n\n1) Load the sample employees\n2) Start with no employees\n")
	for {
		answer, ok := console.ask("Choose 1 or 2: ")
		switch {
		case !ok:
			return nil
		case answer == "1":
			employees := sampleEmployees()
			fmt.Fprintf(out, "Loaded %d sample employees.\n", len(employees))
			return console.addMore(employees)
		case answer == "2":
			return console.typeEmployees(nil)
		default:
			fmt.Fprintln(out, "Enter 1 or 2.")
		}
	}
}

func (console console) addMore(employees []Employee) []Employee {
	for {
		answer, ok := console.ask("Add more employees? (y/n): ")
		if !ok {
			return employees
		}
		switch strings.ToLower(answer) {
		case "n":
			return employees
		case "y":
			return console.typeEmployees(employees)
		default:
			fmt.Fprintln(console.out, "Enter y or n.")
		}
	}
}

// typeEmployees asks for employees until a blank name or the end of input. The end of input drops an employee
// whose days are not all entered.
func (console console) typeEmployees(employees []Employee) []Employee {
	for {
		name, ok := console.ask("Employee name (blank to finish): ")
		if !ok || name == "" {
			return employees
		}
		isDuplicate := func(employee Employee) bool { return strings.EqualFold(employee.Name, name) }
		if slices.ContainsFunc(employees, isDuplicate) {
			fmt.Fprintf(console.out, "Name already exists: %s.\n", name)
			continue
		}
		preferences, ok := console.askPreferences(name)
		if !ok {
			return employees
		}
		employees = append(employees, Employee{Name: name, Preferences: preferences})
	}
}

// askPreferences asks a ranking for each day until the employee has 5 working days, even when the fifth is Sunday.
// It returns false at the end of input.
func (console console) askPreferences(name string) (map[Day]Preference, bool) {
	preferences := map[Day]Preference{}
	for day := range Day(daysPerWeek) {
		for {
			answer, ok := console.ask(fmt.Sprintf("%v ranking (M A E, or - for a day off): ", day))
			if !ok {
				return nil, false
			}
			preference, isDayOff, err := parseRanking(answer)
			if err != nil {
				fmt.Fprintln(console.out, err)
				continue
			}
			if !isDayOff {
				preferences[day] = preference
			}
			break
		}
		if len(preferences) == maxDaysPerWeek {
			fmt.Fprintf(console.out, "%s has %d working days. The remaining days are days off.\n", name, maxDaysPerWeek)
			break
		}
	}
	return preferences, true
}

// slotMaximum reads the value of MAX_EMPLOYEES_PER_SLOT. It returns the default of 3 when the variable is not set,
// and returns false when the value is not a whole number of 2 or more.
func slotMaximum(value string, isSet bool) (int, bool) {
	if !isSet {
		return defaultEmployeesPerSlot, true
	}
	maximum, err := strconv.Atoi(value)
	if err != nil || maximum < minEmployeesPerSlot {
		return defaultEmployeesPerSlot, false
	}
	return maximum, true
}

// parseRanking reads one day's ranking: 1 to 3 different letters from M, A, and E in any case, or - for a day off.
// Spaces do not count. Missing shifts go after the given letters, in the order morning, afternoon, evening.
func parseRanking(text string) (preference Preference, isDayOff bool, err error) {
	letters := strings.ToUpper(strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, text))
	if letters == "-" {
		return Preference{}, true, nil
	}
	if len(letters) == 0 || len(letters) > shiftsPerDay {
		return Preference{}, false, errInvalidRanking
	}
	var ranked []Shift
	for _, letter := range letters {
		shift, known := shiftLetters[letter]
		if !known || slices.Contains(ranked, shift) {
			return Preference{}, false, errInvalidRanking
		}
		ranked = append(ranked, shift)
	}
	for shift := range Shift(shiftsPerDay) {
		if !slices.Contains(ranked, shift) {
			ranked = append(ranked, shift)
		}
	}
	return Preference(ranked), false, nil
}

// sampleEmployees returns the 13 sample employees in entry order. With the default maximum of 3, one run shows every
// change type.
func sampleEmployees() []Employee {
	mae := Preference{Morning, Afternoon, Evening}
	mea := Preference{Morning, Evening, Afternoon}
	ame := Preference{Afternoon, Morning, Evening}
	aem := Preference{Afternoon, Evening, Morning}
	ema := Preference{Evening, Morning, Afternoon}
	eam := Preference{Evening, Afternoon, Morning}
	return []Employee{
		{"Ana", map[Day]Preference{Monday: mae, Tuesday: mae, Wednesday: mae, Saturday: ame, Sunday: ema}},
		{"Ben", map[Day]Preference{Monday: mae, Tuesday: mae, Thursday: mae, Saturday: ema, Sunday: mae}},
		{"Cara", map[Day]Preference{Monday: mea, Tuesday: mae, Friday: mae, Saturday: ame, Sunday: ema}},
		{"Dan", map[Day]Preference{Monday: mae, Wednesday: mae, Thursday: mae, Saturday: mae, Sunday: ame}},
		{"Eve", map[Day]Preference{Monday: ame, Wednesday: mae, Friday: mae, Saturday: mae, Sunday: ema}},
		{"Finn", map[Day]Preference{Monday: aem, Wednesday: ame, Thursday: mae, Saturday: ema, Sunday: mae}},
		{"Gia", map[Day]Preference{Monday: eam, Wednesday: ame, Friday: ame, Saturday: mae, Sunday: ame}},
		{"Hugo", map[Day]Preference{Monday: aem, Thursday: ame, Friday: ame, Saturday: ame, Sunday: mae}},
		{"Ivy", map[Day]Preference{Monday: ema, Thursday: ame, Friday: ema, Saturday: ema, Sunday: ame}},
		{"Jay", map[Day]Preference{Tuesday: ame}},
		{"Kai", map[Day]Preference{Monday: eam, Wednesday: ame, Thursday: ema, Friday: ema}},
		{"Leo", map[Day]Preference{Monday: mae, Tuesday: ame, Wednesday: ema, Thursday: ema, Friday: mae}},
		{"Mia", map[Day]Preference{Tuesday: ame}},
	}
}
