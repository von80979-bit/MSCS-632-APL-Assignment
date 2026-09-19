// The 13 shared test cases from the specification.

package main

import (
	"fmt"
	"maps"
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
)

var mae = Preference{Morning, Afternoon, Evening}

const seedCount = 50

func seededRandom(seed uint64) *rand.Rand {
	return rand.New(rand.NewPCG(seed, seed))
}

// employeesRanking gives count employees named after prefix, each with the same preference on the given days.
func employeesRanking(prefix string, count int, preference Preference, days ...Day) []Employee {
	employees := make([]Employee, count)
	for i := range employees {
		employees[i] = employeeRanking(fmt.Sprintf("%s%d", prefix, i+1), preference, days...)
	}
	return employees
}

func employeeRanking(name string, preference Preference, days ...Day) Employee {
	preferences := map[Day]Preference{}
	for _, day := range days {
		preferences[day] = preference
	}
	return Employee{Name: name, Preferences: preferences}
}

// shiftOn returns the shift the employee works on the day, or false for no assignment.
func shiftOn(schedule Schedule, name string, day Day) (Shift, bool) {
	for shift, slot := range schedule.Slots[day] {
		for _, index := range slot {
			if schedule.Employees[index].Name == name {
				return Shift(shift), true
			}
		}
	}
	return 0, false
}

func TestRankFallback(t *testing.T) {
	employees := employeesRanking("E", 4, mae, Monday)

	schedule := buildSchedule(employees, 3, seededRandom(1))

	if shift, _ := shiftOn(schedule, "E4", Monday); shift != Afternoon {
		t.Fatalf("E4 works Monday %v, want afternoon", shift)
	}
	want := RankFallback{Employee: "E4", Day: Monday, FirstChoice: Morning, Placed: Afternoon, Choice: 2}
	if !slices.Contains(schedule.Changes, Change(want)) {
		t.Fatalf("changes %v do not contain %v", schedule.Changes, want)
	}
}

func TestDayMoveToNextDay(t *testing.T) {
	employees := employeesRanking("E", 10, mae, Monday)

	schedule := buildSchedule(employees, 3, seededRandom(1))

	if shift, placed := shiftOn(schedule, "E10", Tuesday); !placed || shift != Morning {
		t.Fatalf("E10 works Tuesday %v (placed %v), want morning", shift, placed)
	}
	want := DayMove{Employee: "E10", RankedDay: Monday, Day: Tuesday, Shift: Morning}
	if !slices.Contains(schedule.Changes, Change(want)) {
		t.Fatalf("changes %v do not contain %v", schedule.Changes, want)
	}
}

func TestDayMovePastWorkingDay(t *testing.T) {
	employees := employeesRanking("E", 9, mae, Monday)
	employees = append(employees, employeeRanking("E10", mae, Monday, Tuesday))

	schedule := buildSchedule(employees, 3, seededRandom(1))

	want := DayMove{Employee: "E10", RankedDay: Monday, Day: Wednesday, Shift: Morning}
	if !slices.Contains(schedule.Changes, Change(want)) {
		t.Fatalf("changes %v do not contain %v", schedule.Changes, want)
	}
}

func TestUnplaced(t *testing.T) {
	employees := employeesRanking("Full", 9, mae, Monday, Saturday, Sunday)
	employees = append(employees, employeeRanking("Late", mae, Monday, Tuesday, Wednesday, Thursday, Friday))

	schedule := buildSchedule(employees, 3, seededRandom(1))

	for _, day := range []Day{Monday, Saturday, Sunday} {
		if _, placed := shiftOn(schedule, "Late", day); placed {
			t.Fatalf("Late works %v, want no assignment", day)
		}
	}
	want := Unplaced{Employee: "Late", Day: Monday}
	if !slices.Contains(schedule.Changes, Change(want)) {
		t.Fatalf("changes %v do not contain %v", schedule.Changes, want)
	}
}

func TestRandomFill(t *testing.T) {
	employees := []Employee{
		employeeRanking("Placed", mae, Monday),
		employeeRanking("Free1", mae),
		employeeRanking("Free2", mae),
		employeeRanking("FullWeek", mae, Tuesday, Wednesday, Thursday, Friday, Saturday),
		employeeRanking("SameDay", Preference{Afternoon, Morning, Evening}, Monday),
	}
	eligible := []string{"Free1", "Free2"}

	for seed := range uint64(seedCount) {
		schedule := buildSchedule(employees, 3, seededRandom(seed))

		if got := len(schedule.Slots[Monday][Morning]); got != 2 {
			t.Fatalf("seed %d: Monday morning has %d employees, want 2", seed, got)
		}
		added := schedule.Employees[schedule.Slots[Monday][Morning][1]].Name
		if !slices.Contains(eligible, added) {
			t.Fatalf("seed %d: random fill added %s to Monday morning, want Free1 or Free2", seed, added)
		}
		want := RandomFill{Employee: added, Day: Monday, Shift: Morning}
		if !slices.Contains(schedule.Changes, Change(want)) {
			t.Fatalf("seed %d: changes %v do not contain %v", seed, schedule.Changes, want)
		}
	}
}

func TestTransfer(t *testing.T) {
	employees := slices.Concat(
		employeesRanking("Morning", 3, mae, Monday),
		employeesRanking("Afternoon", 2, Preference{Afternoon, Morning, Evening}, Monday),
		employeesRanking("Evening", 1, Preference{Evening, Morning, Afternoon}, Monday),
	)

	for seed := range uint64(seedCount) {
		schedule := buildSchedule(employees, 3, seededRandom(seed))

		if got := len(schedule.Slots[Monday][Evening]); got != 2 {
			t.Fatalf("seed %d: Monday evening has %d employees, want 2", seed, got)
		}
		if got := len(schedule.Slots[Monday][Morning]); got != 2 {
			t.Fatalf("seed %d: Monday morning has %d employees, want 2", seed, got)
		}
		moved := schedule.Employees[schedule.Slots[Monday][Evening][1]].Name
		want := Transfer{Employee: moved, FromDay: Monday, FromShift: Morning, ToDay: Monday, ToShift: Evening}
		if !slices.Contains(schedule.Changes, Change(want)) {
			t.Fatalf("seed %d: changes %v do not contain %v", seed, schedule.Changes, want)
		}
	}
}

func TestNoDonor(t *testing.T) {
	employees := []Employee{employeeRanking("Solo", mae, Monday)}

	schedule := buildSchedule(employees, 3, seededRandom(1))

	if got := len(schedule.Slots[Monday][Morning]); got != 1 {
		t.Fatalf("Monday morning has %d employees, want 1", got)
	}
	want := StillShort{Day: Monday, Shift: Morning, Employees: 1}
	if !slices.Contains(schedule.Changes, Change(want)) {
		t.Fatalf("changes %v do not contain %v", schedule.Changes, want)
	}
	if line, wantLine := changeLine(want), "Monday morning is still short: 1 of 2 employees."; line != wantLine {
		t.Fatalf("still-short line %q, want %q", line, wantLine)
	}
}

func TestOneShiftPerDay(t *testing.T) {
	for seed := range uint64(seedCount) {
		schedule := buildSchedule(sampleEmployees(), 3, seededRandom(seed))

		for day := range Day(daysPerWeek) {
			var names []string
			for _, slot := range schedule.Slots[day] {
				for _, employee := range slot {
					name := schedule.Employees[employee].Name
					if slices.Contains(names, name) {
						t.Fatalf("seed %d: %s works twice on %v", seed, name, day)
					}
					names = append(names, name)
				}
			}
		}
	}
}

func TestFiveDayLimit(t *testing.T) {
	for seed := range uint64(seedCount) {
		schedule := buildSchedule(sampleEmployees(), 3, seededRandom(seed))

		days := map[string]int{}
		for day := range Day(daysPerWeek) {
			for _, slot := range schedule.Slots[day] {
				for _, employee := range slot {
					days[schedule.Employees[employee].Name]++
				}
			}
		}
		for name, count := range days {
			if count > 5 {
				t.Fatalf("seed %d: %s works %d days, want 5 at most", seed, name, count)
			}
		}
	}
}

func TestSampleCoverage(t *testing.T) {
	for seed := range uint64(seedCount) {
		schedule := buildSchedule(sampleEmployees(), 3, seededRandom(seed))

		counts := map[string]int{}
		for _, change := range schedule.Changes {
			counts[fmt.Sprintf("%T", change)]++
		}
		want := map[string]int{
			"main.RankFallback": 2,
			"main.DayMove":      1,
			"main.Unplaced":     1,
			"main.RandomFill":   1,
			"main.Transfer":     1,
		}
		if !maps.Equal(counts, want) {
			t.Fatalf("seed %d: change counts %v, want %v", seed, counts, want)
		}
		dayMove := DayMove{Employee: "Kai", RankedDay: Monday, Day: Tuesday, Shift: Evening}
		unplaced := Unplaced{Employee: "Leo", Day: Monday}
		for _, change := range []Change{dayMove, unplaced} {
			if !slices.Contains(schedule.Changes, change) {
				t.Fatalf("seed %d: changes %v do not contain %v", seed, schedule.Changes, change)
			}
		}
	}
}

func TestSlotMaximumVariable(t *testing.T) {
	cases := []struct {
		value     string
		isSet     bool
		want      int
		wantValid bool
	}{
		{"4", true, 4, true},
		{"", false, 3, true},
		{"abc", true, 3, false},
		{"1", true, 3, false},
		{"-2", true, 3, false},
	}
	for _, c := range cases {
		got, valid := parseSlotMaximum(c.value, c.isSet)
		if got != c.want || valid != c.wantValid {
			t.Errorf("parseSlotMaximum(%q, %v) = %d, %v, want %d, %v",
				c.value, c.isSet, got, valid, c.want, c.wantValid)
		}
	}
}

func TestRankingInput(t *testing.T) {
	accepted := []struct {
		text string
		want Preference
	}{
		{"ma", Preference{Morning, Afternoon, Evening}},
		{"e", Preference{Evening, Morning, Afternoon}},
		{"m a", Preference{Morning, Afternoon, Evening}},
		{"E M A", Preference{Evening, Morning, Afternoon}},
	}
	for _, c := range accepted {
		got, isDayOff, err := parseRanking(c.text)
		if err != nil || isDayOff || got != c.want {
			t.Errorf("parseRanking(%q) = %v, %v, %v, want %v", c.text, got, isDayOff, err, c.want)
		}
	}

	for _, text := range []string{"mm", "x", "mae e", "", "  "} {
		if _, _, err := parseRanking(text); err == nil {
			t.Errorf("parseRanking(%q) gives no error, want a rejection", text)
		}
	}

	if _, isDayOff, err := parseRanking("-"); err != nil || !isDayOff {
		t.Errorf("parseRanking(\"-\") = day off %v, %v, want a day off", isDayOff, err)
	}
}

func TestPreferenceTable(t *testing.T) {
	employees := []Employee{
		employeeRanking("Ana", mae, Monday, Sunday),
		employeeRanking("Bo", Preference{Evening, Morning, Afternoon}, Tuesday),
	}

	var table strings.Builder
	writeGrid(&table, preferenceRows(employees))

	want := strings.Join([]string{
		"Employee | Monday | Tuesday | Wednesday | Thursday | Friday | Saturday | Sunday",
		"---------+--------+---------+-----------+----------+--------+----------+-------",
		"Ana      | M A E  | -       | -         | -        | -      | -        | M A E",
		"Bo       | -      | E M A   | -         | -        | -      | -        | -",
		"",
	}, "\n")
	if table.String() != want {
		t.Fatalf("preference table\n%s\nwant\n%s", table.String(), want)
	}
}
