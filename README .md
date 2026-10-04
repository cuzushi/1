# go-reloaded

A small text editing and auto-correction tool written in Go.
It reads a text file, applies a set of formatting rules, and writes the result to another file.
Only the Go standard library is used.

## Contents

1. [Usage](#usage)
2. [Supported rules](#supported-rules)
3. [How it works](#how-it-works)
4. [Cases that work](#cases-that-work)
5. [Cases that do not work (known limitations)](#cases-that-do-not-work-known-limitations)
6. [Design choices](#design-choices)
7. [Error handling](#error-handling)
8. [Project structure](#project-structure)
9. [Tests](#tests)

## Usage

```bash
go run . <input_file> <output_file>
```

Example:

```bash
$ cat sample.txt
Simply add 42 (hex) and 10 (bin) and you will see the result is 68.

$ go run . sample.txt result.txt

$ cat result.txt
Simply add 66 and 2 and you will see the result is 68.
```

The program needs exactly two arguments. The output file is created if it does not exist, and **overwritten** if it does.

## Supported rules

| Rule | Input | Output |
|------|-------|--------|
| `(hex)` converts the previous word from hexadecimal to decimal | `1E (hex) files were added` | `30 files were added` |
| `(bin)` converts the previous word from binary to decimal | `It has been 10 (bin) years` | `It has been 2 years` |
| `(up)` makes the previous word uppercase | `Ready, set, go (up) !` | `Ready, set, GO!` |
| `(low)` makes the previous word lowercase | `I should stop SHOUTING (low)` | `I should stop shouting` |
| `(cap)` capitalizes the previous word | `Welcome to the Brooklyn bridge (cap)` | `Welcome to the Brooklyn Bridge` |
| `(up, N)` `(low, N)` `(cap, N)` change the previous N words | `This is so exciting (up, 2)` | `This is SO EXCITING` |
| `. , ! ? : ;` stick to the previous word, with one space after | `there ,and then BAMM !!` | `there, and then BAMM!!` |
| Groups of punctuation stay together | `I was thinking ... You were right` | `I was thinking... You were right` |
| `' '` quote marks go right next to the words inside | `' I am the most well-known '` | `'I am the most well-known'` |
| `a` becomes `an` before a vowel or `h` | `There it was. A amazing rock!` | `There it was. An amazing rock!` |

## How it works

The text goes through a pipeline of four stages. The order matters: each stage expects the previous one to be finished.

```
input file
   |  read the whole file into a string
   v
split into lines, then each line into words (tokens)
   |
   |  Stage 1: commands      (hex) (bin) (up) (low) (cap) and the "(up, N)" versions
   |  Stage 2: a -> an       done after the commands, so the next word is final
   v
join the words back into lines, and the lines into one text
   |  Stage 3: quote()               attach ' marks to the words
   |  Stage 4: formatPunctuation()   fix the spaces around . , ! ? : ;
   v
output file
```

### Stage 1: commands

- Each line is split into words with `strings.Fields`, so several spaces count as one.
- The program walks through the words from left to right. When a word is **exactly** a command (`(up)`, `(hex)`, `(up,` ...), it changes the previous word(s) and then deletes the command from the list.
- `(hex)` and `(bin)` use `strconv.ParseUint`, so values up to 64 bits (`18446744073709551615`) work. If the previous word is not a valid number, it stays unchanged and the command is still removed.
- `(up, 2)` is split into two words, `(up,` and `2)`. The number is read from the second one. If N is bigger than the number of words before the command, it applies to all of them (the loop is clamped, so a huge N is fast).
- Commands only affect the **current line**. `(up, 5)` never changes words from the previous line.
- Chaining works, because commands are applied from left to right: `bRoOkLyN (up) (cap)` gives `Brooklyn`.
- `(cap)` works on characters (runes), not bytes, so it is safe for `élan`, `über` and similar words.

### Stage 2: `a` to `an`

A second loop runs **after** all commands are removed. For each word that is exactly `a` or `A`, the program looks at the next word, skipping opening quote and bracket characters (`' " ( [ {`). If its first letter is a vowel or `h`, then `a` becomes `an` and `A` becomes `An`.

Because it runs after the commands, `a (cap) apple` gives `An apple`.

### Stage 3: quotes

`quote()` works line by line and uses a flag `inquote`. A word that is exactly `'` opens the quote the first time and closes it the second time. The opening mark is glued to the front of the next word, and the closing mark to the end of the previous word. An apostrophe inside a word (`don't`) is never touched.

### Stage 4: punctuation

Three regular expressions run in this order:

1. join punctuation that was split by spaces (`. . .` becomes `...`), repeated until nothing changes;
2. remove spaces and tabs before a punctuation mark;
3. add one space between a punctuation run and a following letter or digit.

The expressions use `[ \t]` and not `\s`, so a line break is never removed. They use `\p{L}\p{N}` and not `[a-zA-Z0-9]`, so accented and non-Latin letters work.

### A complete example

```
i have a (up) apple ,and ' 1E (hex) ' !
```

1. Words: `i` `have` `a` `(up)` `apple` `,and` `'` `1E` `(hex)` `'` `!`
2. Commands: `(up)` makes `a` into `A`; `(hex)` turns `1E` into `30`.
3. `a` to `an`: `A` is followed by `apple` (a vowel), so it becomes `An`.
4. Quotes: `' 30 '` becomes `'30'`.
5. Punctuation: the space before `,` and `!` is removed, and a space is added after `,`.

Result: `i have An apple, and '30'!`

## Cases that work

| Case | Input | Output |
|------|-------|--------|
| All examples from the subject | see [Supported rules](#supported-rules) | as shown |
| Big numbers | `FFFFFFFFFFFFFFFF (hex)` | `18446744073709551615` |
| Lowercase or uppercase hex | `ff (hex)` and `FF (hex)` | `255` |
| Binary | `00101 (bin)` | `5` |
| Several commands in one line | `ff (hex) and 11 (bin)` | `255 and 3` |
| Chained commands | `bRoOkLyN (up) (cap)` | `Brooklyn` |
| A number larger than the line | `hello world (up, 5)` | `HELLO WORLD` |
| A huge number | `hello (up, 9223372036854775807)` | `HELLO` |
| A command never crosses a line | `one` / `two (up, 5)` | `one` / `TWO` |
| A command at the start of a line | `(up) hello` | `hello` |
| Invalid number: the word stays, the command is removed | `zz (hex)` | `zz` |
| Accented letters | `élan (cap)`, `été (up)` | `Élan`, `ÉTÉ` |
| Accented letters after punctuation | `hello,élan` | `hello, élan` |
| Punctuation groups | `what ?!?`, `hmm .... ok`, `no way ...!` | `what?!?`, `hmm.... ok`, `no way...!` |
| Punctuation glued to the next word | `wait...what` | `wait... what` |
| A line starting with punctuation keeps its line break | `line one` / `...and then` | `line one` / `... and then` |
| Quote with punctuation inside | `' really ? ' he asked` | `'really?' he asked` |
| Quote with a command inside | `he said ' go (up) ' !` | `he said 'GO'!` |
| Several quotes | `' yes ' or ' no '` | `'yes' or 'no'` |
| Apostrophes in words are left alone | `don't stop, it's fine` | `don't stop, it's fine` |
| `a` before a quoted word | `i have a 'apple'` | `i have an 'apple'` |
| `a` before a bracket | `i have a (apple)` | `i have an (apple)` |
| `a` before a consonant | `a banana` | `a banana` |
| Words that only contain `a` | `data entry`, `area` | unchanged |
| `a` followed by a command | `a (cap) apple` | `An apple` |
| Normal text in parentheses is not a command | `I (really) like it` | unchanged |
| Unknown command stays in the text | `hello (foo)` | `hello (foo)` |
| Extra spaces and tabs | `hello   ,   world` | `hello, world` |
| Empty file | (nothing) | (nothing) |

## Cases that do not work (known limitations)

These are things the subject does not describe, or that the program does not handle. They are known and documented on purpose.

| # | Case | What happens | What you might expect |
|---|------|--------------|-----------------------|
| 1 | Command glued to a word or a number: `a(up) apple`, `go(up)`, `so exciting (up,2)` | The command is **not applied** and stays in the text. (The punctuation stage may even add a space: `(up,2)` becomes `(up, 2)`.) | The command applied |
| 2 | Command glued to punctuation: `go (up)!` | Not applied, stays in the text | `GO!` |
| 3 | Commands are case sensitive: `(UP)`, `(Hex)` | Not applied | applied |
| 4 | Quote marks glued to a word: `'hello '` | Not recognized, left as it is | `'hello'` |
| 5 | Unpaired quote: `' hello` | The mark is attached to the next word and never closed: `'hello` | |
| 6 | A quote that spans two lines | Not paired | |
| 7 | Decimals and times: `3.14`, `10:30` | `3. 14`, `10: 30` (the rule is applied literally) | `3.14`, `10:30` |
| 8 | Abbreviations and addresses: `U.S.A.`, `example.com` | `U. S. A.`, `example. com` | unchanged |
| 9 | Punctuation split by spaces across two sentences: `Hi! ...ok` | `Hi!... ok` (the two groups are merged) | |
| 10 | `(up, N)` counts lone punctuation and quote marks as words: `so , exciting (up, 2)` | `so, EXCITING` (it changed the `,` and `exciting`) | `SO, EXCITING` |
| 11 | The letter `A` on its own: `Vitamin A is good` | `Vitamin An is good` (the rule applies to every `a`) | unchanged |
| 12 | `a` before a word starting with a non-ASCII vowel: `a école` | unchanged (only ASCII vowels are checked) | `an école` |
| 13 | `an` is never changed back to `a` | `an banana` stays | |
| 14 | Hex with a prefix or a sign: `0x1F (hex)`, `-1F (hex)` | The word stays, the command is removed | `31` |
| 15 | Hex or binary values larger than 64 bits | The word stays, the command is removed | the decimal value |
| 16 | Very large single lines (tens of thousands of commands on one line) | Slow: deleting a command shifts the list, so the time grows with the square of the line size. Normal files are not affected | |

## Design choices

These are decisions where the subject says nothing. They are intentional.

- **`(up)` before a vowel word.** `a (up) apple` gives `An apple`, not `AN apple`. The command makes `a` into `A` first, and then the `a` rule changes `A` into `An`, just like the subject's example `A amazing rock` giving `An amazing rock`.
- **Empty lines are removed.** Only lines with no characters at all are dropped. A line that contains only spaces, or only a `\r` (blank lines in a Windows file), becomes an empty line and is kept.
- **No newline at the end.** The output does not end with a line break, even if the input did.
- **Spacing is normalized.** Several spaces or tabs become one space, and leading and trailing spaces on a line are removed.
- **Windows line endings** (`\r\n`) are converted to `\n`.
- **The `h` rule.** As the subject says, `a` becomes `an` before a word starting with `h` (`an hour`, and also `an house`).

## Error handling

| Situation | Behavior | Exit code |
|-----------|----------|-----------|
| Wrong number of arguments | prints an error message | 1 |
| Input file cannot be read (missing, no permission) | prints the error | 1 |
| Output file cannot be written (for example, the folder does not exist) | prints the error to stderr | 1 |
| Success | no output | 0 |

Every index is checked before it is used, so odd input (a command at the start of a line, a missing number, an unpaired quote) does not crash the program.

## Project structure

```
.
├── main.go               # main, ProcessText (commands and a/an), capitalize, isvowel, nextWordStart
├── quote.go              # quote()               - pairs and places the ' marks
├── formatPunctuation.go  # formatPunctuation()   - spacing around punctuation (regular expressions)
└── *_test.go             # unit tests
```

All three source files are in `package main`. `ProcessText` is a pure function (a string in, a string out), so it is easy to test without touching any files.

## Tests

```bash
go test                       # run all tests
go test -v                    # see every test
go test -v -run TestName      # run one test or group
```

The tests cover the examples from the subject, each rule with edge cases (big numbers, accented letters, a number larger than the line, chained commands), how the rules interact (quotes with punctuation, commands inside quotes), and inputs that must never crash or hang.
