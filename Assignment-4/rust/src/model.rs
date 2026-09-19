//! Days, shifts, employees, preferences, slots, the schedule, and the change types.

pub const MINIMUM_EMPLOYEES_PER_SLOT: usize = 2;
pub const DEFAULT_MAXIMUM_EMPLOYEES_PER_SLOT: usize = 3;
pub const MAXIMUM_DAYS_PER_WEEK: usize = 5;

#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub enum Day {
    Monday,
    Tuesday,
    Wednesday,
    Thursday,
    Friday,
    Saturday,
    Sunday,
}

impl Day {
    pub const ALL: [Day; 7] =
        [Day::Monday, Day::Tuesday, Day::Wednesday, Day::Thursday, Day::Friday, Day::Saturday, Day::Sunday];

    pub fn index(self) -> usize {
        self as usize
    }

    pub fn name(self) -> &'static str {
        match self {
            Day::Monday => "Monday",
            Day::Tuesday => "Tuesday",
            Day::Wednesday => "Wednesday",
            Day::Thursday => "Thursday",
            Day::Friday => "Friday",
            Day::Saturday => "Saturday",
            Day::Sunday => "Sunday",
        }
    }
}

#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub enum Shift {
    Morning,
    Afternoon,
    Evening,
}

impl Shift {
    pub const ALL: [Shift; 3] = [Shift::Morning, Shift::Afternoon, Shift::Evening];

    pub fn index(self) -> usize {
        self as usize
    }

    /// The lower-case name, as a change line uses it after a day: "Monday morning".
    pub fn name(self) -> &'static str {
        match self {
            Shift::Morning => "morning",
            Shift::Afternoon => "afternoon",
            Shift::Evening => "evening",
        }
    }

    /// The capitalized name, as the schedule table header uses it.
    pub fn title(self) -> String {
        let name = self.name();
        name[..1].to_uppercase() + &name[1..]
    }

    /// The letter of the ranking input, as the preference table uses it.
    pub fn letter(self) -> String {
        self.name()[..1].to_uppercase()
    }
}

/// The three shifts of one day in rank order, first choice first.
pub type Preference = [Shift; 3];

#[derive(Clone, Debug, PartialEq, Eq)]
pub struct Employee {
    pub name: String,
    /// One entry per day from Monday to Sunday. `None` is a day off.
    pub preferences: [Option<Preference>; 7],
}

impl Employee {
    pub fn preference_on(&self, day: Day) -> Option<Preference> {
        self.preferences[day.index()]
    }
}

/// An employee is referred to by their position in entry order.
pub type EmployeeId = usize;

#[derive(Clone, Debug, PartialEq, Eq)]
pub struct Slot {
    pub day: Day,
    pub shift: Shift,
    /// The assigned employees in order of assignment.
    pub employees: Vec<EmployeeId>,
}

impl Slot {
    pub fn is_short(&self) -> bool {
        self.employees.len() < MINIMUM_EMPLOYEES_PER_SLOT
    }

    pub fn has_room(&self, slot_maximum: usize) -> bool {
        self.employees.len() < slot_maximum
    }

    /// A donor keeps at least the minimum after it gives one employee away.
    pub fn is_donor(&self) -> bool {
        self.employees.len() > MINIMUM_EMPLOYEES_PER_SLOT
    }
}

#[derive(Clone, Debug, PartialEq, Eq)]
pub enum Change {
    RankFallback { employee: EmployeeId, day: Day, first_choice: Shift, placed: Shift, choice: usize },
    DayMove { employee: EmployeeId, ranked_day: Day, day: Day, shift: Shift },
    Unplaced { employee: EmployeeId, day: Day },
    RandomFill { employee: EmployeeId, day: Day, shift: Shift },
    Transfer { employee: EmployeeId, from_day: Day, from_shift: Shift, to_day: Day, to_shift: Shift },
    StillShort { day: Day, shift: Shift, employee_count: usize },
}

#[derive(Clone, Debug, PartialEq, Eq)]
pub struct Schedule {
    /// The 21 slots, Monday morning first and Sunday evening last.
    pub slots: Vec<Slot>,
    pub changes: Vec<Change>,
}

impl Schedule {
    pub fn empty() -> Self {
        let slots = Day::ALL
            .iter()
            .flat_map(|&day| Shift::ALL.iter().map(move |&shift| Slot { day, shift, employees: Vec::new() }))
            .collect();
        Schedule { slots, changes: Vec::new() }
    }

    pub fn slot(&self, day: Day, shift: Shift) -> &Slot {
        &self.slots[day.index() * Shift::ALL.len() + shift.index()]
    }

    pub fn slot_mut(&mut self, day: Day, shift: Shift) -> &mut Slot {
        &mut self.slots[day.index() * Shift::ALL.len() + shift.index()]
    }

    pub fn is_working_on(&self, employee: EmployeeId, day: Day) -> bool {
        Shift::ALL.iter().any(|&shift| self.slot(day, shift).employees.contains(&employee))
    }

    pub fn days_worked(&self, employee: EmployeeId) -> usize {
        Day::ALL.iter().filter(|&&day| self.is_working_on(employee, day)).count()
    }
}
