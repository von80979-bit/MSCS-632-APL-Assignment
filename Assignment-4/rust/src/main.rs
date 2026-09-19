//! Reads the slot maximum, runs the input, the scheduler, and the output.

mod input;
mod model;
mod output;
mod scheduler;

use std::io::{self, Write};

fn main() -> io::Result<()> {
    let slot_maximum = input::read_slot_maximum();

    let mut stdout = io::stdout().lock();
    let employees = input::collect_employees(&mut io::stdin().lock(), &mut stdout)?;
    if employees.is_empty() {
        writeln!(stdout, "No employees entered.")?;
        return Ok(());
    }

    let schedule = scheduler::build_schedule(&employees, slot_maximum, &mut rand::rng());
    writeln!(stdout)?;
    write!(stdout, "{}", output::render(&schedule, &employees, slot_maximum))
}
