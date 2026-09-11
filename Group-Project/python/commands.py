"""Prompting, validation, error messages, and output formatting.

One function per command. Each prompts for its fields, rejects bad input with the message the
specification fixes, and asks the same question again rather than abandoning the command.
"""

from typing import List, Optional
import model
import store


def format_table(expense_list: List[dict]) -> str:
    """Format a list of expense dicts into the standard table output."""
    header = f"{'ID':>4}  {'DATE':<10}  {'AMOUNT':>9}  {'CATEGORY':<15}  {'DESCRIPTION'}"
    rows = [header]
    for exp in expense_list:
        amount_str = f"${exp['amount']:.2f}"
        row = (
            f"{exp['id']:>4}  "
            f"{str(exp['date']):<10}  "
            f"{amount_str:>9}  "
            f"{exp['category']:<15}  "
            f"{exp['description']}"
        )
        rows.append(row)
    return "\n".join(rows)


def format_summary(cat_totals: List[tuple], overall: float) -> str:
    """Format category totals and overall total into the summary block."""
    lines = []
    for cat, total in cat_totals:
        amount_str = f"${total:.2f}"
        lines.append(f"{cat:<15}{amount_str:>10}")
    lines.append("")
    overall_str = f"${overall:.2f}"
    lines.append(f"{'Overall':<15}{overall_str:>10}")
    return "\n".join(lines)


def add() -> None:
    """Prompt for fields, validate each with re-prompting, and record a new expense."""
    while True:
        raw_date = input("Date (YYYY-MM-DD): ").strip()
        try:
            model.parse_date(raw_date)
            break
        except ValueError as exc:
            print(exc)

    while True:
        raw_amount = input("Amount: ").strip()
        try:
            model.parse_amount(raw_amount)
            break
        except ValueError as exc:
            print(exc)

    while True:
        raw_category = input("Category (Food, Transport, Utilities, Entertainment, Other): ").strip()
        try:
            model.parse_category(raw_category)
            break
        except ValueError as exc:
            print(exc)

    raw_description = input("Description: ").strip()

    new_id = store.add_expense(raw_date, raw_amount, raw_category, raw_description)
    print(f"Added expense {new_id}.")


def list_expenses() -> None:
    """Print all expenses sorted by id ascending, or 'No expenses.' when empty."""
    items = store.list_expenses()
    if not items:
        print("No expenses.")
    else:
        print(format_table(items))


def filter_expenses() -> None:
    """Prompt for date range and category, then display matching expenses."""
    start_date = None
    while True:
        raw = input("Start date (YYYY-MM-DD, blank for none): ").strip()
        if not raw:
            break
        try:
            start_date = model.parse_date(raw)
            break
        except ValueError as exc:
            print(exc)

    end_date = None
    while True:
        raw = input("End date (YYYY-MM-DD, blank for none): ").strip()
        if not raw:
            break
        try:
            end_date = model.parse_date(raw)
            break
        except ValueError as exc:
            print(exc)

    filter_cat = None
    while True:
        raw = input("Category (blank for all): ").strip()
        if not raw:
            break
        try:
            filter_cat = model.parse_category(raw)
            break
        except ValueError as exc:
            print(exc)

    matches = store.filter_expenses(start_date=start_date, end_date=end_date, category=filter_cat)
    if not matches:
        print("No expenses match.")
    else:
        print(format_table(matches))


def summary() -> None:
    """Print category totals and overall total."""
    cat_totals, overall = store.get_summary()
    print(format_summary(cat_totals, overall))


def edit() -> None:
    """Prompt for id and each editable field, keeping existing values on blank."""
    raw_id = input("Expense id: ").strip()
    try:
        exp_id = int(raw_id)
    except ValueError:
        print(f"No expense with id {raw_id}.")
        return

    expense = store.get_expense(exp_id)
    if expense is None:
        print(f"No expense with id {raw_id}.")
        return

    cur_date = str(expense["date"])
    cur_amount = f"{expense['amount']:.2f}"
    cur_cat = expense["category"]
    cur_desc = expense["description"]

    while True:
        raw = input(f"Date [{cur_date}] (blank to keep): ").strip()
        if not raw:
            new_date = cur_date
            break
        try:
            model.parse_date(raw)
            new_date = raw
            break
        except ValueError as exc:
            print(exc)

    while True:
        raw = input(f"Amount [{cur_amount}] (blank to keep): ").strip()
        if not raw:
            new_amount = cur_amount
            break
        try:
            model.parse_amount(raw)
            new_amount = raw
            break
        except ValueError as exc:
            print(exc)

    while True:
        raw = input(f"Category [{cur_cat}] (blank to keep): ").strip()
        if not raw:
            new_cat = cur_cat
            break
        try:
            model.parse_category(raw)
            new_cat = raw
            break
        except ValueError as exc:
            print(exc)

    raw = input(f"Description [{cur_desc}] (blank to keep): ").strip()
    if not raw:
        new_desc = cur_desc
    else:
        new_desc = raw

    store.update_expense(exp_id, new_date, new_amount, new_cat, new_desc)
    print(f"Updated expense {exp_id}.")


def delete() -> None:
    """Prompt for id and remove that expense, or print not-found message."""
    raw_id = input("Expense id: ").strip()
    try:
        exp_id = int(raw_id)
    except ValueError:
        print(f"No expense with id {raw_id}.")
        return

    if not store.delete_expense(exp_id):
        print(f"No expense with id {raw_id}.")
        return

    print(f"Deleted expense {exp_id}.")
