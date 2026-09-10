"""The collection of expenses, the category index, and the operations over them.

Expenses are held in a dict keyed by id. A second dict maps a category name to the list of
expenses in it, and summary and category filtering read that index rather than scanning
everything. Nothing is written to disk; the data lives for one run.
"""
