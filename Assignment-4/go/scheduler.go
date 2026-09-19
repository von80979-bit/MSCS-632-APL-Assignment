// The four passes: preferences, day moves, random fill, and transfer.

package main

import (
	"math/rand/v2"
	"slices"
)

type placement struct {
	employee int
	day      Day
	shift    Shift
}

// overflow is an employee whose ranked day was full in pass 1.
type overflow struct {
	employee int
	day      Day
}

// buildSchedule runs the four passes over the week. The caller passes the random source, so a test can seed it.
func buildSchedule(employees []Employee, maxPerSlot int, random *rand.Rand) Schedule {
	schedule := Schedule{Employees: employees, MaxPerSlot: maxPerSlot}
	overflows := schedule.placePreferences()
	schedule.moveDays(overflows)
	schedule.fillAtRandom(random)
	schedule.transfer(random)
	return schedule
}

// placePreferences is pass 1. It places each employee in the best ranked shift with room, in entry order, and
// returns the employees whose ranked day was full.
func (schedule *Schedule) placePreferences() []overflow {
	var overflows []overflow
	for day := range Day(daysPerWeek) {
		for employee, employeeData := range schedule.Employees {
			preference, ranked := employeeData.Preferences[day]
			if !ranked {
				continue
			}
			if !schedule.placeByRank(employee, day, preference) {
				overflows = append(overflows, overflow{employee: employee, day: day})
			}
		}
	}
	return overflows
}

// placeByRank places the employee in the first shift of the preference that has room, and records a rank fallback
// when that shift is not the first choice.
func (schedule *Schedule) placeByRank(employee int, day Day, preference Preference) bool {
	rank, found := schedule.firstRankWithRoom(day, preference)
	if !found {
		return false
	}
	schedule.assign(employee, day, preference[rank])
	if rank > 0 {
		schedule.Changes = append(schedule.Changes, RankFallback{
			Employee:    schedule.Employees[employee].Name,
			Day:         day,
			FirstChoice: preference[0],
			Placed:      preference[rank],
			Choice:      rank + 1,
		})
	}
	return true
}

// moveDays is pass 2. It places each overflow on the first other day with room, starting from the next day and
// wrapping to Monday, because the schedule covers one week. An overflow with no such day stays unplaced.
func (schedule *Schedule) moveDays(overflows []overflow) {
	for _, overflow := range overflows {
		if !schedule.moveToOtherDay(overflow) {
			schedule.Changes = append(schedule.Changes, Unplaced{
				Employee: schedule.Employees[overflow.employee].Name,
				Day:      overflow.day,
			})
		}
	}
}

func (schedule *Schedule) moveToOtherDay(overflow overflow) bool {
	preference := schedule.Employees[overflow.employee].Preferences[overflow.day]
	for offset := 1; offset < daysPerWeek; offset++ {
		day := (overflow.day + Day(offset)) % daysPerWeek
		if schedule.hasAssignment(overflow.employee, day) {
			continue
		}
		rank, found := schedule.firstRankWithRoom(day, preference)
		if !found {
			continue
		}
		shift := preference[rank]
		schedule.assign(overflow.employee, day, shift)
		schedule.Changes = append(schedule.Changes, DayMove{
			Employee:  schedule.Employees[overflow.employee].Name,
			RankedDay: overflow.day,
			Day:       day,
			Shift:     shift,
		})
		return true
	}
	return false
}

// fillAtRandom is pass 3. It adds random candidates to each short slot until the slot has 2 employees or no
// candidate is left. A candidate has no assignment on that day and fewer than 5 days.
func (schedule *Schedule) fillAtRandom(random *rand.Rand) {
	for day := range Day(daysPerWeek) {
		for shift := range Shift(shiftsPerDay) {
			for schedule.isShort(day, shift) {
				var candidates []int
				for employee := range schedule.Employees {
					if !schedule.hasAssignment(employee, day) && schedule.daysWorked(employee) < maxDaysPerWeek {
						candidates = append(candidates, employee)
					}
				}
				if len(candidates) == 0 {
					break
				}
				employee := candidates[random.IntN(len(candidates))]
				schedule.assign(employee, day, shift)
				schedule.Changes = append(schedule.Changes, RandomFill{
					Employee: schedule.Employees[employee].Name,
					Day:      day,
					Shift:    shift,
				})
			}
		}
	}
}

// transfer is pass 4. It moves employees from donor slots into each slot that is still short, until the slot has 2
// employees or no donor is left. A donor slot has more than 2 employees, so it never becomes short. A slot without a
// donor stays short.
func (schedule *Schedule) transfer(random *rand.Rand) {
	for day := range Day(daysPerWeek) {
		for shift := range Shift(shiftsPerDay) {
			for schedule.isShort(day, shift) {
				candidates := schedule.transferCandidates(day, shift)
				if len(candidates) == 0 {
					schedule.Changes = append(schedule.Changes, StillShort{
						Day:       day,
						Shift:     shift,
						Employees: len(schedule.Slots[day][shift]),
					})
					break
				}
				moved := candidates[random.IntN(len(candidates))]
				donor := &schedule.Slots[moved.day][moved.shift]
				*donor = slices.DeleteFunc(*donor, func(employee int) bool { return employee == moved.employee })
				schedule.assign(moved.employee, day, shift)
				schedule.Changes = append(schedule.Changes, Transfer{
					Employee:  schedule.Employees[moved.employee].Name,
					FromDay:   moved.day,
					FromShift: moved.shift,
					ToDay:     day,
					ToShift:   shift,
				})
			}
		}
	}
}

// transferCandidates returns the employees in the donor slots of the same day. A move within the day keeps the
// one-shift-per-day rule. Without a donor on the same day, it returns the employees in donor slots on other days
// who have no assignment on the target day.
func (schedule *Schedule) transferCandidates(day Day, shift Shift) []placement {
	var candidates []placement
	for donorShift := range Shift(shiftsPerDay) {
		if donorShift != shift && schedule.isDonor(day, donorShift) {
			for _, employee := range schedule.Slots[day][donorShift] {
				candidates = append(candidates, placement{employee: employee, day: day, shift: donorShift})
			}
		}
	}
	if len(candidates) > 0 {
		return candidates
	}
	for donorDay := range Day(daysPerWeek) {
		if donorDay == day {
			continue
		}
		for donorShift := range Shift(shiftsPerDay) {
			if !schedule.isDonor(donorDay, donorShift) {
				continue
			}
			for _, employee := range schedule.Slots[donorDay][donorShift] {
				if !schedule.hasAssignment(employee, day) {
					candidates = append(candidates, placement{employee: employee, day: donorDay, shift: donorShift})
				}
			}
		}
	}
	return candidates
}

func (schedule *Schedule) isDonor(day Day, shift Shift) bool {
	return len(schedule.Slots[day][shift]) > minEmployeesPerSlot
}

func (schedule *Schedule) daysWorked(employee int) int {
	days := 0
	for day := range Day(daysPerWeek) {
		if schedule.hasAssignment(employee, day) {
			days++
		}
	}
	return days
}

func (schedule *Schedule) isShort(day Day, shift Shift) bool {
	return len(schedule.Slots[day][shift]) < minEmployeesPerSlot
}

// firstRankWithRoom returns the rank of the first shift in the preference with room on the day.
func (schedule *Schedule) firstRankWithRoom(day Day, preference Preference) (int, bool) {
	for rank, shift := range preference {
		if schedule.hasRoom(day, shift) {
			return rank, true
		}
	}
	return 0, false
}

func (schedule *Schedule) hasAssignment(employee int, day Day) bool {
	for _, slot := range schedule.Slots[day] {
		if slices.Contains(slot, employee) {
			return true
		}
	}
	return false
}

func (schedule *Schedule) hasRoom(day Day, shift Shift) bool {
	return len(schedule.Slots[day][shift]) < schedule.MaxPerSlot
}

func (schedule *Schedule) assign(employee int, day Day, shift Shift) {
	schedule.Slots[day][shift] = append(schedule.Slots[day][shift], employee)
}
