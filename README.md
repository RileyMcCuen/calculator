# calculator
Simple precedence adhering command line calculator. Created after wishing I had one in many occasions that worked the same on all platforms.

# Usage

Compile or install directly with Go.

The calculator has two modes:

## Evaluate Expression

Evaluates any expressions specified as args to the calculator. Variables defined in earlier expressions can be used in later expressions.

example: `calculator "1+2^.5/6"`

example: `calculator "x=1+2^.5/6" "x+10.2"`

## Interactive

If there are no extra arguments interactive mode is started.
Starts an interactive calculator session that (depending on terminal implementation) allows you to cycle through your past expressions.
Will contintue evaluating expressions until the keyword 'exit' is entered, in which case the calculator program exits
Can use variables to save past calculations.

example of session
```
# start the calculator
calculator

# input
1+2

#output
=3

x=5
=5
x
=5
y=x+1
=6
y
=6
```
