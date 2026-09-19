//! The four checks: preferences, day moves, random fill, and transfer.

use rand::Rng;
use rand::seq::IndexedRandom;

use crate::model::{Change, Day, Employee, EmployeeId, MAXIMUM_DAYS_PER_WEEK, Preference, Schedule, Shift};

#[cfg(test)]
mod tests;

/// An employee whose ranked day was full in check 1.
struct Overflow {
    employee: EmployeeId,
    ranked_day: Day,
    preference: Preference,
}

pub fn build_schedule(employees: &[Employee], slot_maximum: usize, random: &mut impl Rng) -> Schedule {
    let mut schedule = Schedule::empty();
    let overflows = place_preferences(&mut schedule, employees, slot_maximum);
    move_overflows_to_other_days(&mut schedule, &overflows, slot_maximum);
    fill_short_slots_at_random(&mut schedule, employees.len(), random);
    transfer_into_short_slots(&mut schedule, random);
    schedule
}

/// Check 1. Finishes the whole week first, so a later day move never takes a slot from an employee who ranked it.
fn place_preferences(schedule: &mut Schedule, employees: &[Employee], slot_maximum: usize) -> Vec<Overflow> {
    let mut overflows = Vec::new();
    for ranked_day in Day::ALL {
        for (id, employee) in employees.iter().enumerate() {
            let Some(preference) = employee.preference_on(ranked_day) else {
                continue;
            };
            let open_choice =
                preference.iter().position(|&shift| schedule.slot(ranked_day, shift).has_room(slot_maximum));
            match open_choice {
                Some(choice_index) => {
                    let placed = preference[choice_index];
                    schedule.slot_mut(ranked_day, placed).employees.push(id);
                    if choice_index > 0 {
                        schedule.changes.push(Change::RankFallback {
                            employee: id,
                            day: ranked_day,
                            first_choice: preference[0],
                            placed,
                            choice: choice_index + 1,
                        });
                    }
                }
                None => overflows.push(Overflow { employee: id, ranked_day, preference }),
            }
        }
    }
    overflows
}

/// Check 2.
fn move_overflows_to_other_days(schedule: &mut Schedule, overflows: &[Overflow], slot_maximum: usize) {
    for overflow in overflows {
        let open_slot = other_days_in_search_order(overflow.ranked_day)
            .filter(|&day| !schedule.is_working_on(overflow.employee, day))
            .find_map(|day| {
                overflow
                    .preference
                    .iter()
                    .find(|&&shift| schedule.slot(day, shift).has_room(slot_maximum))
                    .map(|&shift| (day, shift))
            });
        match open_slot {
            Some((day, shift)) => {
                schedule.slot_mut(day, shift).employees.push(overflow.employee);
                schedule.changes.push(Change::DayMove {
                    employee: overflow.employee,
                    ranked_day: overflow.ranked_day,
                    day,
                    shift,
                });
            }
            None => {
                schedule.changes.push(Change::Unplaced { employee: overflow.employee, day: overflow.ranked_day });
            }
        }
    }
}

/// The other days of the week in search order: the days after the ranked day up to Sunday, then Monday onward.
fn other_days_in_search_order(ranked_day: Day) -> impl Iterator<Item = Day> {
    let next_index = ranked_day.index() + 1;
    Day::ALL[next_index..].iter().chain(&Day::ALL[..ranked_day.index()]).copied()
}

/// Check 3.
fn fill_short_slots_at_random(schedule: &mut Schedule, employee_count: usize, random: &mut impl Rng) {
    for (day, shift) in all_slots() {
        while schedule.slot(day, shift).is_short() {
            let candidates: Vec<EmployeeId> = (0..employee_count)
                .filter(|&employee| {
                    !schedule.is_working_on(employee, day) && schedule.days_worked(employee) < MAXIMUM_DAYS_PER_WEEK
                })
                .collect();
            let Some(&employee) = candidates.choose(random) else {
                break;
            };
            schedule.slot_mut(day, shift).employees.push(employee);
            schedule.changes.push(Change::RandomFill { employee, day, shift });
        }
    }
}

/// Check 4.
fn transfer_into_short_slots(schedule: &mut Schedule, random: &mut impl Rng) {
    for (day, shift) in all_slots() {
        while schedule.slot(day, shift).is_short() {
            let candidates = transfer_candidates(schedule, day);
            let Some(&(employee, from_day, from_shift)) = candidates.choose(random) else {
                schedule.changes.push(Change::StillShort {
                    day,
                    shift,
                    employee_count: schedule.slot(day, shift).employees.len(),
                });
                break;
            };
            schedule.slot_mut(from_day, from_shift).employees.retain(|&other| other != employee);
            schedule.slot_mut(day, shift).employees.push(employee);
            schedule.changes.push(Change::Transfer { employee, from_day, from_shift, to_day: day, to_shift: shift });
        }
    }
}

type TransferCandidate = (EmployeeId, Day, Shift);

/// Same-day donors come first, because a move within the day keeps the one-shift-per-day rule.
fn transfer_candidates(schedule: &Schedule, target_day: Day) -> Vec<TransferCandidate> {
    let candidates_from = |same_day: bool| -> Vec<TransferCandidate> {
        donor_slots(schedule)
            .filter(|&(donor_day, _)| (donor_day == target_day) == same_day)
            .flat_map(|(donor_day, donor_shift)| {
                schedule
                    .slot(donor_day, donor_shift)
                    .employees
                    .iter()
                    .map(move |&employee| (employee, donor_day, donor_shift))
            })
            .filter(|&(employee, _, _)| same_day || !schedule.is_working_on(employee, target_day))
            .collect()
    };

    let same_day_candidates = candidates_from(true);
    if same_day_candidates.is_empty() { candidates_from(false) } else { same_day_candidates }
}

fn donor_slots(schedule: &Schedule) -> impl Iterator<Item = (Day, Shift)> + '_ {
    all_slots().filter(|&(day, shift)| schedule.slot(day, shift).is_donor())
}

/// Every slot of the week, from Monday morning to Sunday evening.
fn all_slots() -> impl Iterator<Item = (Day, Shift)> {
    Day::ALL.into_iter().flat_map(|day| Shift::ALL.into_iter().map(move |shift| (day, shift)))
}
