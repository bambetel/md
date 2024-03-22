# Test heading
Just a container nesting test. Consists only of the container level elements.
To simplify, only single numerc list markers are put here.
Also blockqutoes are normalized to '>'

    This is meant for the first phase of Md parsing, no block structure apart of li markers are checked.

1. A list
1. Item
    2. subitem
  what?
1. Lorem 

1. Simple
    1. Nesting
    \1. escaped li - this should be merged to previous li
    1. Continued
            this too
    2. This too
        should be merged

        > With a blockquote
        >    Test

1. List
        1. Not a list item.
    1. Subitem
1. Item

1. Phat 

    Container first block (prefix = ____)

    1. List inside a phat element 
        1. Indent - simple nesting
        1. Continue 
        a. test

>> Blockquote normalized
>
>> Inside
> Paragraph 
>> 1. List
> 1. Outer list
>
>     > with a quote (?)
>     > spaces below:
>          
>         > this is some code
>         > yes, literal gt char, not marker!
>     1. nested lists
>         2. simple nesting
>          3. yes or no?
>    
>     1. list nested in container
          1. equivalent spaces but no marker!
>
> 1. Outer list
