"""The collection of expenses, the category index, and the operations over them.

Expenses are held in a dict keyed by id. A second dict maps a category name to the list of
expenses in it, and summary and category filtering read that index rather than scanning
everything. Nothing is written to disk; the data lives for one run.
"""

from typing import Dict, List, Optional, Tuple
import datetime

from model import CATEGORIES, parse_expense_fields

# Dict keyed by id: int -> dict
expenses: Dict[int, dict] = {}

# Secondary index: category name -> list of expenses
category_index: Dict[str, List[dict]] = {cat: [] for cat in CATEGORIES}

# Auto-incrementing ID counter (never reused)
_next_id: int = 1


def seed_data() -> None:
    """Seed the in-memory store with the exactly specified twelve initial expenses."""
    global _next_id
    expenses.clear()
    for cat in CATEGORIES:
        category_index[cat].clear()

    seeds = [
        (1, "2026-01-05", "12.50", "Food", "Coffee and pastry"),
        (2, "2026-01-08", "45.00", "Transport", "Monthly transit top-up"),
        (3, "2026-01-14", "120.75", "Utilities", "Electricity bill"),
        (4, "2026-01-22", "18.99", "Entertainment", "Cinema ticket"),
        (5, "2026-01-28", "32.40", "Food", "Grocery run"),
        (6, "2026-02-03", "9.25", "Other", "Stationery"),
        (7, "2026-02-11", "62.10", "Food", "Dinner with friends"),
        (8, "2026-02-17", "28.00", "Transport", "Airport taxi"),
        (9, "2026-02-25", "95.40", "Utilities", "Internet and phone"),
        (10, "2026-03-04", "15.00", "Entertainment", "Streaming subscription"),
        (11, "2026-03-12", "41.85", "Food", "Weekly groceries"),
        (12, "2026-03-20", "22.30", "Other", "Pharmacy"),
    ]

    for exp_id, d, a, c, desc in seeds:
        record = parse_expense_fields(d, a, c, desc, expense_id=exp_id)
        expenses[exp_id] = record
        category_index[record["category"]].append(record)

    _next_id = 13


def add_expense(date_str: str, amount_str: str, category_str: str, description_str: str) -> int:
    """Add a new expense using the next available id and update category index."""
    global _next_id
    new_id = _next_id
    _next_id += 1
    record = parse_expense_fields(date_str, amount_str, category_str, description_str, expense_id=new_id)
    expenses[new_id] = record
    category_index[record["category"]].append(record)
    return new_id


def get_expense(expense_id: int) -> Optional[dict]:
    """Retrieve an expense by id, or None if not found."""
    return expenses.get(expense_id)


def list_expenses() -> List[dict]:
    """Return every expense sorted by id ascending."""
    return sorted(expenses.values(), key=lambda e: e["id"])


def filter_expenses(
    start_date: Optional[datetime.date] = None,
    end_date: Optional[datetime.date] = None,
    category: Optional[str] = None,
) -> List[dict]:
    """Filter expenses by date range and category.

    When category is given, reads category_index directly rather than scanning everything.
    """
    if category is not None:
        candidates = category_index.get(category, [])
    else:
        candidates = list(expenses.values())

    results = []
    for exp in candidates:
        if start_date is not None and exp["date"] < start_date:
            continue
        if end_date is not None and exp["date"] > end_date:
            continue
        results.append(exp)

    return sorted(results, key=lambda e: e["id"])


def get_summary() -> Tuple[List[Tuple[str, float]], float]:
    """Calculate totals for each category in fixed order and an overall total.

    Reads category_index rather than scanning everything.
    """
    cat_totals: List[Tuple[str, float]] = []
    overall = 0.0
    for cat in CATEGORIES:
        total = sum(exp["amount"] for exp in category_index[cat])
        cat_totals.append((cat, total))
        overall += total
    return cat_totals, overall


def update_expense(
    expense_id: int,
    date_str: str,
    amount_str: str,
    category_str: str,
    description_str: str,
) -> bool:
    """Update an existing expense in place, maintaining category index and preserving id."""
    if expense_id not in expenses:
        return False
    old_record = expenses[expense_id]
    old_cat = old_record["category"]

    updated = parse_expense_fields(date_str, amount_str, category_str, description_str, expense_id=expense_id)
    new_cat = updated["category"]

    if old_cat != new_cat:
        if old_record in category_index[old_cat]:
            category_index[old_cat].remove(old_record)
        category_index[new_cat].append(updated)
    else:
        idx = category_index[old_cat].index(old_record)
        category_index[old_cat][idx] = updated

    expenses[expense_id] = updated
    return True


def delete_expense(expense_id: int) -> bool:
    """Remove an expense by id from the primary dict and category index."""
    if expense_id not in expenses:
        return False
    record = expenses.pop(expense_id)
    if record in category_index[record["category"]]:
        category_index[record["category"]].remove(record)
    return True
