// The schedule table, the changes, and the employee summary.

package main

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

func writeSchedule(out io.Writer, schedule Schedule) {
	fmt.Fprintf(out, "Weekly schedule (maximum %d employees per slot)\n\n", schedule.MaxPerSlot)
	writeTable(out, schedule)
	fmt.Fprintln(out)
	writeChanges(out, schedule.Changes)
	fmt.Fprintln(out)
	writeSummary(out, schedule)
}

func writeTable(out io.Writer, schedule Schedule) {
	rows := [][]string{{"Day", "Morning", "Afternoon", "Evening"}}
	for day := range Day(daysPerWeek) {
		row := []string{day.String()}
		for _, slot := range schedule.Slots[day] {
			row = append(row, cell(schedule, slot))
		}
		rows = append(rows, row)
	}

	widths := make([]int, len(rows[0]))
	for _, row := range rows {
		for column, text := range row {
			widths[column] = max(widths[column], utf8.RuneCountInString(text))
		}
	}

	// The first and last columns have one space of padding next to the dashes, and the middle columns have two.
	separator := make([]string, len(widths))
	for column, width := range widths {
		padding := 2
		if column == 0 || column == len(widths)-1 {
			padding = 1
		}
		separator[column] = strings.Repeat("-", width+padding)
	}

	writeRow(out, rows[0], widths)
	fmt.Fprintln(out, strings.Join(separator, "+"))
	for _, row := range rows[1:] {
		writeRow(out, row, widths)
	}
}

func writeRow(out io.Writer, row []string, widths []int) {
	cells := make([]string, len(row))
	for column, text := range row {
		cells[column] = fmt.Sprintf("%-*s", widths[column], text)
	}
	fmt.Fprintln(out, strings.TrimRight(strings.Join(cells, " | "), " "))
}

func cell(schedule Schedule, slot Slot) string {
	names := make([]string, len(slot))
	for i, employee := range slot {
		names[i] = schedule.Employees[employee].Name
	}
	switch len(names) {
	case 0:
		return "(empty)"
	case 1:
		return names[0] + " (short)"
	default:
		return strings.Join(names, ", ")
	}
}

func writeChanges(out io.Writer, changes []Change) {
	fmt.Fprintln(out, "Changes")
	if len(changes) == 0 {
		fmt.Fprintln(out, "No changes.")
	}
	for _, change := range changes {
		fmt.Fprintln(out, changeLine(change))
	}
}

func changeLine(change Change) string {
	switch change := change.(type) {
	case RankFallback:
		return fmt.Sprintf("%s: %v %v was full. Placed in %v %v (choice %d).",
			change.Employee, change.Day, change.FirstChoice, change.Day, change.Placed, change.Choice)
	case DayMove:
		return fmt.Sprintf("%s: %v was full. Moved to %v %v.",
			change.Employee, change.RankedDay, change.Day, change.Shift)
	case Unplaced:
		return fmt.Sprintf("%s: %v was full and no other day had room. Not placed for %v.",
			change.Employee, change.Day, change.Day)
	case RandomFill:
		return fmt.Sprintf("%s: %v %v was short. Added at random.", change.Employee, change.Day, change.Shift)
	case Transfer:
		return fmt.Sprintf("%s: Moved from %v %v to %v %v to cover a short slot.",
			change.Employee, change.FromDay, change.FromShift, change.ToDay, change.ToShift)
	case StillShort:
		return fmt.Sprintf("%v %v is still short: %d of %d employees.",
			change.Day, change.Shift, change.Employees, minEmployeesPerSlot)
	default:
		panic(fmt.Sprintf("unknown change type %T", change))
	}
}

// writeSummary prints the days worked by each employee, in entry order.
func writeSummary(out io.Writer, schedule Schedule) {
	fmt.Fprintln(out, "Employees")
	nameWidth := 0
	for _, employee := range schedule.Employees {
		nameWidth = max(nameWidth, utf8.RuneCountInString(employee.Name))
	}
	for index, employee := range schedule.Employees {
		days := schedule.daysWorked(index)
		unit := "days"
		if days == 1 {
			unit = "day"
		}
		fmt.Fprintf(out, "%-*s    %d %s\n", nameWidth, employee.Name, days, unit)
	}
}
