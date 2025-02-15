# COOL Language Grammar and Symbol Explanation

program ::= [class;]+

class ::= class TYPE [inherits TYPE] { [feature;]* }

feature ::= ID([formal [, formal]*]): TYPE { expr }
          | ID: TYPE [<- expr]

formal ::= ID: TYPE

expr ::= ID <- expr
       | expr[@TYPE].ID([expr [, expr]*])
       | ID([expr [, expr]*])
       | if expr then expr else expr fi
       | while expr loop expr pool
       | { [expr;]+ }
       | let ID: TYPE [<- expr] [, ID: TYPE [<- expr]]* in expr
       | case expr of [ID: TYPE => expr;]+ esac
       | new TYPE
       | isvoid expr
       | expr + expr
       | expr - expr
       | expr * expr
       | expr / expr
       | ~expr
       | expr < expr
       | expr <= expr
       | expr = expr
       | not expr
       | (expr)
       | ID
       | integer
       | string
       | true
       | false

Figure 1: Cool syntax.

---

**Explanation of Symbols:**

*   **`::=`**  This symbol means "is defined as". It separates the name of a grammar rule (on the left) from its definition (on the right). For example, `program ::= [class;]+` means "a program is defined as one or more classes followed by semicolons".

*   **`|`** This symbol means "or". It indicates alternatives in a grammar rule. For example, `feature ::= ID([formal [, formal]*]): TYPE { expr } | ID: TYPE [<- expr]` means a feature can be either a method definition or an attribute definition.

*   **`[]`** Square brackets indicate optional elements.  For example, `class TYPE [inherits TYPE]` means that in a class definition, the `inherits TYPE` part is optional.  Similarly, `[feature;]*` means that a class can contain zero or more features, each followed by a semicolon.

*   **`+`** The plus sign `+` means "one or more repetitions". For example, `[class;]+` means there must be at least one class definition, and there can be more, each followed by a semicolon.  `[expr;]+` within curly braces `{}` means there must be at least one expression, and there can be more, each followed by a semicolon. `[formal [, formal]*]` means there must be at least one formal parameter, and there can be more, separated by commas. `[ID: TYPE => expr;]+` in the `case` expression means there must be at least one case branch, and there can be more, each followed by a semicolon.

*   **`*`** The asterisk `*` means "zero or more repetitions". For example, `[feature;]*` in a class definition means a class can have zero or more features inside its curly braces. `[formal [, formal]*]` in method definitions means a method can have zero or more formal parameters inside its parentheses.

*   **`()`** Parentheses are used for grouping. For example, `(expr)` in the grammar indicates that an expression can be enclosed in parentheses for precedence control.

*   **`{}`** Curly braces are used to group a sequence of items.  In `class TYPE [inherits TYPE] { [feature;]* }`, curly braces enclose the list of features within a class definition.  In `{ [expr;]+ }`, curly braces enclose a block of expressions.

*   **`,`** Comma is used as a separator, for example, to separate formal parameters in a method definition `[formal [, formal]*]` or in `let` expressions `[, ID: TYPE [<- expr]]*`.

*   **`;`** Semicolon is used as a terminator. For example, `[class;]+` means each class definition is terminated by a semicolon. `[feature;]*` within a class means each feature is terminated by a semicolon. `[expr;]+` within curly braces means each expression in a block is terminated by a semicolon. `[ID: TYPE => expr;]+` in the `case` expression means each case branch is terminated by a semicolon.

*   **`@`**  The `@` symbol in `expr[@TYPE].ID([expr [, expr]*])` is used for static dispatch, allowing a method to be called from a specific ancestor class.

*   **`<-`** The `<-` symbol is used for assignment in expressions like `ID <- expr` and for initialization in attribute definitions `ID: TYPE [<- expr]` and `let` expressions `let ID: TYPE [<- expr]`.

*   **`=>`** The `=>` symbol is used in `case` expressions to separate the type and identifier from the expression to be evaluated in each case branch: `ID: TYPE => expr`.

*   **`~`** The `~` symbol represents integer negation (unary minus).

*   **Keywords and Terminals**: Words like `class`, `inherits`, `TYPE`, `feature`, `ID`, `formal`, `expr`, `if`, `then`, `else`, `fi`, `while`, `loop`, `pool`, `let`, `in`, `case`, `of`, `esac`, `new`, `isvoid`, `true`, `false`, `integer`, `string` are keywords or terminal symbols of the COOL language.  `TYPE` and `ID` represent placeholders for actual type names and identifiers, respectively. `integer` and `string` represent integer and string literals.