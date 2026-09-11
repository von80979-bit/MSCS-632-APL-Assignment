"""The expense record, the category, the date, and conversion between text and those types.

An expense is a dict with the keys id, date, amount, category and description. Categories are the
five names the specification fixes. Date work goes through datetime, and every conversion from
typed text happens here rather than at the prompt.
"""

import datetime
import re

CATEGORIES = [
    "Food",
    "Transport",
    "Utilities",
    "Entertainment",
    "Other",
]

_CATEGORY_LOOKUP = {cat.lower(): cat for cat in CATEGORIES}


def parse_date(text: str) -> datetime.date:
    """Parse and validate a date string in YYYY-MM-DD form.

    Distinguishes shape errors (Case 3) from non-existent calendar dates (Case 4).
    """
    if not re.match(r"^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$", text):
        raise ValueError("Date must be in YYYY-MM-DD form.")
    try:
        return datetime.datetime.strptime(text, "%Y-%m-%d").date()
    except ValueError:
        raise ValueError("That date does not exist.")


def parse_amount(text: str) -> float:
    """Parse and validate an amount string as a positive float."""
    try:
        amount = float(text)
    except ValueError:
        raise ValueError("Amount must be a number.")
    if amount <= 0:
        raise ValueError("Amount must be greater than zero.")
    return amount


def parse_category(text: str) -> str:
    """Parse and validate a category string against the 5 allowed categories (case-insensitive)."""
    cleaned = text.strip().lower()
    if cleaned not in _CATEGORY_LOOKUP:
        raise ValueError("Category must be one of: Food, Transport, Utilities, Entertainment, Other.")
    return _CATEGORY_LOOKUP[cleaned]


def parse_expense_fields(
    date_str: str,
    amount_str: str,
    category_str: str,
    description_str: str,
    expense_id: int = 0,
) -> dict:
    """Convert typed text fields into an expense dict with no declared class shape.

    Takes the four typed strings and returns the dict. The single place the report
    points at for dynamic typing.
    """
    parsed_date = parse_date(date_str)
    parsed_amount = parse_amount(amount_str)
    parsed_category = parse_category(category_str)
    return {
        "id": expense_id,
        "date": parsed_date,
        "amount": parsed_amount,
        "category": parsed_category,
        "description": description_str,
    }
