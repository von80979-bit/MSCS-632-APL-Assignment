//! Reads the slot maximum, runs the input, the scheduler, and the output.

mod input;
mod model;
mod output;
mod scheduler;

fn main() {
    // Placeholder until the Rust implementation lands. It proves that the Docker run passes the variable through.
    let slot_maximum = std::env::var("MAX_EMPLOYEES_PER_SLOT").unwrap_or_else(|_| String::from("unset"));
    println!("scheduler-rust: not implemented yet. MAX_EMPLOYEES_PER_SLOT={slot_maximum}");
}
