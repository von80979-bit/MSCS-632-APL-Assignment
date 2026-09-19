//! The preference table, the schedule table, the changes, and the employee summary.

use crate::model::{Change, Day, Employee, EmployeeId, MINIMUM_EMPLOYEES_PER_SLOT, Schedule, Shift, Slot};

const SUMMARY_GAP: &str = "    ";

pub fn render(schedule: &Schedule, employees: &[Employee], slot_maximum: usize) -> String {
    let mut lines = vec![String::from("Employee preferences"), String::new()];
    lines.extend(grid(&preference_rows(employees)));
    lines.push(String::new());
    lines.push(format!("Weekly schedule (maximum {slot_maximum} employees per slot)"));
    lines.push(String::new());
    lines.extend(grid(&schedule_rows(schedule, employees)));
    lines.push(String::new());
    lines.push(String::from("Changes"));
    if schedule.changes.is_empty() {
        lines.push(String::from("No changes."));
    }
    lines.extend(schedule.changes.iter().map(|change| change_line(change, employees)));
    lines.push(String::new());
    lines.push(String::from("Employees"));
    lines.extend(employee_summary(schedule, employees));

    let mut text = lines.join("\n");
    text.push('\n');
    text
}

/// The header and one row per employee in entry order, with the ranking of each day.
fn preference_rows(employees: &[Employee]) -> Vec<Vec<String>> {
    let header: Vec<String> =
        std::iter::once(String::from("Employee")).chain(Day::ALL.iter().map(|day| day.name().to_string())).collect();
    let rows = employees.iter().map(|employee| {
        std::iter::once(employee.name.clone()).chain(Day::ALL.iter().map(|&day| ranking_cell(employee, day))).collect()
    });
    std::iter::once(header).chain(rows).collect()
}

/// The shift letters in rank order, for example "M A E", or "-" for a day off.
fn ranking_cell(employee: &Employee, day: Day) -> String {
    match employee.preference_on(day) {
        None => String::from("-"),
        Some(preference) => preference.iter().map(|shift| shift.letter()).collect::<Vec<_>>().join(" "),
    }
}

/// The header and one row per day, with the employees of each slot.
fn schedule_rows(schedule: &Schedule, employees: &[Employee]) -> Vec<Vec<String>> {
    let header: Vec<String> =
        std::iter::once(String::from("Day")).chain(Shift::ALL.iter().map(|shift| shift.title())).collect();
    let rows = Day::ALL.iter().map(|&day| {
        std::iter::once(day.name().to_string())
            .chain(Shift::ALL.iter().map(|&shift| slot_cell(schedule.slot(day, shift), employees)))
            .collect()
    });
    std::iter::once(header).chain(rows).collect()
}

/// Renders the rows as a table. The first row is the header, and each column is as wide as its longest cell.
fn grid(rows: &[Vec<String>]) -> Vec<String> {
    let (header, body) = rows.split_first().expect("a table has a header row");

    let widths: Vec<usize> =
        (0..header.len()).map(|column| rows.iter().map(|row| row[column].chars().count()).max().unwrap_or(0)).collect();

    let last_column = widths.len() - 1;
    let separator = widths
        .iter()
        .enumerate()
        .map(|(column, &width)| {
            let spaces_around_bar = if column == 0 || column == last_column { 1 } else { 2 };
            "-".repeat(width + spaces_around_bar)
        })
        .collect::<Vec<_>>()
        .join("+");

    let table_row = |row: &Vec<String>| {
        row.iter()
            .zip(&widths)
            .map(|(cell, &width)| format!("{cell:<width$}"))
            .collect::<Vec<_>>()
            .join(" | ")
            .trim_end()
            .to_string()
    };

    let mut lines = vec![table_row(header), separator];
    lines.extend(body.iter().map(table_row));
    lines
}

fn slot_cell(slot: &Slot, employees: &[Employee]) -> String {
    match slot.employees.as_slice() {
        [] => String::from("(empty)"),
        [only] => format!("{} (short)", employees[*only].name),
        assigned => assigned.iter().map(|&employee| employees[employee].name.as_str()).collect::<Vec<_>>().join(", "),
    }
}

fn change_line(change: &Change, employees: &[Employee]) -> String {
    let name_of = |employee: EmployeeId| employees[employee].name.as_str();
    match *change {
        Change::RankFallback { employee, day, first_choice, placed, choice } => format!(
            "{}: {} {} was full. Placed in {} {} (choice {choice}).",
            name_of(employee),
            day.name(),
            first_choice.name(),
            day.name(),
            placed.name()
        ),
        Change::DayMove { employee, ranked_day, day, shift } => {
            format!("{}: {} was full. Moved to {} {}.", name_of(employee), ranked_day.name(), day.name(), shift.name())
        }
        Change::Unplaced { employee, day } => format!(
            "{}: {} was full and no other day had room. Not placed for {}.",
            name_of(employee),
            day.name(),
            day.name()
        ),
        Change::RandomFill { employee, day, shift } => {
            format!("{}: {} {} was short. Added at random.", name_of(employee), day.name(), shift.name())
        }
        Change::Transfer { employee, from_day, from_shift, to_day, to_shift } => format!(
            "{}: Moved from {} {} to {} {} to cover a short slot.",
            name_of(employee),
            from_day.name(),
            from_shift.name(),
            to_day.name(),
            to_shift.name()
        ),
        Change::StillShort { day, shift, employee_count } => format!(
            "{} {} is still short: {employee_count} of {MINIMUM_EMPLOYEES_PER_SLOT} employees.",
            day.name(),
            shift.name()
        ),
    }
}

fn employee_summary(schedule: &Schedule, employees: &[Employee]) -> Vec<String> {
    let name_width = employees.iter().map(|employee| employee.name.chars().count()).max().unwrap_or(0);
    employees
        .iter()
        .enumerate()
        .map(|(id, employee)| {
            let days = schedule.days_worked(id);
            let unit = if days == 1 { "day" } else { "days" };
            format!("{:<name_width$}{SUMMARY_GAP}{days} {unit}", employee.name)
        })
        .collect()
}
