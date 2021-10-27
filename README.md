# calculator
Simple precedence adhering command line calculator. Created after wishing I had one in many occasions that worked the same on all platforms.

# Usage

Compile or install directly with Go.

The calculator has two modes:

## Evaluate Expression

Evaluates the single expression specified by flag "e" and prints the value out.

example: `calculator -e "1+2^.5/6"`

## Interactive

Starts an interactive calculator session that (depending on terminal implementation) allows you to cycle through your past expressions.
Will contintue evaluating expressions until the keyword 'exit' is entered, in which case the calculator program exits
Can use variables to save past calculations.

```
calculator
1+2
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
