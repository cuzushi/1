# Go Text Tool

A small command-line tool written in Go that reads a text file, applies a set of editing and auto-correction rules, and writes the result to a new `.txt` file.

## Usage

```bash
go run . input.txt output.txt
```

- `input.txt`: the file to read.
- `output.txt`: the file to write. It must end in `.txt` (the tool refuses to write anything else, so you can't accidentally overwrite a `.go` file).

## What it can do

### 1. Number conversion

| Tag | What it does | Example |
|-----|--------------|---------|
| `(hex)` | Converts the previous word from hexadecimal to decimal | `1E (hex) files` → `30 files` |
| `(bin)` | Converts the previous word from binary to decimal | `10 (bin) years` → `2 years` |

If the previous word is not a valid number, it is left unchanged and the tag is simply removed.

Hexadecimal (base 16)

• Maximum: 7fffffffffffffff

• Minimum: -8000000000000000

Binary (base 2)

• Maximum: 111111111111111111111111111111111111111111111111111111111111111

• Minimum: -1000000000000000000000000000000000000000000000000000000000000000 



### 2. Case conversion

| Tag | What it does | Example |
|-----|--------------|---------|
| `(up)` | Uppercases the previous word | `go (up)` → `GO` |
| `(low)` | Lowercases the previous word | `SHOUTING (low)` → `shouting` |
| `(cap)` | Capitalizes the previous word | `bridge (cap)` → `Bridge` |

Add a number to apply the change to the previous **N** words:

```
This is so exciting (up, 2)   →   This is SO EXCITING
```

Works the same with `(low, N)` and `(cap, N)`.

You should respect this format to work " (low, nb) " (u can put one space or  multiple spaces between low, and nb) and don't put a space between nb and ")" same thing with low,up and cap 

you should respect the space between (low, nb) and other words bcz any attached with a word it will treat as part of the world

### 3. Article correction

`a` becomes `an` (and `A` becomes `An`) when the next word starts with a vowel (`a, e, i, o, u`) or `h`.

```
There it was. A amazing rock!   →   There it was. An amazing rock!
```

Quotes and brackets in front of the next word are ignored when checking.

### 4. Quote formatting

A single quote `'` standing alone is attached to the word next to it, so quotes hug the text they surround:

```
He said: ' I am happy '   →   He said: 'I am happy'
```
it shouldn't be attached to any word or it would treated like a part of the word 
it wouldn't work if there's any attached with another characters

### 5. Punctuation formatting

The characters `. , ! ? : ;` are:

- attached to the word before them (no space before),
- followed by a single space when a word comes after them,
- kept together when grouped, like `...`, `!!` or `!?`.

```
I was sitting over there ,and then BAMM !!   →   I was sitting over there, and then BAMM!!
```

### 6. Cleanup

Extra spaces between words are collapsed into one, and empty lines are removed.

## Full example

**input.txt**

```
it (cap) was the best of times ,it was the worst of times (up, 3) .
Simply add 42 (hex) and 10 (bin) ,then a amazing result appears...
As someone said: ' stay curious '
```

**output.txt**

```
It was the best of times, it was THE WORST OF TIMES.
Simply add 66 and 2, then an amazing result appears...
As someone said: 'stay curious'
```

## Project structure

| File | Role |
|------|------|
| `main.go` | Entry point, argument checks, reads/writes files, and `ProcessText` (tags, number conversion, case changes, `a` → `an`) |
| `quote.go` | `quote()` — handles single-quote formatting |
| `formatPunctuation.go` | `formatPunctuation()` — handles punctuation spacing with regular expressions |

## Processing order

1. Read the input file.
2. Apply the tags (`hex`, `bin`, `up`, `low`, `cap`) line by line.
3. Fix `a` → `an`.
4. Format quotes.
5. Format punctuation.
6. Write the result to the output file.

## Error handling

The program exits with an error message if:

- the number of arguments is not exactly 2,
- the output file does not end with `.txt`,
- the input file cannot be read or the output file cannot be written.


