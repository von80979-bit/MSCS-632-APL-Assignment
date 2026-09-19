# Assignment 4: Employee scheduler

Two weekly employee schedule programs built in Rust and Go. 

## What the application does

The company works 7 days a week, and each day has three shifts: morning, afternoon, and evening. One shift on one day is a slot, so the week has 21 slots.

The application collects employee names and their shift preferences, builds the schedule for one week, and prints four parts: 
- The preference table, which shows the input the schedule comes from
- The schedule table
- The changes the scheduler made
- The days worked by each employee. 

The program exits after it prints the schedule.

Each employee ranks the three shifts for each day they want to work.

## Scheduling rules

The scheduler enforces three rules:

- An employee works at most one shift per day.
- An employee works at most 5 days per week.
- Each slot holds at least 2 employees.

The scheduler runs four checks over the week, in this order:

1. **Preferences.** Each employee goes to their first choice on each day they ranked. If that slot is full, they go to the second choice, then the third. A placement in a lower choice is a rank fallback.
2. **Day moves.** If every slot on the ranked day is full, the employee moves to another day with room. The search starts at the next day and goes to Sunday, and then it starts again at Monday. For example, an overflow on Thursday searches Friday, Saturday, Sunday, Monday, Tuesday, and Wednesday. An employee with no free day stays unplaced.
3. **Random fill.** Each slot with fewer than 2 employees gets random employees who are free on that day and who work fewer than 5 days.
4. **Transfer.** Each slot that is still short takes an employee from a slot with more than 2 employees. A donor slot on the same day is used first.

## Assumptions

The requirements are unclear on certain points, so the two programs share these assumptions:

1. A slot holds at most 3 employees, because the requirement gives a minimum of 2 but no maximum. The environment variable `MAX_EMPLOYEES_PER_SLOT` changes this maximum.
2. The scheduler places employees first come, first served, in the order of entry.
3. The name identifies an employee, so two employees must not have the same name. The comparison ignores letter case and outer spaces.
4. Weekends are working days, so an employee chooses any 5 days of the week.
5. A day without a preference means "no preference", not "unavailable". The random fill and the transfer can use such a day.
6. The schedule covers one week, so a Sunday overflow goes to another day of the same week.

## Run with Docker

Docker is the primary run path. One image holds both programs. Run these commands in this folder.

Build the image:

```sh
docker build -t employee-scheduler .
```

Run the Rust program:

```sh
docker run -it --rm employee-scheduler scheduler-rust
```

Run the Go program:

```sh
docker run -it --rm employee-scheduler scheduler-go
```

Change the slot maximum with the environment variable:

```sh
docker run -it --rm -e MAX_EMPLOYEES_PER_SLOT=4 employee-scheduler scheduler-rust
```

If the value is not a whole number of 2 or more, the program prints a warning to standard error and uses 3.

## Run without Docker

The Rust program needs Rust 1.85 or later, because it uses edition 2024. The Go program needs Go 1.24 or later.

Run the Rust program:

```sh
cd rust
cargo run --release
```

Run the Go program:

```sh
cd go
go run .
```

## Run the tests

Each language has the same 13 test cases for the scheduling rules, the input parsers, and the preference table.

```sh
cd rust
cargo test
```

```sh
cd go
go test ./...
```

## Use the application

At the start, the program offers two paths:

```
Employee Scheduler

1) Load the sample employees
2) Start with no employees
Choose 1 or 2:
```

Option 1 loads 13 sample employees, and then asks whether you want to add more. Option 2 starts with no employees.

For each new employee, the program asks for a name, and then for one ranking per day:

```
Employee name (blank to finish): Ana
Monday ranking (M A E, or - for a day off): mae
```

- Give 1 to 3 different letters from `M`, `A`, and `E`, in the order of your choice. Spaces and letter case do not count.
- A preference always ranks all three shifts. If you give fewer than three letters, the program adds the shifts you left out at the end, in the order morning, afternoon, evening. For example, `e` gives evening, then morning, then afternoon.
- Enter `-` for a day off.
- After the fifth working day, the program makes the remaining days days off.
- A blank name, or the end of input with Ctrl-D, ends the input and starts the scheduler.

## Sample output

This run uses the sample employees and the default maximum of 3:

```
Employee preferences

Employee | Monday | Tuesday | Wednesday | Thursday | Friday | Saturday | Sunday
---------+--------+---------+-----------+----------+--------+----------+-------
Ana      | M A E  | M A E   | M A E     | -        | -      | A M E    | E M A
Ben      | M A E  | M A E   | -         | M A E    | -      | E M A    | M A E
Cara     | M E A  | M A E   | -         | -        | M A E  | A M E    | E M A
Dan      | M A E  | -       | M A E     | M A E    | -      | M A E    | A M E
Eve      | A M E  | -       | M A E     | -        | M A E  | M A E    | E M A
Finn     | A E M  | -       | A M E     | M A E    | -      | E M A    | M A E
Gia      | E A M  | -       | A M E     | -        | A M E  | M A E    | A M E
Hugo     | A E M  | -       | -         | A M E    | A M E  | A M E    | M A E
Ivy      | E M A  | -       | -         | A M E    | E M A  | E M A    | A M E
Jay      | -      | A M E   | -         | -        | -      | -        | -
Kai      | E A M  | -       | A M E     | E M A    | E M A  | -        | -
Leo      | M A E  | A M E   | E M A     | E M A    | M A E  | -        | -
Mia      | -      | A M E   | -         | -        | -      | -        | -

Weekly schedule (maximum 3 employees per slot)

Day       | Morning         | Afternoon       | Evening
----------+-----------------+-----------------+---------------
Monday    | Ana, Ben, Cara  | Dan, Eve, Finn  | Gia, Hugo, Ivy
Tuesday   | Ben, Cara       | Jay, Leo, Mia   | Kai, Ana
Wednesday | Ana, Dan, Eve   | Finn, Gia, Kai  | Leo, Mia
Thursday  | Ben, Dan, Finn  | Hugo, Ivy       | Kai, Leo
Friday    | Cara, Eve, Leo  | Gia, Hugo       | Ivy, Kai
Saturday  | Dan, Eve, Gia   | Ana, Cara, Hugo | Ben, Finn, Ivy
Sunday    | Ben, Finn, Hugo | Dan, Gia, Ivy   | Ana, Cara, Eve

Changes
Dan: Monday morning was full. Placed in Monday afternoon (choice 2).
Hugo: Monday afternoon was full. Placed in Monday evening (choice 2).
Kai: Monday was full. Moved to Tuesday evening.
Leo: Monday was full and no other day had room. Not placed for Monday.
Mia: Wednesday evening was short. Added at random.
Ana: Moved from Tuesday morning to Tuesday evening to cover a short slot.

Employees
Ana     5 days
Ben     5 days
Cara    5 days
Dan     5 days
Eve     5 days
Finn    5 days
Gia     5 days
Hugo    5 days
Ivy     5 days
Jay     1 day
Kai     4 days
Leo     4 days
Mia     2 days
```

The preference table repeats the input, so a reader can check each assignment against the ranking that asked for it. A day off shows `-`.

A slot with one employee shows the name and `(short)`. An empty slot shows `(empty)`.

The sample always gives two rank fallbacks, one day move, one unplaced employee, one random fill, and one transfer. The random picks change between runs, so the two programs can name different employees in the last two lines of the changes.

## Project layout

```
Assignment-4/
  Dockerfile        both programs in one image
  README.md
  rust/
    src/
      main.rs       reads the slot maximum, and runs the input, the scheduler, and the output
      model.rs      days, shifts, employees, preferences, slots, the schedule, and the change types
      input.rs      the start menu, the typed input, the ranking parser, and the sample employees
      scheduler.rs  the four checks
      output.rs     the schedule table, the changes, and the employee summary
      scheduler/
        tests.rs    the 13 test cases
  go/
    main.go         the same five parts as the Rust program, one file each
    model.go
    input.go
    scheduler.go
    output.go
    scheduler_test.go
```

The Rust program uses the `rand` crate for the random picks. The Go program uses `math/rand/v2` from the standard library. Each scheduler receives its random source as a parameter, so the tests give it a seed.
