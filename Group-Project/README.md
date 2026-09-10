# Expense Tracker

MSCS-632 group project, deliverable 1. One terminal application, implemented in Python and C++.

The application records expenses, lists them, filters them by date range and category, and totals them by category and overall. 
Both implementations offer the same feature sets over the same record structure. The data lives in memory for the length of one run.

The purpose of building the same application in 2 different languages is to compare how the two languages handle the same job, in data structures, memory management, and error handling. 
The comparison is written up separately in the project report.

## Running it

**Python** needs Python 3.8 or later.

```sh
cd python
python3 main.py
```

**C++** needs a compiler with C++17 support.

```sh
cd cpp
make
./expense-tracker
```

`make clean` removes the binary.

## Commands

| Command | Does |
| --- | --- |
| `add` | prompts for date, amount, category and description, then adds an expense |
| `list` | prints every expense, sorted by id |
| `filter` | prompts for a start date, an end date and a category, any of which may be left blank |
| `summary` | prints a total for each category and an overall total |
| `edit` | prompts for an id, then for each field, where blank keeps the current value |
| `delete` | prompts for an id and removes that expense |
| `quit` | ends the program |


## Layout and ownership

| Directory | Holds | Owned by |
| --- | --- | --- |
| `python/` | the Python implementation | Sabin Ranabhat |
| `cpp/` | the C++ implementation | Nguyen Vo |


Both implementations use the same four units, which is what lets the report put the two implementations side by side.

| Unit | Holds |
| --- | --- |
| `model` | the expense record, the category, the date, and conversion between text and those types |
| `store` | the collection, the category index, and the operations over them |
| `commands` | prompting, validation, error messages, and output formatting |
| `main` | the command loop |
