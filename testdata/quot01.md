# Sample Markdown Document with Quoted Text

This is a sample Markdown document demonstrating the use of quoted text.
The paragraph is hard-wrapped and 
spans for a few lines.

_ _ _
## Introduction

Lorem ipsum dolor sit amet, consectetur adipiscing elit. 
1. Nulla ##troll facilisi. 
2. Duis at velit ullamcorper, blandit sem eget, varius metus. 
3. Proin at efficitur ex. 

- ul simple
- second
test line join (second)
- third

- [ ] check simple
- [ ] second
- [ ] third

Not a settext H2
--- --

----

Sed non nisi nec mauris dapibus pharetra. Nam malesuada erat in elit consectetur, eget aliquet ligula sollicitudin. Quisque eu dui at libero convallis venenatis.

1. Nulla ##troll facilisi. 
    a. sublist a
    b. sublist a second
    c. sublist a third
2. Duis at velit ullamcorper, blandit sem eget, varius metus. 
3. Proin at efficitur ex. 

### Subintroduction

1. Nulla ##troll facilisi. 
    
    A p inside the li.

    a. sublist inside the li a
        - further Nesting
        - is possible
    b. sublist inside the li b
    c. sublist inside the li c

2. Duis at velit ullamcorper, blandit sem eget, varius metus. 
3. Proin at efficitur ex. 


> ### Quotes 
> 1. "In the end, it's not the years in your life that count. It's the life in your years."
>     - Abraham Lincoln
> 2. "To be or not to be"
>    Line join in a blockquote.

> Nesting 
> > Nested blockquote.
> > 
> > Test line join
> > this should be one line with "Test line join".
> 
> > Test continuity
> > Other subquote.

## Main Content
Vivamus rutrum magna at justo porta, et volutpat eros pharetra [longnote]. Nulla facilisi. Sed vel nisi at sapien ultrices consequat. Pellentesque nec urna et velit rutrum fermentum. Integer lacinia lacus at libero convallis, ut tempor sapien consequat. Should not be joined, h3 is non-breakable!

Łabędź
: Ptak niemy lub krzykliwy
Kaczka
: Ptak kwaczący (first definition)
: Nocnik (second definition)
Gęś
: Ptak gęgający

> "The only way to do great work is to love what you do." - Steve Jobs

## Conclusion

~~~javascript
function MarkdownIsGreat(i) {
    console.log("Markdown is greater than ", i)

    // blank line not to rely only on block merging!
    let i = 2;
}
~~~

> ```Notclosed
>  should not spill out of the bq
>  Final line
New p.

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus vehicula tristique sapien, et sollicitudin urna tincidunt et. Integer accumsan mi ac enim varius, a venenatis nulla viverra. Nullam auctor, justo id mattis condimentum, neque sapien consectetur justo, at feugiat urna nunc sed turpis. [1]

> "The future belongs to those who believe in the beauty of their dreams." - Eleanor Roosevelt

> ```javascript
> function MarkdownIsGreat(i) {
> > // the second markers are in fence prefix, therefore a part of the 
> > // code block, not nested blockquote!
> >     console.log("Markdown is greater than ", i)
> }
> ```

[1]: Reference 1
[longnote]: Reference 2

    Can have a container.
```Crashtest
