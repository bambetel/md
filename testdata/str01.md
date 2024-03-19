# Markdown structure test file

This is a paragraph Lorem ipsum dolor sit amet, officia excepteur ex fugiat
reprehenderit enim labore culpa sint ad nisi Lorem pariatur mollit ex esse
exercitation amet. **Nisi anim** cupidatat excepteur officia. Reprehenderit nostrud
nostrud _ipsum_ Lorem est aliquip amet voluptate voluptate dolor minim nulla est
proident. 

## Test for
no heading hard-wrapping

Just join if more than 4 spaces
     > like here.
      
Although single H1 per document is advised
===
The above should be a settext H1.

NOTE: GFM 3 space rule may apply to both heading text and underscore.

Settext H2 test
---------------

Lorem ipsum dolor sit amet, officia excepteur ex fugiat reprehenderit enim labore culpa sint ad nisi Lorem pariatur mollit ex esse exercitation amet. Nisi anim cupidatat excepteur officia. Reprehenderit nostrud nostrud ipsum Lorem est aliquip amet voluptate voluptate dolor minim nulla est proident. 

A wrapped paragraph with
underline should also yield a H2!
---------

Or just 3, 2, 1 spaces
   like here, however no apparent marker allowed.

> Space after `>` is required (?).
> 
> Blockquote second paragraph.
>
>     Indented code inside

>Test space optionality - no spaces here.
> Neither here.
>  This line should have 1.

Nostrud officia pariatur ut officia. Sit irure elit esse ea nulla
     1. this should not be treated as a li, because of 5 space indent
    2. but this is a list in some plugins, however it shouldn't!
    3. NOT last item.
   4. If 3 spaces allowed, this would make a list.

A section containing this p-heading for a list:
- unordered 
- list 
- items

> Hello world
>---
> this is a:
> 1. broken
>> Nesting example - a blockquote interrupts the list
> 2. a cool list 
> 3. inside
>> LOL
>> Another quote
>>> This VIM plugin doesn't work

    This is an indented text
and this is a next paragraph.

    Indented text

    countinued after a blank line.

1. list item
2. second list item
    
    A paragraph inside a li.

    The second paragraph.

3. Third list item.
    a. simple nesting
    b. list but phat item

        With a paragraph Lorem ipsum dolor sit amet, officia excepteur ex
        fugiat reprehenderit enim labore culpa sint ad nisi Lorem pariatur
        mollit ex esse exercitation amet. Nisi anim cupidatat excepteur
        officia. Reprehenderit nostrud nostrud ipsum Lorem est aliquip amet
        voluptate voluptate dolor minim nulla est proident. 

            functionem monstratCodem(x) (i int){
                decrementis(x)
                returnum 
            }

        Nostrud officia pariatur ut officia. Sit irure elit esse ea nulla sunt
        ex occaecat reprehenderit commodo officia dolor Lorem duis laboris
        cupidatat officia voluptate.
    c. list end
4. Unlike GFM, there are compact/spread lis, not lists!

1. Test for more nesting
2. Compound item
    a. sublist 1 first
    b. sublist 1 second
        - subsublist 1 first
        - subsublist 1 second

3. A new list
4. Because previous items were compact nested

***

\# not a heading!
1. Wrapped
list item
2. Should be treated like a p.
    1. sublist
    2. escaping
    \3. test
3. Also wrapped, but in much more
   elegant a way
    a. with nesting
    b. works
\4. Escaped li

Markdown
: Quite messy, but useful and popular
: Another Markdown standard

Headings
: They are allegedly handled by "outline algorithm"

Edge
case
: A situation provoked to prove Markdown dumb



[1]: Reference

    Can have also a container, just like a li.

[other]: Second reference 
    [2]: In the early implementation
    [3]: This behave like simple list nesting 


