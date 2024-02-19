# Tests

Test data for:
- block recognition
- flow parser
- container recognition

# Prefix

Blockquote marker trailing space present or not

CHECK:
Code in a blockquote catch:
- checking prefix the current way may result in what is 
  supposed to be a `> ~~~` code line to close the code block.

# Frontmatter parsing

Needed: key → value, key → list (tags etc.)

Implementation:
- YAML
- Other?

# Paragraph vs indented code 

- local prefix difference?

# Line pattern matching

## Lists

1. { LI(p,t)} - open a list
2. { LI(p,t), LI(p,t) } - sibling list item
3. { LI(p), blank, ANY(p in p) } - compound LI body
4. ELSE { LI, blank, ANY other } - end list 

## Definition lists 

{ `<p>`, `<dd>`\* } -> `dl > (dt + dd*)+` // one line p only?

## HTML output 

- attributes
- auto TOC, heading/part ids.

# Heading level handling 

## In blockquotes

1. Shift heading number to be higher than the parent element's heading.
2. Strip headings?!

# Blockquote handling 

Recursion or a stack?
