"""Section 2: three semantic differences between Python, JavaScript and C++, shown by what this program prints.

The three are name binding, closure capture and type checking, in that order.

Run: python3 scope-closures.py
"""

# Name binding. Python resolves a global name on every call, not when the function is defined.


def report_measurement_limit():
    # measurement_limit does not exist yet. Python accepts the function anyway and looks the name up on each call.
    return measurement_limit


try:
    print("name binding: before the name exists ->", report_measurement_limit())
except NameError as missing_name:
    print("name binding: NameError ->", missing_name)

measurement_limit = 10
print("name binding: after the name exists ->", report_measurement_limit())

measurement_limit = 99
print("name binding: same function, new answer ->", report_measurement_limit())


# Closure capture. Python closes over the variable itself, not over the value it held at the time.


def build_index_reporters():
    reporters = []
    for index in range(3):
        reporters.append(lambda: index)  # all three lambdas share the one loop variable, which ends holding 2
    return reporters


def build_index_reporters_by_value():
    reporters = []
    for index in range(3):
        reporters.append(lambda captured_index=index: captured_index)  # a default argument copies the value now
    return reporters


print("closure capture: shared loop variable ->", [reporter() for reporter in build_index_reporters()])
print("closure capture: value copied per lambda ->", [reporter() for reporter in build_index_reporters_by_value()])


# Type checking. Both parts above have already printed, because nothing checked this line until the interpreter
# reached it.

try:
    print("type checking: \"1\" + 1 ->", "1" + 1)
except TypeError as mismatched_operands:
    print("type checking: TypeError ->", mismatched_operands)

print("type checking: with the conversion written out ->", "1" + str(1))
