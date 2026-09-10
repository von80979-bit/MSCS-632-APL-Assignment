"""The expense record, the category, the date, and conversion between text and those types.

An expense is a dict with the keys id, date, amount, category and description. Categories are the
five names the specification fixes. Date work goes through datetime, and every conversion from
typed text happens here rather than at the prompt.
"""
