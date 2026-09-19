//! The 13 shared test cases from the specification.

use rand::SeedableRng;
use rand::rngs::StdRng;

use super::build_schedule;
use crate::input::{InvalidSlotMaximum, parse_ranking, parse_slot_maximum, sample_employees};
use crate::model::{
    Change, DEFAULT_MAXIMUM_EMPLOYEES_PER_SLOT, Day, Employee, MAXIMUM_DAYS_PER_WEEK, MINIMUM_EMPLOYEES_PER_SLOT,
    Schedule, Shift,
};
use crate::output::render;

const MAXIMUM: usize = DEFAULT_MAXIMUM_EMPLOYEES_PER_SLOT;
const SEEDS: std::ops::Range<u64> = 0..50;

/// Builds an employee from one ranking per day, Monday first. `-` is a day off.
fn employee(name: &str, rankings: [&str; 7]) -> Employee {
    Employee { name: name.to_string(), preferences: rankings.map(|ranking| parse_ranking(ranking).unwrap()) }
}

fn schedule_with_seed(employees: &[Employee], seed: u64) -> Schedule {
    build_schedule(employees, MAXIMUM, &mut StdRng::seed_from_u64(seed))
}

fn id_of(employees: &[Employee], name: &str) -> usize {
    employees.iter().position(|employee| employee.name == name).unwrap()
}

#[test]
fn case_1_one_shift_per_day() {
    let employees = sample_employees();
    for seed in SEEDS {
        let schedule = schedule_with_seed(&employees, seed);
        for day in Day::ALL {
            for (id, employee) in employees.iter().enumerate() {
                let shifts_on_day = Shift::ALL
                    .iter()
                    .map(|&shift| schedule.slot(day, shift).employees.iter().filter(|&&other| other == id).count())
                    .sum::<usize>();
                assert!(shifts_on_day <= 1, "seed {seed}: {} works {shifts_on_day} shifts on {day:?}", employee.name);
            }
        }
    }
}

#[test]
fn case_2_five_day_limit() {
    let employees = sample_employees();
    for seed in SEEDS {
        let schedule = schedule_with_seed(&employees, seed);
        for (id, employee) in employees.iter().enumerate() {
            assert!(schedule.days_worked(id) <= MAXIMUM_DAYS_PER_WEEK, "seed {seed}: {}", employee.name);
        }
        for slot in &schedule.slots {
            assert!(
                slot.employees.len() <= MAXIMUM,
                "seed {seed}: {:?} {:?} is over the maximum",
                slot.day,
                slot.shift
            );
        }
    }
}

#[test]
fn case_3_rank_fallback() {
    let employees: Vec<Employee> =
        ["A", "B", "C", "D"].iter().map(|name| employee(name, ["MAE", "-", "-", "-", "-", "-", "-"])).collect();
    let schedule = schedule_with_seed(&employees, 0);

    // Check 4 later moves one morning employee into the short afternoon, so only the fourth placement is fixed.
    assert!(!schedule.slot(Day::Monday, Shift::Morning).employees.contains(&3));
    assert_eq!(schedule.slot(Day::Monday, Shift::Afternoon).employees.first(), Some(&3));
    assert!(schedule.changes.contains(&Change::RankFallback {
        employee: 3,
        day: Day::Monday,
        first_choice: Shift::Morning,
        placed: Shift::Afternoon,
        choice: 2,
    }));
}

fn ten_monday_employees(tenth_tuesday: &str) -> Vec<Employee> {
    let mut employees: Vec<Employee> =
        (1..=9).map(|number| employee(&format!("E{number}"), ["MAE", "-", "-", "-", "-", "-", "-"])).collect();
    employees.push(employee("Tenth", ["MAE", tenth_tuesday, "-", "-", "-", "-", "-"]));
    employees
}

#[test]
fn case_4_day_move_to_the_next_day() {
    let employees = ten_monday_employees("-");
    let schedule = schedule_with_seed(&employees, 0);

    assert!(schedule.slot(Day::Tuesday, Shift::Morning).employees.contains(&9));
    assert!(!schedule.is_working_on(9, Day::Monday));
    assert!(schedule.changes.contains(&Change::DayMove {
        employee: 9,
        ranked_day: Day::Monday,
        day: Day::Tuesday,
        shift: Shift::Morning,
    }));
}

#[test]
fn case_5_day_move_past_a_working_day() {
    let employees = ten_monday_employees("MAE");
    let schedule = schedule_with_seed(&employees, 0);

    assert!(schedule.slot(Day::Tuesday, Shift::Morning).employees.contains(&9));
    assert!(schedule.slot(Day::Wednesday, Shift::Morning).employees.contains(&9));
    assert!(schedule.changes.contains(&Change::DayMove {
        employee: 9,
        ranked_day: Day::Monday,
        day: Day::Wednesday,
        shift: Shift::Morning,
    }));
}

#[test]
fn case_6_unplaced() {
    // Late works Tuesday to Friday, so only Monday, Saturday, and Sunday are possible, and nine employees fill them.
    let mut employees: Vec<Employee> =
        (1..=9).map(|number| employee(&format!("Full{number}"), ["MAE", "-", "-", "-", "-", "MAE", "MAE"])).collect();
    employees.push(employee("Late", ["MAE", "MAE", "MAE", "MAE", "MAE", "-", "-"]));
    let late = id_of(&employees, "Late");
    let schedule = schedule_with_seed(&employees, 0);

    for day in [Day::Monday, Day::Saturday, Day::Sunday] {
        assert!(!schedule.is_working_on(late, day), "Late works {day:?}");
    }
    assert!(schedule.changes.contains(&Change::Unplaced { employee: late, day: Day::Monday }));
}

#[test]
fn case_7_random_fill() {
    let employees = vec![
        employee("Placed", ["MAE", "-", "-", "-", "-", "-", "-"]),
        employee("Free1", ["-", "-", "-", "-", "-", "-", "-"]),
        employee("Free2", ["-", "-", "-", "-", "-", "-", "-"]),
        employee("FullWeek", ["-", "MAE", "MAE", "MAE", "MAE", "MAE", "-"]),
        employee("SameDay", ["AME", "-", "-", "-", "-", "-", "-"]),
    ];
    let eligible = [id_of(&employees, "Free1"), id_of(&employees, "Free2")];

    for seed in SEEDS {
        let schedule = schedule_with_seed(&employees, seed);
        let monday_morning = &schedule.slot(Day::Monday, Shift::Morning).employees;

        assert_eq!(monday_morning.len(), MINIMUM_EMPLOYEES_PER_SLOT, "seed {seed}");
        assert!(eligible.contains(&monday_morning[1]), "seed {seed}: added employee {}", monday_morning[1]);
        assert!(schedule.changes.contains(&Change::RandomFill {
            employee: monday_morning[1],
            day: Day::Monday,
            shift: Shift::Morning,
        }));
    }
}

#[test]
fn case_8_transfer() {
    // Every employee works Monday, so no Monday slot has a candidate for the random fill.
    let employees = vec![
        employee("M1", ["MAE", "-", "-", "-", "-", "-", "-"]),
        employee("M2", ["MAE", "-", "-", "-", "-", "-", "-"]),
        employee("M3", ["MAE", "-", "-", "-", "-", "-", "-"]),
        employee("A1", ["AME", "-", "-", "-", "-", "-", "-"]),
        employee("A2", ["AME", "-", "-", "-", "-", "-", "-"]),
        employee("E1", ["EMA", "-", "-", "-", "-", "-", "-"]),
    ];

    for seed in SEEDS {
        let schedule = schedule_with_seed(&employees, seed);
        let monday_evening = &schedule.slot(Day::Monday, Shift::Evening).employees;

        assert_eq!(monday_evening.len(), MINIMUM_EMPLOYEES_PER_SLOT, "seed {seed}");
        assert_eq!(schedule.slot(Day::Monday, Shift::Morning).employees.len(), MINIMUM_EMPLOYEES_PER_SLOT);
        assert!([0, 1, 2].contains(&monday_evening[1]), "seed {seed}");
        assert!(schedule.changes.contains(&Change::Transfer {
            employee: monday_evening[1],
            from_day: Day::Monday,
            from_shift: Shift::Morning,
            to_day: Day::Monday,
            to_shift: Shift::Evening,
        }));
    }
}

#[test]
fn case_9_no_donor() {
    let employees = vec![employee("Only", ["MAE", "-", "-", "-", "-", "-", "-"])];
    let schedule = schedule_with_seed(&employees, 0);

    assert_eq!(schedule.slot(Day::Monday, Shift::Morning).employees, [0]);
    assert!(schedule.changes.contains(&Change::StillShort {
        day: Day::Monday,
        shift: Shift::Morning,
        employee_count: 1,
    }));
}

#[test]
fn case_10_slot_maximum_variable() {
    assert_eq!(parse_slot_maximum(Some("4")), Ok(4));
    assert_eq!(parse_slot_maximum(None), Ok(3));
    for invalid in ["abc", "1", "-2"] {
        assert_eq!(parse_slot_maximum(Some(invalid)), Err(InvalidSlotMaximum), "{invalid}");
        assert_eq!(parse_slot_maximum(Some(invalid)).unwrap_or(DEFAULT_MAXIMUM_EMPLOYEES_PER_SLOT), 3, "{invalid}");
    }
}

#[test]
fn case_11_ranking_input() {
    assert_eq!(parse_ranking("ma"), Ok(Some([Shift::Morning, Shift::Afternoon, Shift::Evening])));
    assert_eq!(parse_ranking("e"), Ok(Some([Shift::Evening, Shift::Morning, Shift::Afternoon])));
    assert_eq!(parse_ranking("m a"), parse_ranking("MA"));
    assert_eq!(parse_ranking("-"), Ok(None));
    for rejected in ["mm", "x", "mae e", ""] {
        assert!(parse_ranking(rejected).is_err(), "{rejected:?}");
    }
}

#[test]
fn case_12_sample_coverage() {
    let employees = sample_employees();
    for seed in SEEDS {
        let changes = schedule_with_seed(&employees, seed).changes;
        let count = |matches: fn(&Change) -> bool| changes.iter().filter(|&change| matches(change)).count();

        assert_eq!(count(|change| matches!(change, Change::RankFallback { .. })), 2, "seed {seed}: {changes:?}");
        assert_eq!(count(|change| matches!(change, Change::DayMove { .. })), 1, "seed {seed}: {changes:?}");
        assert_eq!(count(|change| matches!(change, Change::Unplaced { .. })), 1, "seed {seed}: {changes:?}");
        assert_eq!(count(|change| matches!(change, Change::RandomFill { .. })), 1, "seed {seed}: {changes:?}");
        assert_eq!(count(|change| matches!(change, Change::Transfer { .. })), 1, "seed {seed}: {changes:?}");
        assert_eq!(count(|change| matches!(change, Change::StillShort { .. })), 0, "seed {seed}: {changes:?}");

        let kai = id_of(&employees, "Kai");
        let leo = id_of(&employees, "Leo");
        assert!(changes.contains(&Change::DayMove {
            employee: kai,
            ranked_day: Day::Monday,
            day: Day::Tuesday,
            shift: Shift::Evening,
        }));
        assert!(changes.contains(&Change::Unplaced { employee: leo, day: Day::Monday }));
    }
}

#[test]
fn case_13_preference_table() {
    let employees = vec![
        employee("Ana", ["MAE", "-", "-", "-", "-", "-", "MAE"]),
        employee("Bo", ["-", "EMA", "-", "-", "-", "-", "-"]),
    ];
    let text = render(&Schedule::empty(), &employees, MAXIMUM);

    let table = "Employee preferences\n\n\
        Employee | Monday | Tuesday | Wednesday | Thursday | Friday | Saturday | Sunday\n\
        ---------+--------+---------+-----------+----------+--------+----------+-------\n\
        Ana      | M A E  | -       | -         | -        | -      | -        | M A E\n\
        Bo       | -      | E M A   | -         | -        | -      | -        | -\n";
    assert!(text.starts_with(table), "{text}");
}
