# GO Markdown parser TODO

## Current

### Interface MdTok -> MdTree (and MdFlow TODO)

- Separation of concerns
- Where does flow content processing go?

1. MdTok
    - tell what each line is
        - is in blockquote 
        - is a literal pre text
        - is a heading, paragraph
        - a list item (incl. logical nesting)
2. MdTree
    - use MdTok output directly
    - process flow content (?)
    - output AST

### Idea: Tags vs tokens.

- Tags: acceptable output HTML tags
- Tokens:
    - hr 
    - hr/hst2 -> test neighbours
    - hst1
    - li+type
    - pre
    - p - default

#### Fenced code vs pre - don't merge!

Separate token.

### Processing

- List type handling
- dl, dt, dd grouping
    - TODO: how
        - could just get a regular paragraph and check if it matches a pattern then
        - rules: term in single line?
        - wrappable definitions?

            Snow
            goose
            : a aprticular species of geese
                that lives somewhere and eats something
            : a goose made of snow

- "smart" block patterns
    - figures 
    - heading + list sections
- Global data 
    - metadata from frontmatter 
    - references
        - handling repeated ids
        - in blockquotes (?)
- Nesting rules
    - headings in blockquotes, lis 

## Create readme

- Custom Md version
- Supported extensions

When:
- interfaces concept ready

## Versions

### Recursive approach

Basically works. Using GO slice magic.

### TODO: strictly linear approach 

State:
- container 
- block

1. get the first line and classify into container/block
2. Repeat keeping track of state:
    1. check if line is block continuation
        - yes: join to previous line
        - no: maybe container end
    2. else check if new container 
        - yes: push state
    3. start a new block, flag if can be continued (hard wrapping)
 
## Interfaces 

### MdTok

Handling fenced code blocks.

## Tests

Test data for:
- block recognition
- flow parser
- container recognition

Concept: HTML/Md normalization

### Unit tests

- MdTok() always keeps line number if a text file is given

MdTok functions:
- isEmptyLine()
- getLinePrefix()

## MdTok() and MdTree() consistency

MdTok() squashes fenced code into one "line", other lines remain separate entries?
- should it stay and also indented code should be put together?
- should the lines stay separate (fences and code) with a tag (`pre`, `pre > code`)?
- line processing "literal mode" turned on inside indented and fenced blocks?

### MdTree strategy

- merge block lines with Join=true
- merge adjacent `pre` lines into blocks
- build lists
    - ol/ul
    - reference list 
    - dl, dt, dd

## Prefix

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

## Validation, normalization, correction

- MdTok output:
    - error
    - optional normalization
    - optional corrections 
    - optional hints

- MdTok - should give enough info for syntax highlighting
