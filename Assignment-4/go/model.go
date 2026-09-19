// Days, shifts, employees, preferences, slots, the schedule, and the change types.

package main

import "strings"

// Day is one day of the week, from Monday to Sunday.
type Day int

const (
	Monday Day = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

const daysPerWeek = 7

var dayNames = [daysPerWeek]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

func (day Day) String() string { return dayNames[day] }

// Shift is one period of a day: morning, afternoon, or evening.
type Shift int

const (
	Morning Shift = iota
	Afternoon
	Evening
)

const shiftsPerDay = 3

var shiftNames = [shiftsPerDay]string{"morning", "afternoon", "evening"}

func (shift Shift) String() string { return shiftNames[shift] }

// title is the shift name as the schedule table header shows it.
func (shift Shift) title() string {
	name := shift.String()
	return strings.ToUpper(name[:1]) + name[1:]
}

// letter is the shift letter of the ranking input, as the preference table shows it.
func (shift Shift) letter() string {
	return strings.ToUpper(shift.String()[:1])
}

const (
	minEmployeesPerSlot     = 2
	defaultEmployeesPerSlot = 3
	maxDaysPerWeek          = 5
)

// Preference is the three shifts of one day in rank order, first choice first.
type Preference [shiftsPerDay]Shift

// Employee has a name and a preference for each day the employee wants to work. A day without a preference is a
// day off.
type Employee struct {
	Name        string
	Preferences map[Day]Preference
}

// Slot holds the indexes into Schedule.Employees of the employees assigned to one shift on one day, in order of
// assignment.
type Slot []int

// Schedule is the week's 21 slots, and the changes the scheduler made in the order of the checks.
type Schedule struct {
	Employees   []Employee
	SlotMaximum int
	Slots       [daysPerWeek][shiftsPerDay]Slot
	Changes     []Change
}

// Change is one change the scheduler made to the preferences. The unexported method keeps the set of change types
// closed to this package, so a type switch over the types below covers every case.
type Change interface {
	isChange()
}

// RankFallback places an employee in the second or third choice, because a better choice is full.
type RankFallback struct {
	Employee    string
	Day         Day
	FirstChoice Shift
	Placed      Shift
	Choice      int
}

// DayMove places an employee on another day, because every slot on the ranked day is full.
type DayMove struct {
	Employee  string
	RankedDay Day
	Day       Day
	Shift     Shift
}

// Unplaced records an employee whose ranked day is full and for whom no other day has room.
type Unplaced struct {
	Employee string
	Day      Day
}

// RandomFill adds a random eligible employee to a short slot.
type RandomFill struct {
	Employee string
	Day      Day
	Shift    Shift
}

// Transfer moves an employee from a slot with more than 2 employees into a short slot.
type Transfer struct {
	Employee  string
	FromDay   Day
	FromShift Shift
	ToDay     Day
	ToShift   Shift
}

// StillShort records a slot that stays short after every check.
type StillShort struct {
	Day       Day
	Shift     Shift
	Employees int
}

func (RankFallback) isChange() {}
func (DayMove) isChange()      {}
func (Unplaced) isChange()     {}
func (RandomFill) isChange()   {}
func (Transfer) isChange()     {}
func (StillShort) isChange()   {}
