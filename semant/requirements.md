# The Cool Reference Manual∗

## Contents

1Introduction3

2Getting Started3

3Classes..................................................4

3.1Features . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .4

3.2Inheritance. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .5

4Types..................................................6

4.1SELF TYPE. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .6

4.2Type Checking . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .7

5Attributes8

5.1Void . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .8

6Methods8

7Expressions...............................................9

7.1Constants . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .9

7.2Identiﬁers . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .9

7.3Assignment . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .9

7.4Dispatch . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .10

7.5Conditionals . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .10

7.6Loops. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .11

7.7Blocks . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .11

7.8Let . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .11

7.9Case . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .12

7.10 New . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .12

7.11 Isvoid. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .12

7.12 Arithmetic and Comparison Operations. . . . . . . . . . . . . . . . . . . . . . . . . . . .13

∗Copyright c⃝1995-2000 by Alex Aiken. All rights reserved.

8Basic Classes13

8.1Object . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .13

8.2IO . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .13

8.3Int . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .14

8.4String . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .14

8.5Bool . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .14

9Main Class14

10 Lexical Structure.................................................14

10.1 Integers, Identiﬁers, and Special Notation . . . . . . . . . . . . . . . . . . . . . . . . . . .15

10.2 Strings . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .15

10.3 Comments . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .15

10.4 Keywords . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .15

10.5 White Space . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .15

11 Cool Syntax....................................................17

11.1 Precedence. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .17

### 12 Type RulesType Rules17

12.1 Type Environments . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .17

12.2 Type Checking Rules . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .18

## 13 Operational Semantics.............................................22

13.1 Environment and the Store. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .22

13.2 Syntax for Cool Objects . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .24

13.3 Class deﬁnitions................................................................24

13.4 Operational Rules. . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . .25

14 Acknowledgements30

## 1Introduction

This manual describes the programming language Cool: the Classroom Object-Oriented Language. Cool is a small language that can be implemented with reasonable eﬀort in a one semester course. Still, Cool retains many of the features of modern programming languages including objects, static typing, and automatic memory management.

Cool programs are sets of classes. A class encapsulates the variables and procedures of a data type.Instances of a class are objects. In Cool, classes and types are identiﬁed; i.e., every class deﬁnes a type.Classes permit programmers to deﬁne new types and associated procedures (or methods) speciﬁc to those types. Inheritance allows new types to extend the behavior of existing types.

Cool is an expression language. Most Cool constructs are expressions, and every expression has a value and a type. Cool is type safe: procedures are guaranteed to be applied to data of the correct type.While static typing imposes a strong discipline on programming in Cool, it guarantees that no runtime type errors can arise in the execution of Cool programs.

This manual is divided into informal and formal components. For a short, informal overview, the ﬁrst half (through Section 9) suﬃces. The formal description begins with Section 10.

## 2Getting Started

The reader who wants to get a sense for Cool at the outset should begin by reading and running the example programs in the directory \~cs164/examples. Cool source ﬁles have extension .cl and Cool assembly ﬁles have extension .s. The Cool compiler is \~cs164/bin/coolc. To compile a program:

coolc [ -o fileout ] file1.cl file2.cl ... filen.cl

The compiler compiles the ﬁles file1.cl through filen.cl as if they were concatenated together.Each ﬁle must deﬁne a set of complete classes—class deﬁnitions may not be split across ﬁles. The -o option speciﬁes an optional name to use for the output assembly code. If fileout is not supplied, the output assembly is named file1.s.

The coolc compiler generates MIPS assembly code. Because not all of the machines the course is using are MIPS-based, Cool programs are run on a MIPS simulator called spim. To run a cool program,type

% spim

(spim) load "file.s"

(spim) run

To run a diﬀerent program during the same spim session, it is necessary to reinitialize the state of the simulator before loading the new assembly ﬁle:

(spim) reinit

An alternative—and faster—way to invoke spim is with a ﬁle:

spim -file file.s

This form loads the ﬁle, runs the program, and exits spim when the program terminates. Be sure that spim is invoked using the script \~cs164/bin/spim. There may be another version of spim installed in on some systems, but it will not execute Cool programs. An easy way to be sure of getting the correct

version is to alias spim to \~cs164/bin/spim. The spim manual is available on the course Web page and in the course reader.

The following is a complete transcript of the compilation and execution of \~cs164/examples/list.cl.This program is very silly, but it does serve to illustrate many of the features of Cool.

% coolc list.cl

% spim -file list.s

SPIM Version 5.6 of January 18, 1995

Copyright 1990-1994 by James R. Larus (larus@cs.wisc.edu).

All Rights Reserved.

See the file README a full copyright notice.

Loaded: /home/ee/cs164/lib/trap.handler

5 4 3 2 1

4 3 2 1

3 2 1

2 1

1

COOL program successfully executed

%

### 3Classes

All code in Cool is organized into classes. Each class deﬁnition must be contained in a single source ﬁle,but multiple classes may be deﬁned in the same ﬁle. Class deﬁnitions have the form:

class &lt;type&gt; [ inherits &lt;type&gt; ] {

&lt;feature_list&gt;

};

The notation [ ...] denotes an optional construct. All class names are globally visible. Class names begin with an uppercase letter. Classes may not be redeﬁned.

#### 3.1Features

The body of a class deﬁnition consists of a list of feature deﬁnitions. A feature is either an attribute or a method. An attribute of class A speciﬁes a variable that is part of the state of objects of class A. A method of class A is a procedure that may manipulate the variables and objects of class A.

One of the major themes of modern programming languages is information hiding, which is the idea that certain aspects of a data type’s implementation should be abstract and hidden from users of the data type. Cool supports information hiding through a simple mechanism: all attributes have scope local to the class, and all methods have global scope. Thus, the only way to provide access to object state in Cool is through methods.

Feature names must begin with a lowercase letter. No method name may be deﬁned multiple times in a class, and no attribute name may be deﬁned multiple times in a class, but a method and an attribute may have the same name.

A fragment from list.cl illustrates simple cases of both attributes and methods:

class Cons inherits List {

xcar : Int;

xcdr : List;

isNil() : Bool { false };

init(hd : Int, tl : List) : Cons {

{

xcar &lt;- hd;

xcdr &lt;- tl;

self;

}

}

...

};

In this example, the class Cons has two attributes xcar and xcdr and two methods isNil and init.Note that the types of attributes, as well as the types of formal parameters and return types of methods,are explicitly declared by the programmer.

Given object c of class Cons and object l of class List, we can set the xcar and xcdr ﬁelds by using the method init:

c.init(1,l)

This notation is object-oriented dispatch. There may be many deﬁnitions of init methods in many diﬀerent classes. The dispatch looks up the class of the object c to decide which init method to invoke.Because the class of c is Cons, the init method in the Cons class is invoked. Within the invocation, the variables xcar and xcdr refer to c’s attributes. The special variable self refers to the object on which the method was dispatched, which, in the example, is c itself.

There is a special form new C that generates a fresh object of class C. An object can be thought of as a record that has a slot for each of the attributes of the class as well as pointers to the methods of the class. A typical dispatch for the init method is:

(new Cons).init(1,new Nil)

This example creates a new cons cell and initializes the "Car" of the cons cell to be 1 and the "cdr" to be new Nil.1 There is no mechanism in Cool for programmers to deallocate objects. Cool has automatic memory management; objects that cannot be used by the program are deallocated by a runtime garbage collector.

Attributes are discussed further in Section 5 and methods are discussed further in Section 6.

### 3.2Inheritance

If a class deﬁnition has the form

class C inherits P { ... };

1In this example, Nil is assumed to be a subtype of List.

then class C inherits the features of P. In this case P is the parent class of C and C is a child class of P.

The semantics of C inherits P is that C has all of the features deﬁned in P in addition to its own features. In the case that a parent and child both deﬁne the same method name, then the deﬁnition given in the child class takes precedence. It is illegal to redeﬁne attribute names. Furthermore, for type safety, it is necessary to place some restrictions on how methods may be redeﬁned (see Section 6).

There is a distinguished class Object. If a class deﬁnition does not specify a parent class, then the class inherits from Object by default. A class may inherit only from a single class; this is aptly called “single inheritance."2 The parent-child relation on classes deﬁnes a graph. This graph may not contain cycles. For example, if C inherits from P, then P must not inherit from C. Furthermore, if C inherits from P, then P must have a class deﬁnition somewhere in the program. Because Cool has single inheritance, it follows that if both of these restrictions are satisﬁed, then the inheritance graph forms a tree with Object as the root.

In addition to Object, Cool has four other basic classes: Int, String, Bool, and IO. The basic classes are discussed in Section 8.

## 4Types

In Cool, every class name is also a type. In addition, there is a type SELF TYPE that can be used in special circumstances.

A type declaration has the form x:C, where x is a variable and C is a type. Every variable must have a type declaration at the point it is introduced, whether that is in a let, case, or as the formal parameter of a method. The types of all attributes must also be declared.

The basic type rule in Cool is that if a method or variable expects a value of type P, then any value of type C may be used instead, provided that P is an ancestor of C in the class hierarchy. In other words,if C inherits from P, either directly or indirectly, then a C can be used wherever a P would suﬃce.

When an object of class C may be used in place of an object of class P, we say that C conforms to P or thatc≤p (think: C is lower down in the inheritance tree). As discussed above, conformance is deﬁnedin terms of the inheritance graph.

Deﬁnition 4.1 (Conformance) Let A, C, and P be types.

• A≤A for all types A

• if C inherits from P, thenc≤p

• ifA≤cand c≤p thenA≤P

Because Object is the root of the class hierarchy, it follows that A ≤Object for all types A.

### 4.1SELF TYPE

The type SELF TYPE is used to refer to the type of the self variable. This is useful in classes that will be inherited by other classes, because it allows the programmer to avoid specifying a ﬁxed ﬁnal type at the time the class is written. For example, the program

2Some object-oriented languages allow a class to inherit from multiple classes, which is equally aptly called “multiple inheritance.”

class Silly {

copy() : SELF_TYPE { self };

};

class Sally inherits Silly { };

class Main {

x:Sa11y&lt;-(nevsal1y).copy();

mina∈():Sa11y[x];

};

Because SELF TYPE is used in the deﬁnition of the copy method, we know that the result of copy is the same as the type of the self parameter. Thus, it follows that (new Sally).copy() has type Sally,which conforms to the declaration of attribute x.

Note that the meaning of SELF TYPE is not ﬁxed, but depends on the class in which it is used. In general, SELF TYPE may refer to the class C in which it appears, or any class that conforms to C. When it is useful to make explicit what SELF TYPE may refer to, we use the name of the class C in which SELF TYPE appears as an index SELF TYPEC. This subscript notation is not part of Cool syntax—it is used merely to make clear in what class a particular occurrence of SELF TYPE appears.

From Deﬁnition 4.1, it follows that SELF TYPEX ≤SELF TYPEX. There is also a special conformancerule for SELF TYPE:

$SELF_TYPE_{C}\leq P$ ifc≤p

Finally, SELF TYPE may be used in the following places: new SELF TYPE, as the return type of a method, as the declared type of a let variable, or as the declared type of an attribute. No other uses of SELF TYPE are permitted.

### 4.2Type Checking

The Cool type system guarantees at compile time that execution of a program cannot result in runtime type errors. Using the type declarations for identiﬁers supplied by the programmer, the type checker infers a type for every expression in the program.

It is important to distinguish between the type assigned by the type checker to an expression at compile time, which we shall call the static type of the expression, and the type(s) to which the expression may evaluate during execution, which we shall call the dynamic types.

The distinction between static and dynamic types is needed because the type checker cannot, at compile time, have perfect information about what values will be computed at runtime. Thus, in general,the static and dynamic types may be diﬀerent. What we require, however, is that the type checker’s static types be sound with respect to the dynamic types.

Deﬁnition 4.2 For any expression e, let De be a dynamic type of e and let Se be the static type inferred by the type checker. Then the type checker is sound if for all expressions e it is the case that$D_{e}\leq S_{e}.$

Put another way, we require that the type checker err on the side of overestimating the type of an expression in those cases where perfect accuracy is not possible. Such a type checker will never accept a program that contains type errors. However, the price paid is that the type checker will reject some programs that would actually execute without runtime errors.

An attribute deﬁnition has the form

&lt;id&gt; : &lt;type&gt; [ &lt;- &lt;expr&gt; ];

The expression is optional initialization that is executed when a new object is created. The static type of the expression must conform to the declared type of the attribute. If no initialization is supplied, then the default initialization is used (see below).

When a new object of a class is created, all of the inherited and local attributes must be initialized.Inherited attributes are initialized ﬁrst in inheritance order beginning with the attributes of the greatest ancestor class. Within a given class, attributes are initialized in the order they appear in the source text.

Attributes are local to the class in which they are deﬁned or inherited. Inherited attributes cannot be redeﬁned.

### 5.1Void

All variables in Cool are initialized to contain values of the appropriate type. The special value void is a member of all types and is used as the default initialization for variables where no initialization is supplied by the user. (void is used where one would use NULL in C or null in Java; Cool does not have anything equivalent to C’s or Java’s void type.) Note that there is no name for void in Cool; the only way to create a void value is to declare a variable of some class other than Int, String, or Bool and allow the default initialization to occur, or to store the result of a while loop.

There is a special form isvoid expr that tests whether a value is void (see Section 7.11). In addition,void values may be tested for equality. A void value may be passed as an argument, assigned to a variable,or otherwise used in any context where any value is legitimate, except that a dispatch to or case on void generates a runtime error.

Variables of the basic classes Int, Bool, and String are initialized specially; see Section 8.

## 6Methods

A method deﬁnition has the form

&lt;id&gt;(&lt;id&gt; : &lt;type&gt;,...,&lt;id&gt; : &lt;type&gt;): &lt;type&gt; {&lt;expx&gt;};

There may be zero or more formal parameters. The identiﬁers used in the formal parameter list must be distinct. The type of the method body must conform to the declared return type. When a method is invoked, the formal parameters are bound to the actual arguments and the expression is evaluated; the resulting value is the meaning of the method invocation. A formal parameter hides any deﬁnition of an attribute of the same name.

To ensure type safety, there are restrictions on the redeﬁnition of inherited methods. The rule is simple: If a class C inherits a method f from an ancestor class P, then C may override the inherited deﬁnition of f provided the number of arguments, the types of the formal parameters, and the return type are exactly the same in both deﬁnitions.

To see why some restriction is necessary on the redeﬁnition of inherited methods, consider the following example:

class P {

f(): Int { 1 };

};

class C inherits P {

f(): String { "1"];

};

Let p be an object with dynamic type P. Then

p.f()+1

is a well-formed expression with value 2. However, we cannot substitute a value of type C for p, as it would result in adding a string to a number. Thus, if methods can be redeﬁned arbitrarily, then subclasses may not simply extend the behavior of their parents, and much of the usefulness of inheritance, as well as type safety, is lost.

## 7Expressions

Expressions are the largest syntactic category in Cool.

### 7.1Constants

The simplest expressions are constants. The boolean constants are true and false. Integer constants are unsigned strings of digits such as 0, 123, and 007. String constants are sequences of characters enclosed in double quotes, such as "This is a string." String constants may be at most 1024 characters long.There are other restrictions on strings; see Section 10.

The constants belong to the basic classes Bool, Int, and String. The value of a constant is an object of the appropriate basic class.

### 7.2Identiﬁers

The names of local variables, formal parameters of methods, self, and class attributes are all expressions.The identiﬁer self may be referenced, but it is an error to assign to self or to bind self in a let, a case, or as a formal parameter. It is also illegal to have attributes named self.

Local variables and formal parameters have lexical scope. Attributes are visible throughout a class in which they are declared or inherited, although they may be hidden by local declarations within expres sions. The binding of an identiﬁer reference is the innermost scope that contains a declaration for that identiﬁer, or to the attribute of the same name if there is no other declaration. The exception to this rule is the identiﬁer self, which is implicitly bound in every class.

### 7.3Assignment

An assignment has the form

&lt;id&gt; &lt;- &lt;expr&gt;

The static type of the expression must conform to the declared type of the identiﬁer. The value is the value of the expression. The static type of an assignment is the static type of &lt;expr&gt;.

There are three forms of dispatch (i.e. method call) in Cool. The three forms diﬀer only in how the called method is selected. The most commonly used form of dispatch is

&lt;expr&gt;.&lt;id&gt;(&lt;expr&gt;,...,&lt;expr&gt;)

Consider the dispatch $e_{0}.f(e_{1},\cdots ,e_{n})$. To evaluate this expression, the arguments are evaluated in left to-right order, from $e_{1}$ to $e_{n}.$ Next, $e_{0}$ is evaluated and its class C noted (if $e_{0}$ is void a runtime error is generated). Finally, the method f in class C is invoked, with the value of e0 bound to self in the body of f and the actual arguments bound to the formals as usual. The value of the expression is the value returned by the method invocation.

Type checking a dispatch involves several steps. Assume e0 has static type A. (Recall that this type is not necessarily the same as the type C above. A is the type inferred by the type checker; C is the class of the object computed at runtime, which is potentially any subclass of A.) Class A must have a method f, the dispatch and the deﬁnition of f must have the same number of arguments, and the static type of the ith actual parameter must conform to the declared type of the ith formal parameter.

If f has return type B and B is a class name, then the static type of the dispatch is B. Otherwise, if f has return type SELF TYPE, then the static type of the dispatch is A. To see why this is sound, note that the self parameter of the method f conforms to type A. Therefore, because f returns SELF TYPE, we can infer that the result must also conform to A. Inferring accurate static types for dispatch expressions is what justiﬁes including SELF TYPE in the Cool type system.

The other forms of dispatch are:

&lt;id&gt;(&lt;expr&gt;,...,&lt;expr&gt;)

&lt;expr&gt;@&lt;type&gt;.id(&lt;expr&gt;,...,&lt;expr&gt;)

The ﬁrst form is shorthand for self.&lt;id&gt;(&lt;expr&gt;,...,&lt;expr&gt;).

The second form provides a way of accessing methods of parent classes that have been hidden by redeﬁnitions in child classes.Instead of using the class of the leftmost expression to determine the method, the method of the class explicitly speciﬁed is used. For example, e@B.f() invokes the method f in class B on the object that is the value of e. For this form of dispatch, the static type to the left of “@”must conform to the type speciﬁed to the right of “@”.

### 7.5Conditionals

A conditional has the form

if &lt;expr&gt; then &lt;expr&gt; else &lt;expr&gt; fi

The semantics of conditionals is standard. The predicate is evaluated ﬁrst. If the predicate is true,then the then branch is evaluated. If the predicate is false, then the else branch is evaluated. The value of the conditional is the value of the evaluated branch.

The predicate must have static type Bool. The branches may have any static types. To specify the static type of the conditional, we deﬁne an operation ⊔(pronounced “join”) on types as follows. LetA,B,D be any types other than SELF TYPE. The least type of a set of types means the least element with respect to the conformance relation ≤.

A∪B=the least type C such thatA≤candB≤c

A∪A=A (idempotent)

A∪B=B∪A (commutative)

$SELF_TYPE_{D}\cup A=D\cup A$

Let T and F be the static types of the branches of the conditional.Then the static type of the conditional is T ⊔F. (think: Walk towards Object from each of T and F until the paths meet.)

### 7.6Loops

A loop has the form

while &lt;expr&gt; loop &lt;expr&gt; pool

The predicate is evaluated before each iteration of the loop. If the predicate is false, the loop terminates and void is returned. If the predicate is true, the body of the loop is evaluated and the process repeats.

The predicate must have static type Bool. The body may have any static type. The static type of a loop expression is Object.

### 7.7Blocks

A block has the form

{ &lt;expr&gt;; ... &lt;expr&gt;; }

The expressions are evaluated in left-to-right order. Every block has at least one expression; the value of a block is the value of the last expression. The expressions of a block may have any static types. The static type of a block is the static type of the last expression.

An occasional source of confusion in Cool is the use of semi-colons (“;”).Semi-colons are used as terminators in lists of expressions (e.g., the block syntax above) and not as expression separators.Semi-colons also terminate other Cool constructs, see Section 11 for details.

### 7.8Let

A let expression has the form

let &lt;id1&gt; : &lt;type1&gt; [ &lt;- &lt;expr1&gt; ], ..., &lt;idn&gt; : &lt;typen&gt; [ &lt;- &lt;exprn&gt; ] in &lt;expr&gt;

The optional expressions are initialization; the other expression is the body. A let is evaluated as follows. First &lt;expr1&gt; is evaluated and the result bound to &lt;id1&gt;. Then &lt;expr2&gt; is evaluated and the result bound to &lt;id2&gt;, and so on, until all of the variables in the let are initialized. (If the initialization of &lt;idk&gt; is omitted, the default initialization of type &lt;typek&gt; is used.) Next the body of the let is evaluated. The value of the let is the value of the body.

The let identiﬁers &lt;id1&gt;,...,&lt;idn&gt; are visible in the body of the let. Furthermore, identiﬁers &lt;id1&gt;,...,&lt;idk&gt; are visible in the initialization of &lt;idm&gt; for any m &gt; k.

If an identiﬁer is deﬁned multiple times in a let, later bindings hide earlier ones. Identiﬁers introduced by let also hide any deﬁnitions for the same names in containing scopes. Every let expression must introduce at least one identiﬁer.

The type of an initialization expression must conform to the declared type of the identiﬁer. The type of let is the type of the body.

The &lt;expr&gt; of a let extends as far (encompasses as many tokens) as the grammar allows.

### 7.9Case

A case expression has the form

case &lt;expr0&gt; of

&lt;id1&gt;:&lt;type1&gt;=&gt; &lt;eexpr1&gt;;

. . .

&lt;idn&gt;:&lt;typen&gt;=&gt;&lt;exprn&gt;;

esac

Case expressions provide runtime type tests on objects. First, expr0 is evaluated and its dynamic type C noted (if expr0 evaluates to void a run-time error is produced). Next, from among the branches the branch with the least type &lt;typek&gt; such that C ≤&lt;typek&gt; is selected. The identiﬁer &lt;idk&gt; is boundto the value of &lt;expr0&gt; and the expression &lt;exprk&gt; is evaluated. The result of the case is the value of &lt;exprk&gt;.If no branch can be selected for evaluation, a run-time error is generated.Every case expression must have at least one branch.

For each branch, let $T_{i}$.The identiﬁer id introduced by a branch of a case hides any variable or attribute deﬁnition for id visible be the static type of &lt;expri&gt;. The static type of a case expression is $U_{1\leq i\leq nT_{i}}$in the containing scope.

The case expression has no special construct for a “default” or “otherwise” branch. The same aﬀect is achieved by including a branch

x:0bject⇒⋯

because every type is≤to Object.

The case expression provides programmers a way to insert explicit runtime type checks in situa tions where static types inferred by the type checker are too conservative. A typical situation is that a programmer writes an expression e and type checking infers that e has static type P. However, the programmer may know that, in fact, the dynamic type of e is always C for some. This informationcan be captured using a case expression:c≤P.

case e ofx:C⇒⋯

In the branch the variable x is bound to the value of e but has the more speciﬁc static type C.

### 7.10New

A new expression has the form

new &lt;type&gt;

The value is a fresh object of the appropriate class. If the type is SELF TYPE, then the value is a fresh object of the class of self in the current scope. The static type is &lt;type&gt;.

### 7.11Isvoid

The expression

isvoid expr

evaluates to true if expr is void and evaluates to false if expr is not void.

### 7.12Arithmetic and Comparison Operations

Cool has four binary arithmetic operations: +, -, *, /. The syntax is

expr1 &lt;op&gt; expr2

To evaluate such an expression ﬁrst expr1 is evaluated and then expr2. The result of the operation is the result of the expression.

The static types of the two sub-expressions must be Int. The static type of the expression is Int.Cool has only integer division.

Cool has three comparison operations: &lt;,〈=,=.For&lt;and &lt;=the rules are exactly the same as for the binary arithmetic operations, except that the result is a Bool. The comparison = is a special case. If either &lt;expr1&gt; or &lt;expr2&gt; has static type Int, Bool, or String, then the other must have the same static type. Any other types, including SELF TYPE, may be freely compared. On non-basic objects,equality simply checks for pointer equality (i.e., whether the memory addresses of the objects are the same). Equality is deﬁned for void.

In principle, there is nothing wrong with permitting equality tests between, for example, Bool and Int. However, such a test must always be false and almost certainly indicates some sort of programming error. The Cool type checking rules catch such errors at compile-time instead of waiting until runtime.

Finally, there is one arithmetic and one logical unary operator. The expression \~&lt;expr&gt; is the integer complement of &lt;expr&gt;.The expression &lt;expr&gt; must have static type Int and the entire expression has static type Int. The expression not &lt;expr&gt; is the boolean complement of &lt;expr&gt;. The expression &lt;expr&gt; must have static type Bool and the entire expression has static type Bool.

## 8Basic Classes

### 8.1Object

The Object class is the root of the inheritance graph.Methods with the following declarations are deﬁned:

abort() : Object

type_name() : String

copy() : SELF_TYPE

The method abort halts program execution with an error message. The method type name returns a string with the name of the class of the object. The method copy produces a shallow copy of the object.3

### 8.2IO

The IO class provides the following methods for performing simple input and output operations:

out_string(x : String) : SELF_TYPE

out_int(x : Int) : SELF_TYPE

in_string() : String

in_int() : Int

$^{3}A$shallow copy of a copies a itself, but does not recursively copy objects that a points to.

The methods out string and out int print their argument and return their self parameter.The method in string reads a string from the standard input, up to but not including a newline character.The method in int reads a single integer, which may be preceded by whitespace. Any characters following the integer, up to and including the next newline, are discarded by in int.

A class can make use of the methods in the IO class by inheriting from IO. It is an error to redeﬁne the IO class.

### 8.3Int

The Int class provides integers. There are no methods special to Int. The default initialization for variables of type Int is 0 (not void). It is an error to inherit from or redeﬁne Int.

### 8.4String

The String class provides strings. The following methods are deﬁned:

length() : Int

concat(s : String) : String

substr(i : Int, l : Int) : String

The method length returns the length of the self parameter. The method concat returns the string formed by concatenating s after self. The method substr returns the substring of its self parameter beginning at position iwith length l; string positions are numbered beginning at 0. A runtime error is generated if the speciﬁed substring is out of range.

The default initialization for variables of type String is "" (not void). It is an error to inherit from or redeﬁne String.

### 8.5Bool

The Bool class provides true and false. The default initialization for variables of type Bool is false (not void). It is an error to inherit from or redeﬁne Bool.

## 9Main Class

Every program must have a class Main. Furthermore, the Main class must have a method main that takes no formal parameters. The main method must be deﬁned in class Main (not inherited from another class). A program is executed by evaluating (new Main).main().

The remaining sections of this manual provide a more formal deﬁnition of Cool. There are four sections covering lexical structure (Section 10), grammar (Section 11), type rules (Section 12), and operational semantics (Section 13).

## 10Lexical Structure

The lexical units of Cool are integers, type identiﬁers, object identiﬁers, special notation, strings, key words, and white space.

### 10.1Integers, Identiﬁers, and Special Notation

Integers are non-empty strings of digits 0-9. Identiﬁers are strings (other than keywords) consisting of letters, digits, and the underscore character. Type identiﬁers begin with a capital letter; object identiﬁers begin with a lower case letter. There are two other identiﬁers, self and SELF TYPE that are treated specially by Cool but are not treated as keywords. The special syntactic symbols (e.g., parentheses,assignment operator, etc.) are given in Figure 1.

### 10.2Strings

Strings are enclosed in double quotes "...". Within a string, a sequence ‘\c’ denotes the character ‘c’,with the exception of the following:

\bbackspace

\ttab

\nnewline

\fformfeed

A non-escaped newline character may not appear in a string:

"This \

is OK"

"This is not

OK"

A string may not contain EOF. A string may not contain the null (character \0). Any other charactermay be included in a string. Strings cannot cross ﬁle boundaries.

### 10.3Comments

There are two forms of comments in Cool. Any characters between two dashes “--” and the next newline (or EOF, if there is no next newline) are treated as comments. Comments may also be written by enclosing text in(*,⋯*). The latter form of comment may be nested. Comments cannot cross ﬁle boundaries.

### 10.4Keywords

The keywords of cool are: class, else, false, ﬁ, if, in, inherits, isvoid, let, loop, pool, then, while,case, esac, new, of, not, true. Except for the constants true and false, keywords are case insensitive.To conform to the rules for other objects, the ﬁrst letter of true and false must be lowercase; the trailing letters may be upper or lower case.

### 10.5White Space

White space consists of any sequence of the characters: blank (ascii 32), \n (newline, ascii 10), \f (formfeed, ascii 12), \r (carriage return, ascii 13), \t (tab, ascii 9), \v (vertical tab, ascii 11).

program :=$[class;]^{+}$

class ::= class TYPE [inherits TYPE] { [feature; ]∗}

feature := ID( [ formal [, formal]∗] ) : TYPE { expr }

|ID : TYPE [ &lt;- expr ]

formal ::= ID : TYPE

expr ::= ID&lt;-expr

expr[@TYPE].ID( [ expr [, expr]∗] )

ID( [ expr [, expr]∗] )

if expr then expr else expr ﬁ

while expr loop expr pool

{$[expr,]^{+}\}$

let ID : TYPE [&lt;-expr][,ID:TYPE [&lt;-expr .15]]∗in expr

case expr of $[ID:TYPE=>expr,]^{+}esac$

new TYPE

isvoid expr

expr+expr

expr-expr

expr ∗expr

expr / expr

˜expr

expr&lt;expr

expr&lt;=expr

expr=expr

not expr

(expr)

ID

integer

string

true

false

Figure 1: Cool syntax.

Figure 1 provides a speciﬁcation of Cool syntax. The speciﬁcation is not in pure Backus-Naur Form (BNF); for convenience, we also use some regular expression notation. Speciﬁcally, A∗means zero or more A’s in succession; $A^{+}$ means one or more A’s. Items in square brackets [. . .] are optional. Double brackets [[ ]] are not part of Cool; they are used in the grammar as a meta-symbol to show association of grammar symbols (e.g.$a[bc]^{+}$ means a followed by one or more bc pairs).

### 11.1Precedence

The precedence of inﬁx binary and preﬁx unary operations, from highest to lowest, is given by the following table:

.

@

\~

isvoid

* /

+ 

&lt;=&lt;=

not

&lt;-

All binary operations are left-associative, with the exception of assignment, which is right-associative,and the three comparison operations, which do not associate.

## 12Type Rules

This section formally deﬁnes the type rules of Cool. The type rules deﬁne the type of every Cool expression in a given context. The context is the type environment, which describes the type of every unbound identiﬁer appearing in an expression. The type environment is described in Section 12.1. Section 12.2gives the type rules.

### 12.1Type Environments

To a ﬁrst approximation, type checking in Cool can be thought of as a bottom-up algorithm: the type of an expression e is computed from the (previously computed) types of e’s subexpressions. For example,an integer 1 has type Int; there are no subexpressions in this case. As another example, if $e_{n}$ has type X, then the expression $\{e_{1}j\cdots ;e_{n};\}$ has type X.

A complication arises in the case of an expression v, where v is an object identiﬁer. It is not possible to say what the type of v is in a strictly bottom-up algorithm; we need to know the type declared for v in the larger expression. Such a declaration must exist for every object identiﬁer in valid Cool programs.

To capture information about the types of identiﬁers, we use a type environment. The environment consists of three parts: a method environment M, an object environment O, and the name of the current class in which the expression appears. The method environment and object environment are both functions (also called mappings). The object environment is a function of the form

O(v)=T

which assigns the type T to object identiﬁer v. The method environment is more complex; it is a function of the form

$M(C,f)=(T_{1},\cdots ,T_{n-1},T_{n})$

where C is a class name (a type), f is a method name, and $t_{1},\cdots ,t_{n}$ are types. The tuple of types is the signature of the method. The interpretation of signatures is that in class C the method f has formal parameters of types$(t_{1},\cdots ,t_{n-1})$—in that order—and a return type $t_{n}$.

Two mappings are required instead of one because object names and method names do not clash—i.e.,there may be a method and an object identiﬁer of the same name.

The third component of the type environment is the name of the current class, which is needed for type rules involving SELF TYPE.

Every expression e is type checked in a type environment; the subexpressions of e may be type checked in the same environment or, if e introduces a new object identiﬁer, in a modiﬁed environment.For example, consider the expression

let c : Int &lt;- 33 in

...

The let expression introduces a new variable c with type Int. Let O be the object component of the type environment for the let. Then the body of the let is type checked in the object type environment

O[Int/c]

where the notation O[T/c] is deﬁned as follows:

O[T/c](c)=T

O[T/c](d)=O(d) ifd≠c

### 12.2Type Checking Rules

The general form a type checking rule is:

$\frac {i}{O,M,C\vert -e:T}$

The rule should be read: In the type environment for objects O, methods M, and containing class C,the expression e has type T. The dots above the horizontal bar stand for other statements about the types of sub-expressions of e. These other statements are hypotheses of the rule; if the hypotheses are satisﬁed, then the statement below the bar is true. In the conclusion, the “turnstyle” (“⊢”) separatescontext (O,M,C) from statement (e : T).

The rule for object identiﬁers is simply that if the environment assigns an identiﬁer Id type T, then Id has type T.

$\frac {O(Id)=T}{O,M,C\vert -Id:T}$ [Var]

The rule for assignment to a variable is more complex:

O(Id)=T

$O,M,C\vert -e_{1}:T^{\prime }$

$\frac {T^{\prime }\leq T}{O,M,C\vert -Id<-e_{1}:}$:T [ASSIGN]

Note that this type rule-as well as others-use the conformance relation≤(see Section 3.2). The rule says that the assigned expression e1 must have a type T' that conforms to the type T of the identifier Id in the type environment. The type of the whole expression is T'.

The type rules for constants are all easy:

O,M,C|-true:Bool [True]

O,M,C|-false:Bod [False]

i is an integer constant [Int]O,M,C|-i:Int

$\frac {\sin \alpha \sin g\cos t\sin t}{O,M,C\vert -s:Str\in g}$ [String]

There are two cases for new, one for new SELF_TYPE and one for any other form:

$T^{\prime }=\{\begin{matrix}SELF-TYEC&ifT=SELF.TYFE\\T&Othervise\\O,M,C+newT:T^{\prime }$

[New]

Dispatch expressions are the most complex to type check.

$O,M,C\vert -e_{0}:T_{0}$

$O,M,C\vert -e_{1}:T_{1}$

·..

$O,M,C\vert -e_{n}:T_{n}$

$T_{0}^{\prime }=\{\begin{matrix}C&ifT_{0}=SELF_{-}TYPE_{C}\\T_{0}otherwise\end{matrix}$

$M(T_{0}^{\prime },f)=(T_{1}^{\prime },\cdots ,T_{n}^{\prime },T_{n+1}^{\prime })$

$T_{i}\leq T_{i}^{\prime }$ 1≤i≤n

$T_{n+1}=\{\begin{matrix}T_{0}ifT_{n+1}^{\prime }=SELF_{-}TYPE\\T_{n+1}^{\prime }&otherwise\\O,M,C+e_{0}.f(e_{1},\cdots ,e_{n}):T_{n+1$

[Dispatch]

$O,M,C\vert -e_{0}:T_{0}$

$O,M,C\vert -e_{1}:T_{1}$

$O,M,C\vert -e_{n}:T_{n}$

$T_{0}\leq T$

$M(T,f)=(T_{1}^{\prime },\cdots ,T_{n}^{\prime },T_{n+1}^{\prime })$

$T_{i}\leq T_{i}^{\prime }$ 1≤i≤n

$T_{n+1}=\{\begin{matrix}T_{0}ifT_{n+1}^{\prime }=SELF_{-}TYFE\\T_{n+1}^{\prime }&otherwise\end{matrix}$

[StaticDispatch]

To type check a dispatch, each of the subexpressions must first be type checked.The type$T_{0}$of$e_{0}$determines which declaration of the method f is used. The argument types of the dispatch must conform to the declared argument types. Note that the type of the result of the dispatch is either the declared return type or$T_{0}$in the case that the declared return type is SELF_TYPE. The only difference in type checking a static dispatch is that the class T of the method f is given in the dispatch,and the type To must conform to T.

The type checking rules for if and {-} expressions are straightforward. See Section 7.5 for the definition of the L operation.

$O,M,C\vert -e_{1}:Bod$ $O,M,C\vert -e_{2}:T_{2}$ $\frac {O,M,C+e_{3}:T_{3}}{O,M,CHfe_{1}thene_{2}else_{3}f_{2}\cup T_{3}}$

[If]

$O,M,C\vert -e_{1}:T_{1}$

$O,M,C\vert -e_{2}:T_{2}$

$\frac {O,M,C\vert -en:T_{n}}{O,M,C\vert -\{e_{1};e_{2};\cdots e_{n}\}:T_{n}}$ [Sequence]

The let rule has some interesting aspects.

$T_{0}^{\prime }=\{\begin{matrix}SELF_{-}TYPE_{C}\\T_{0}\end{matrix}$ 、$T_{0}=SELF_{-}TYPE$ otherwise

$O,M,C\vert -e_{1}:T_{1}$

$T_{1}\leq T_{0}^{\prime }$

$\frac {O[T_{0}^{\prime }/x],M,C+e_{2}:T_{2}}{0,M,C+1etx:T_{0}-e_{1}\in e_{2}:T_{2}}$ [Let-Init]

First,the initialization$e_{1}$is type checked in an environment without a new definition for x. Thus, the variable æ cannot be used in ei unless it already has a definition in an outer scope. Second, the body e2is type checked in the environment O extended with the typing$x:T_{0}^{\prime }$.Third,note that the type of x may be SELF_TYPE.

$T_{0}^{\prime }=\{\begin{matrix}SELF_{-}TYPE_{C}\dot {I}_{0}=SELF_{-}TYPE\\T_{0}otherwise\end{matrix}$

$\frac {O[T_{0}^{\prime }/x],M,C\vert -e_{1}:T_{1}}{O,M,C\vert -1etx:T_{0}\in x:T_{1}:T_{1}}$ [Let-No-Init]

The rule for let with no initialization simply omits the conformance requirement. We give type rules only for a let with a single variable. Typing a multiple let

$letx_{1}:T_{1}[-e_{1}],x_{2}:T_{2}[-e2]_{\cdots }x_{n}:T_{n}[-e_{n}]$in e

is defined to be the same as typing

let$x_{1}:T_{1}$ $[\leftarrow e_{1}]$in(let$x_{2}:T_{2}[-e_{2}],\cdots ,x_{n}:T_{n}[-e_{n}]\in e)$

$O,M,C\vert -e_{0}:T_{0}$

$O[T_{1}/x_{1}],M,C\vert -e_{1}:T_{1}^{\prime }$

·..

$O,M,C\vert -\cos eet_{0}0f$ $x_{1}:$ $T_{1}\Rightarrow e_{1};\cdots x_{n}:$ $:T_{n}=$ $e_{n}$ $\csc c:\vert 1\leq i\leq nT_{i}^{\prime }$ [Case]$O[T_{n}/x_{n}],M,C\vert -e_{n}:T_{n}^{\prime }$

Each branch of a case is type checked in an environment where variable$x_{i}$has type$T_{i}.$.The type of the entire case is the join of the types of its branches. The variables declared on each branch of a case must all have distinct types.

$O,M,C\vert -e_{1}:Bod$

$O,M,C\vert -e_{2}:T_{2}$

O,M,C,-$e_{1}$ $\log e_{2}pod:Object$ [Loop]

The predicate of a loop must have type Bool; the type of the entire loop is always Object. An isvoid test has type Bool:

$\frac {O,M,C\vert -e_{1}:T_{1}}{O,M,C\vert -isvoide_{1}:Bool}$ [Isvoid]

With the exception of the rule for equality, the type checking rules for the primitive logical, compar-ison, and arithmetic operations are easy.

$\frac {O,M,C\vert -e_{1}:Bool}{O,M,C\vert --e_{1}:Bool}$ [Not]

[Compare]

$O,M,C\vert -e_{1}:Int$ $O,M,C\vert -e_{2}:Int$ $\frac {op\in \{<,\leq \}}{O,M,C\vert -e_{1}ope_{2}:Bod}$ $\frac {O,M,C\vert -e_{1}:Int}{O,M,C\vert --e_{1}:Int}$ $O,M,C\vert -e_{1}:Int$ $O,M,C\vert -e_{2}:Int$ $\frac {op\in \{*,+,-,}{O,M,C\vert -e_{1}Op}$ $\frac {,-,\beta }{ope_{2}:Int}$

[Neg]

[Arith]

The wrinkle in the rule for equality is that any types may be freely compared except Int, String and Bool, which may only be compared with objects of the same type.

$O,M,C\vert -e_{1}:T_{1}$ $O,M,C\vert -e_{2}:T_{2}$ $T_{1}\in \{Int,Str\in ing,Bool\}VT_{2}\in \{Int,Str\in ng,Bool\}\rightarrow T_{1}=T_{2}$

[Equal]

The final cases are type checking rules for attributes and methods. For a classC,let the object environment$O_{C}$give the types of all attributes of C (including any inherited attributes).More formally,if x is an attribute (inherited or not) of C, and the declaration of x is:T-,then

$Oc(x)=\{\begin{matrix}SELF_{-}TYPE_{C}\dot {I}T=SELF_{-}TYPE\\T&otheTwise\end{matrix}$

The method environment M is global to the entire program and deﬁnes for every class C the signatures of all of the methods of C (including any inherited methods).

The two rules for type checking attribute deﬁninitions are similar the rules for let. The essential diﬀerence is that attributes are visible within their initialization expressions. Note that self is bound in the initialization.

$Oc(x)=T_{0}$

Oc|$[SELF_{-}TYPEC/sel],M,C\vert -e_{1}:$

$T_{1}\leq T_{0}\\O_{C},M,C+x:T_{0}\leftarrow e_{1};$ [Attr-Init]

$\frac {Oc(x)=T}{Oc,M,C\vert -x:T;}$ [Attr-No-Init]

The rule for typing methods checks the body of the method in an environment where$O_{C}$ is extended with bindings for the formal parameters and self. The type of the method body must conform to the declared return type.

$M(C,f)=(T_{1},\cdots ,T_{n},T_{0})$

$Oc[SELF.TYEc/sel射[T_{1}/x_{1}]\cdots [T_{n}/x_{n}],M,C\vert -e:T_{0}^{\prime }$

$T_{0}^{\prime }\leq \{\begin{matrix}SE正F_{-}YPE_{C}&ifT_{0}=SELF_{-}TYPE\\OG,M_{1}C+f(x_{1}:T_{1},\cdots ,x_{n}:T_{n}):T_{0}\{e\};$

[Method]

## 13Operational Semantics

This section contains a mostly formal presentation of the operational semantics for the Cool language. The operational semantics deﬁne for every Cool expression what value it should produce in a given context.The context has three components: an environment, a store, and a self object. These components are described in the next section. Section 13.2 deﬁnes the syntax used to refer to Cool objects, and Section 13.3 deﬁnes the syntax used to refer to class deﬁnitions.

Keep in mind that a formal semantics is a speciﬁcation only—it does not describe an implementation.The purpose of presenting the formal semantics is to make clear all the details of the behavior of Cool expressions. How this behavior is implemented is another matter.

### 13.1Environment and the Store

Before we can present a semantics for Cool we need a number of concepts and a considerable amount of notation. An environment is a mapping of variable identiﬁers to locations. Intuitively, an environment tells us for a given identiﬁer the address of the memory location where that identiﬁer’s value is stored.For a given expression, the environment must assign a location to all identiﬁers to which the expression may refer. For the expression, e.g,a+b, we need an environment that maps a to some location and b to some location. We’ll use the following syntax to describe environments, which is very similar to the syntax of type assumptions used in Section 12.

$E=[a:l_{1},b:l_{2}]$

This environment maps a to location $l_{1}$, and b to location $l_{2}$

The second component of the context for the evaluation of an expression is the store (memory). The store maps locations to values, where values in Cool are just objects. Intuitively, a store tells us what value is stored in a given memory location. For the moment, assume all values are integers. A store is similar to an environment:

$S=[l_{1}\rightarrow 55,l_{2}\rightarrow 77]$

This store maps location $l_{1}$ to value 55 and location $l_{2}$ to value 77.

Given an environment and a store, the value of an identiﬁer can be found by ﬁrst looking up the location that the identiﬁer maps to in the environment and then looking up the location in the store.

$E(a)=l_{1}$

$S(l_{1})=55$

Together, the environment and the store deﬁne the execution state at a particular step of the evaluation of a Cool expression. The double indirection from identiﬁers to locations to values allows us to model variables. Consider what happens if the value 99 is assigned variable a in the environment and store deﬁned above. Assigning to a variable means changing the value to which it refers but not its location.To perform the assignment, we look up the location for a in the environment E and then change the mapping for the obtained location to the new value, giving a new store $S^{\prime }.$

$E(a)=l_{1}\\S^{\prime }=S[99/l_{1}]$

The syntax S[v/l] denotes a new store that is identical to the store S, except that $S^{\prime }$ maps location l to value v. For all locations l′ where$l^{\prime }\neq l,$ we still have $S^{\prime }(l^{\prime })=S(l^{\prime }).$.

The store models the contents of memory of the computer during program execution. Assigning to a variable modiﬁes the store.

There are also situations in which the environment is modiﬁed. Consider the following Cool fragment:

let c : Int &lt;- 33 in

c

When evaluating this expression, we must introduce the new identiﬁer c into the environment before evaluating the body of the let. If the current environment and state are E and S, then we create a new environment E′ and a new store S′ deﬁned by:

=newloc(S)

$E^{\prime }=E[l_{c}/c]$

$S^{\prime }=S[33/lc]$

The ﬁrst step is to allocate a location for the variable c. The location should be fresh, meaning that the current store does not have a mapping for it. The function newloc() applied to a store gives us an unused location in that store. We then create a new environment$E^{\prime }$, which maps c to $l_{c}$ but also contains all of the mappings of E for identiﬁers other than c. Note that if c already has a mapping in E, the new environment E′ hides this old mapping. We must also update the store to map the new location to a value. In this case lc maps to the value 33, which is the initial value for c as deﬁned by the let-expression.

The example in this subsection oversimpliﬁes Cool environments and stores a bit, because simple integers are not Cool values. Even integers are full-ﬂedged objects in Cool.

### 13.2Syntax for Cool Objects

Every Cool value is an object. Objects contain a list of named attributes, a bit like records in C. In addition, each object belongs to a class. We use the following syntax for values in Cool:

$v=X(a_{1}=l_{1},a_{2}=l_{2,\cdots ,a_{n}=l_{n})$

Read the syntax as follows: The value v is a member of class X containing the attributes $a_{1},\cdots ,a_{n}$ whose locations are$l_{1},\cdots ,l_{n}$. Note that the attributes have an associated location. Intuitively this means that there is some space in memory reserved for each attribute.

For base objects of Cool (i.e., Ints, Strings, and Bools) we use a special case of the above syntax.Base objects have a class name, but their attributes are not like attributes of normal classes, because they cannot be modiﬁed. Therefore, we describe base objects using the following syntax:

Int(5)

Bool(true)

$Str\in g(4,^{n}Cool^{n})$

For Ints and Bools, the meaning is obvious. Strings contain two parts, the length and the actual sequence of ASCII characters.

### 13.3Class deﬁnitions

In the rules presented in the next section, we need a way to refer to the deﬁnitions of attributes and methods for classes. Suppose we have the following Cool class deﬁnition:

class B {

s : String &lt;- "Hello";

g (y:String) : Int {

y.concat(s)

};

f (x:Int) : Int {

x+1

};

};

class A inherits B {

a : Int;

b:B&lt;-newB;

f(x:Int) : Int {

x+a

};

};

Two mappings, called class and implementation, are associated with class deﬁnitions.The class mapping is used to get the attributes, as well as their types and initializations, of a particular class:

$dass(A)=(s:Str\in g-^{n}Hell_{0}^{n},a:Int-0,b:B-newB)$

Note that the information for class A contains everything that it inherited from class B, as well as its own deﬁnitions. If B had inherited other attributes, those attributes would also appear in the information for A. The attributes are listed in the order they are inherited and then in source order: all the attributes from the greatest ancestor are listed ﬁrst in the order in which they textually appear, then the attributes of the next greatest ancestor, and so on, on down to the attributes deﬁned in the particular class. We rely on this order in describing how new objects are initialized.

The general form of a class mapping is:

$class(X)=(a_{1}:T_{1}\leftarrow e_{1},\cdots ,a_{n}:T_{n}-e_{n})$

Note that every attribute has an initializing expression, even if the Cool program does not specify one for each attribute. The default initialization for a variable or attribute is the default of its type. The default of Int is 0, the default of String is "", the default of Bool is false, and the default of any other type is void.4 The default of type T is written $D_{T}$.

The implementation mapping gives information about the methods of a class. For the above example,implementation of A is deﬁned as follows:

implementation(A,f)=(x,x+a)

implementation(A,g)=((y,y,cotxdt(s))

In general, for a class X and a method m,

$implcmentation(X,m)=(x_{1},x_{2},\cdots ,x_{n},e_{body})$

speciﬁes that method m when invoked from class X, has formal parameters $x_{1},\cdots ,x_{n},$ and the body of the method is expression $e_{body}.$

### 13.4Operational Rules

Equipped with environments, stores, objects, and class deﬁnitions, we can now attack the operational semantics for Cool. The operational semantics is described by rules similar to the rules used in type checking. The general form of the rules is:

.$\overset\frown{SO,S,E\vert -e_{1}:v,S^{\prime }}$

The rule should be read as: In the context where self is the object so, the store is S, and the environment is E, the expression e1 evaluates to object v and the new store is $S^{\prime }$. The dots above the horizontal bar stand for other statements about the evaluation of sub-expressions of $e_{1}$.

Besides an environment and a store, the evaluation context contains a self object so. The self object is just the object to which the identiﬁer self refers if self appears in the expression. We do not place self in the environment and store because self is not a variable—it cannot be assigned to. Note that the rules specify a new store after the evaluation of an expression. The new store contains all changes to memory resulting as side eﬀects of evaluating expression e1.

4A tiny point: We are allowing void to be used as an expression here. There is no expression for void available to Cool programmers.

The rest of this section presents and briefly discusses each of the operational rules. A few cases are not covered; these are discussed at the end of the section.

$so,S_{1},E\vert -e_{1}:v_{1},S_{2}$ $E(Id)=l_{1}$ $\frac {S_{3}=S_{2}[v_{1}/l_{1}]}{s_{0},S_{1},E\vert -Id<e_{1}:v_{1},S_{3}}$

[Assign]

An assignment first evaluates the expression on the right-hand side, yielding a value v1. This value is stored in memory at the address for the identifier.

The rules for identifier references, self, and constants are straightforward:

E(Id)=l

$\frac {S(l)=v}{so,S,E\vert -Id:v,S}$ [Var]

$\overline {so,S,E\vert -self:so,S}$ [Self]

so,S,E|-true:Bool(true),S [True]

SO,S,E|-falSe:Bod(false),S [False]

$\frac {iisan\inf tegercons\tan t}{so,S,E+i\cdot Int(i),S}$ [Int]

s is a string constant l=length(s)$\frac {l=length(s)}{so,S,E\vert -s:Str\in g(l,s),S}$ [String]

A new expression is more complicated than one might expect:

$T_{0}=\{\begin{matrix}XifT=SELF_{-}TYPEand\\T&otherwise\end{matrix}$80=X(...)

$class(T_{0})=(a_{1}:T_{1}\leftarrow e_{1},$ $T_{n}\leftarrow e_{n})$

$l_{i}=newloc(S_{1})$,fori=1...n and each$l_{i}$is diistinct

$v_{1}=T_{0}(a_{1}=l_{1},\cdots ,a_{n}=l_{n})$

$S_{2}=S_{1}[D_{T1}/l_{1,\cdots ,D_{Tn}/l_{n}]$

$v_{1},S_{2},[a_{1}:l_{1},\cdots ,a_{n}:l_{n}]-\{a_{1}-e_{1};\cdots ;a_{n}-e_{n}\}:v_{2},S_{3}$ [New]

The tricky thing in a new expression is to initialize the attributes in the right order. Note also that,during initialization, attributes are bound to the default of the appropriate class.

$so,S_{1},E\vert -e_{1}:v_{1},S_{2}$

$so,S_{2},E\vert -e_{2}:v_{2},S_{3}$

...

$so,S_{n},E\vert -e_{n}:v_{n},S_{n+1}$

$so,S_{n+1},E\vert -e_{0}:v_{0},S_{n+2}$

$v_{0}=X(a_{1}=l_{a_{1},\cdots ,a_{m}}=l_{a_{m}})$

implementation(X,$,f)=(x_{1},\cdots ,x_{n},e_{n+1})$

$l_{x_{i}}=newloc(S_{n+2}),$,fori=1....nand each$l_{x_{i}}$is distinct

$S_{n+3}=S_{n+2}[v_{1}/l_{x1},\cdots ,v_{n}/l_{xn}]$

$\frac {v_{0},S_{n+3}\vert (a_{1}:l_{a1},\cdots ,a_{nn}:l_{a_{n},x_{1},\cdots ,x_{n}}{s_{0},S_{1},E_{1},EFe_{1},\cdots ,e_{n}):\frac$ [Dispatch]

$so,S_{1},E\vert -e_{1}:v_{1},S_{2}$

80,$S_{2},$ $E\vert -e_{2}:v_{2},S_{3}$

··.

$so,S_{n},E\vert -e_{n}:v_{n},S_{n+1}$

$so,S_{n+1},E\vert -e_{0}:v_{0},S_{n+2}$

$v_{0}=X(a_{1}=l_{a_{1},\cdots ,a_{m}}=l_{a_{m}})$

implementation$(T,f)=(x_{1},\cdots ,x_{n},e_{n+1})$

$l_{x_{i}}=newloc(S_{n+2})$,fori=1..nand each$l_{x_{i}}$is distinct

$S_{n+3}=S_{n+2}[v_{1}/l_{x_{1}},\cdots ,v_{n}/l_{x_{n}}]$

$t_{0},S_{n+3},(a_{1}:l_{21},\cdots ,a_{nn}:l_{an},x_{1}:I_{21},\cdots ,x_{n}:l_{n}\vert \vert -e_{n-1}\vert \cdot v_{n+1}:v_{n+1},S$ [StaticDispatch]

The two dispatch rules do what one would expect. The arguments are evaluated and saved. Next,the expression on the left-hand side of the “” is evaluated. In a normal dispatch, the class of this expression is used to determine the method to invoke; otherwise the class is specified in the dispatch itself.

$80,S_{1},E\vert -e_{1}:Bool(true),S_{2}$

$\frac {s_{1}S_{2},E_{1}-e_{2}:v_{2},S_{3}}{s_{0},S_{1},E_{1}+ife_{1}thene_{2}elsee_{3}f_{2}:v_{2},S_{3}}$ [If-True]

SO,S1,E|-e1:Bool(false),S2

$\frac {s_{2}S_{2},EFe_{3}:v_{3},S_{3}}{s_{0},S_{1},E_{1}+fe_{1}+hene_{2}e_{3}f_{1}:v_{3},S_{3}}$ [If-False]

There are no surprises in the if-then-else rules. Note that value of the predicate is a Bool object, not a boolean.

$so,S_{1},E\vert -e_{1}:v_{1},S_{2}$

$so,S_{2},E\vert -e_{2}:v_{2},S_{3}$

$\frac {s_{0},S_{n},E(-e_{n},v_{n},S_{n+1}}{s_{0},S_{1},E+\{e_{1};e_{2}\cdots ;e_{n}\}:v_{n},S_{n+1}}$ [Sequence]

Blocks are evaluated from the first expression to the last expression, in order. The result is the result of the last expression.

$so,S_{1},E\vert -e_{1}:v_{1},S_{2}$ $l_{1}=newloc(S_{2})$ $S_{3}=S_{2}[v_{1}/l_{1}]$ $E^{\prime }=E[l_{1}/Id]$ $\frac {s_{3}S_{3}E^{\prime }-e_{2}:v_{2},S_{4}}{s_{0},S_{1},E+1etId:T_{1}<e_{1}\sin e_{2}:v_{2},S_{4}}$

[Let]

A let evaluates any initialization code, assigns the result to the variable at a fresh location, and evaluates the body of the let. (If there is no initialization, the variable is initialized to the default value of$T_{1}.)$We give the operational semantics only for the case of let with a single variable. The semantics of a multiple 1et

let$x_{1}:T_{1}<e_{1},x_{2}:T_{2}\leftarrow e_{3}\cdots ,x_{n}:T_{n}<e_{n}\in e$

is defined to be the same as

let$x_{1}:T_{1}\leftarrow e_{1}$in$(letx_{2}:T_{2}\leftarrow e_{2},\cdots ,x_{n}:T_{n}\in e_{n}\in e)$

$80,S_{1},E\vert -e_{0}:v_{2},S_{2}$

$v_{0}=X(..)$

Ti=closest ancestor of X in$\{T_{1},\cdots ,T_{n}\}$

$l_{0}=newloc(S_{2})$

$S_{3}=S_{2}[v_{0}/l_{0}]$

$E^{\prime }=E[l_{0}/Id_{i}]$

$\overline {so,S_{1},E\vert -\cos e}$ $e_{0}$ $ofId_{1}:T_{1}\Rightarrow e_{1};\cdots ;$ $Id_{n}:$ $:T_{n}\Rightarrow e_{n};$ $esac:v_{1},S_{4}$ [Case]$SO,S3,E^{\prime }\vert -e_{i}:v1,S_{4}$

Note that the case rule requires that the class hierarchy be available in some form at runtime, so that the correct branch of the case can be selected. This rule is otherwise straightforward.

$so,S_{1},E\vert -e_{1}:Bod(true),S_{2}$

$so,S_{2},E\vert -e_{2}:v2,S_{3}$

$\frac {so,S_{3}E(-whileee_{1}loope_{2}pool:vod,S_{4}}{so,S_{1},E_{1}-whileee_{1}loope_{2}pod:vod,S_{4}}$ [Loop-True]

$\frac {so,S_{1},E+e_{1}:Bool(false),S_{2}}{so,S_{1},E\vert -whilee1\log e_{2}pood:void,S_{2}}$ [Loop-False]

There are two rules for while: one for the case where the predicate is true and one for the case where the predicate is false. Both cases are straightforward. The two rules for isvoid are also straightforward:

$\frac {so,S_{1},E\vert -e_{1}:void,S_{2}}{so,S_{1},E\vert -isvoide1:Bood(true),S_{2}}$ [IsVoid-True]

$\frac {s_{0},S_{1},E\vert -e_{1}:X(\cdots ),S_{2}}{so,S_{1},E\vert -isvoide_{1}:Bool(false),S_{2}}$ [IsVoid-False]

The remainder of the rules are for the primitive arithmetic, logical, and comparison operations except equality. These are all easy rules.

$so,S_{1},E\vert -e_{1}:Bool(b),S_{2}$

$\frac {v_{1}=Bool(-b)}{so_{2},S_{1},E+note_{1}:v_{1},S_{2}}$ [Not]

$so,S_{1},E\vert -e_{1}:Int(i_{1}),S_{2}$

$so,S_{2},E\vert -e_{2}:Int(i_{2}),S_{3}$

Op∈{≤,&lt;}

$v_{1}=\{\begin{matrix}Bod(true),ifi_{1}opi_{2}\\Bool(false),othervise\\so.S_{1}.E+e_{1}ope_{2}:v_{1},S_{3}}\end{matrix}$

[Comp]

$so,S_{1},E\vert -e_{1}:Int(i_{1}),S_{2}$

$\frac {v_{1}=Int(-i_{1})}{so,S_{1},E\vert -e_{1}:v_{1},S_{2}}$ [Neg]

$so,S_{1},E\vert -e_{1}:Int(i_{1}),S_{2}$

$so,S_{2},E\vert -e_{2}$ $:Int(i_{2}),S_{3}$

Op∈{*,+,-,/}

$\frac {v_{1}=Int(i_{1},opi_{2})}{so_{1}S_{1},E+e_{1}ope_{2}:v_{1},S_{3}}$ [Arith]

Cool Ints are 32-bit two's complement signed integers; the arithmetic operations are defined accordingly.

The notation and rules given above are not powerful enough to describe how objects are tested for equality,or how runtime exceptions are handled. For these cases we resort to an English description.

In$e_{1}=e_{2}$,first ei is evaluated and then e2 is evaluated. The two objects are compared for equality by first comparing their pointers (addresses). If they are the same,the objects are equal. The value void is not equal to any object except itself. If the two objects are of type String, Bool, or Int, their respective contents are compared.

In addition, the operational rules do not specify what happens in the event of a runtime error.Execution aborts when a runtime error occurs. The following list specifies all possible runtime errors.

1.A dispatch (static or dynamic) on void.

2.A case on void.

3.Execution of a case statement without a matching branch.

4.Division by zero.

5. Substring out of range.

6. Heap overflow.

Finally, the rules given above do not explain the execution behaviour for dispatches to primitive methods defined in the Object, I0, or String classes. Descriptions of these primitive methods are given in Sections 8.3-8.5.

Cool is based on Sather164, which is itself based on the language Sather. Portions of this document were cribbed from the Sather164 manual; in turn, portions of the Sather164 manual are based on Sather documentation written by Stephen M. Omohundro.

A number people have contributed to the design and implementation of Cool, including Manuel F¨ahndrich, David Gay, Douglas Hauge, Megan Jacoby, Tendo Kayiira, Carleton Miyamoto, and Michael Stoddart. Joe Darcy updated Cool to the current version.



