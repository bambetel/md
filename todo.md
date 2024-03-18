# GO Markdown parser TODO

## Recursive approach

## Tests

Test data for:
- block recognition
- flow parser
- container recognition

### Unit tests

Prefix functions:
- [ ] prefixInside()
- [ ] equalQuotes()

## MdTok() and MdTree() consistency

MdTok() squashes fenced code into one "line", other lines remain separate entries?
- should it stay and also indented code should be put together?
- should the lines stay separate (fences and code) with a tag (`pre`, `pre > code`)?
- line processing "literal mode" turned on inside indented and fenced blocks?

## Prefix

Blockquote marker trailing space present or not

CHECK:
Code in a blockquote catch:
- checking prefix the current way may result in what is 
  supposed to be a `> ~~~` code line to close the code block.

## Frontmatter parsing

Needed: key → value, key → list (tags etc.)

Implementation:
- YAML
- Other?

## HTML output 

- attributes
- auto TOC, heading/part ids.
- use GO templates to make a complete document with metadata, links, styles etc.

## Heading level handling 

### In blockquotes

1. Shift heading number to be higher than the parent element's heading.
2. Strip headings?!

## Blockquote handling 

Recursion or a stack?


## Line joining

In MdTok() (???)

## For LSP (?)

- MdTok - enough to highlight?
- TODO try differentiate indented code from regular paragraph
