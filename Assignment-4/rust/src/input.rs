//! The start menu, the typed input, the ranking parser, the slot maximum variable, and the sample employees.

use std::io::{self, BufRead, Write};

use crate::model::{
    DEFAULT_MAXIMUM_EMPLOYEES_PER_SLOT, Day, Employee, MAXIMUM_DAYS_PER_WEEK, MINIMUM_EMPLOYEES_PER_SLOT, Preference,
    Shift,
};

pub const SLOT_MAXIMUM_VARIABLE: &str = "MAX_EMPLOYEES_PER_SLOT";

const RANKING_HELP: &str = "Use 1 to 3 different letters from M, A, and E, or - for a day off.";

/// Entry order is table order. `-` is a day off.
const SAMPLE_EMPLOYEES: [(&str, [&str; 7]); 13] = [
    ("Ana", ["MAE", "MAE", "MAE", "-", "-", "AME", "EMA"]),
    ("Ben", ["MAE", "MAE", "-", "MAE", "-", "EMA", "MAE"]),
    ("Cara", ["MEA", "MAE", "-", "-", "MAE", "AME", "EMA"]),
    ("Dan", ["MAE", "-", "MAE", "MAE", "-", "MAE", "AME"]),
    ("Eve", ["AME", "-", "MAE", "-", "MAE", "MAE", "EMA"]),
    ("Finn", ["AEM", "-", "AME", "MAE", "-", "EMA", "MAE"]),
    ("Gia", ["EAM", "-", "AME", "-", "AME", "MAE", "AME"]),
    ("Hugo", ["AEM", "-", "-", "AME", "AME", "AME", "MAE"]),
    ("Ivy", ["EMA", "-", "-", "AME", "EMA", "EMA", "AME"]),
    ("Jay", ["-", "AME", "-", "-", "-", "-", "-"]),
    ("Kai", ["EAM", "-", "AME", "EMA", "EMA", "-", "-"]),
    ("Leo", ["MAE", "AME", "EMA", "EMA", "MAE", "-", "-"]),
    ("Mia", ["-", "AME", "-", "-", "-", "-", "-"]),
];

#[derive(Debug, PartialEq, Eq)]
pub struct InvalidRanking;

#[derive(Debug, PartialEq, Eq)]
pub struct InvalidSlotMaximum;

/// Returns the preference for one day, or `None` for a day off.
pub fn parse_ranking(text: &str) -> Result<Option<Preference>, InvalidRanking> {
    let compact: String = text.chars().filter(|character| !character.is_whitespace()).collect();
    if compact == "-" {
        return Ok(None);
    }

    let mut ranked = Vec::with_capacity(Shift::ALL.len());
    for letter in compact.chars() {
        let shift = match letter.to_ascii_uppercase() {
            'M' => Shift::Morning,
            'A' => Shift::Afternoon,
            'E' => Shift::Evening,
            _ => return Err(InvalidRanking),
        };
        if ranked.contains(&shift) {
            return Err(InvalidRanking);
        }
        ranked.push(shift);
    }
    if ranked.is_empty() {
        return Err(InvalidRanking);
    }

    for shift in Shift::ALL {
        if !ranked.contains(&shift) {
            ranked.push(shift);
        }
    }
    Ok(Some([ranked[0], ranked[1], ranked[2]]))
}

/// `None` is a variable that is not set.
pub fn parse_slot_maximum(value: Option<&str>) -> Result<usize, InvalidSlotMaximum> {
    match value {
        None => Ok(DEFAULT_MAXIMUM_EMPLOYEES_PER_SLOT),
        Some(text) => text
            .parse::<usize>()
            .ok()
            .filter(|&maximum| maximum >= MINIMUM_EMPLOYEES_PER_SLOT)
            .ok_or(InvalidSlotMaximum),
    }
}

pub fn read_slot_maximum() -> usize {
    let value = std::env::var(SLOT_MAXIMUM_VARIABLE).ok();
    parse_slot_maximum(value.as_deref()).unwrap_or_else(|InvalidSlotMaximum| {
        eprintln!(
            "{SLOT_MAXIMUM_VARIABLE} must be a whole number of {MINIMUM_EMPLOYEES_PER_SLOT} or more. \
             Using {DEFAULT_MAXIMUM_EMPLOYEES_PER_SLOT}."
        );
        DEFAULT_MAXIMUM_EMPLOYEES_PER_SLOT
    })
}

pub fn sample_employees() -> Vec<Employee> {
    SAMPLE_EMPLOYEES
        .iter()
        .map(|(name, rankings)| Employee {
            name: name.to_string(),
            preferences: rankings.map(|ranking| parse_ranking(ranking).expect("the sample rankings are valid")),
        })
        .collect()
}

/// Runs the start menu and the typed input. An empty list means that no employee was entered.
pub fn collect_employees(input: &mut impl BufRead, output: &mut impl Write) -> io::Result<Vec<Employee>> {
    writeln!(output, "Employee Scheduler")?;
    writeln!(output)?;
    writeln!(output, "1) Load the sample employees")?;
    writeln!(output, "2) Start with no employees")?;

    loop {
        let Some(answer) = prompt(input, output, "Choose 1 or 2: ")? else {
            return Ok(Vec::new());
        };
        match answer.trim() {
            "1" => return load_sample_then_ask(input, output),
            "2" => return read_typed_employees(input, output, Vec::new()),
            _ => writeln!(output, "Enter 1 or 2.")?,
        }
    }
}

fn load_sample_then_ask(input: &mut impl BufRead, output: &mut impl Write) -> io::Result<Vec<Employee>> {
    let employees = sample_employees();
    writeln!(output, "Loaded {} sample employees.", employees.len())?;

    loop {
        let Some(answer) = prompt(input, output, "Add more employees? (y/n): ")? else {
            return Ok(employees);
        };
        match answer.trim().to_ascii_lowercase().as_str() {
            "y" => return read_typed_employees(input, output, employees),
            "n" => return Ok(employees),
            _ => writeln!(output, "Enter y or n.")?,
        }
    }
}

fn read_typed_employees(
    input: &mut impl BufRead,
    output: &mut impl Write,
    mut employees: Vec<Employee>,
) -> io::Result<Vec<Employee>> {
    loop {
        let Some(line) = prompt(input, output, "Employee name (blank to finish): ")? else {
            return Ok(employees);
        };
        let name = line.trim();
        if name.is_empty() {
            return Ok(employees);
        }

        let lower_case_name = name.to_lowercase();
        if employees.iter().any(|employee| employee.name.to_lowercase() == lower_case_name) {
            writeln!(output, "Name already exists: {name}.")?;
            continue;
        }

        // End of input partway through the days drops this unfinished employee.
        let Some(preferences) = read_preferences(input, output, name)? else {
            return Ok(employees);
        };
        employees.push(Employee { name: name.to_string(), preferences });
    }
}

/// Returns `None` at end of input.
fn read_preferences(
    input: &mut impl BufRead,
    output: &mut impl Write,
    name: &str,
) -> io::Result<Option<[Option<Preference>; 7]>> {
    let mut preferences = [None; 7];
    let mut working_days = 0;

    for day in Day::ALL {
        let label = format!("{} ranking (M A E, or - for a day off): ", day.name());
        loop {
            let Some(answer) = prompt(input, output, &label)? else {
                return Ok(None);
            };
            match parse_ranking(&answer) {
                Ok(preference) => {
                    preferences[day.index()] = preference;
                    if preference.is_some() {
                        working_days += 1;
                    }
                    break;
                }
                Err(InvalidRanking) => writeln!(output, "{RANKING_HELP}")?,
            }
        }

        if working_days == MAXIMUM_DAYS_PER_WEEK {
            writeln!(output, "{name} has {MAXIMUM_DAYS_PER_WEEK} working days. The remaining days are days off.")?;
            break;
        }
    }
    Ok(Some(preferences))
}

/// Returns the line without its line ending, or `None` at end of input.
fn prompt(input: &mut impl BufRead, output: &mut impl Write, label: &str) -> io::Result<Option<String>> {
    write!(output, "{label}")?;
    output.flush()?;

    let mut line = String::new();
    if input.read_line(&mut line)? == 0 {
        // Ends the prompt line, so the next output starts on a new line.
        writeln!(output)?;
        return Ok(None);
    }
    Ok(Some(line.trim_end_matches(['\n', '\r']).to_string()))
}
