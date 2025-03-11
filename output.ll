%Object = type { i8* }
%IO = type { i8* }
%Int = type { i8* }
%String = type { i8* }
%Bool = type { i8* }
%Node = type { i8*, i32, %Node* }
%LinkedList = type { i8*, %Node*, %Node*, i32 }
%Main = type { i8* }

@.str.empty = constant [1 x i8] c"\00"
@vtable.Object = global [3 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.IO = global [7 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%IO*)* @IO.in_int to i8*), i8* bitcast (i8* (%IO*)* @IO.in_string to i8*), i8* bitcast (%IO* (%IO*, i32)* @IO.out_int to i8*), i8* bitcast (%IO* (%IO*, i8*)* @IO.out_string to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Int = global [3 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.String = global [6 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (i8* (%String*, i8*)* @String.concat to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%String*)* @String.length to i8*), i8* bitcast (i8* (%String*, i32, i32)* @String.substr to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Bool = global [3 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Node = global [7 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%Node*)* @Node.getData to i8*), i8* bitcast (%Node* (%Node*)* @Node.getNext to i8*), i8* bitcast (%Node* (%Node*, i32, %Node*)* @Node.init to i8*), i8* bitcast (%Node* (%Node*, %Node*)* @Node.setNext to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.LinkedList = global [10 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%LinkedList* (%LinkedList*, i32)* @LinkedList.addFirst to i8*), i8* bitcast (%LinkedList* (%LinkedList*, i32)* @LinkedList.addLast to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%LinkedList*, i32)* @LinkedList.get to i8*), i8* bitcast (%LinkedList* (%LinkedList*)* @LinkedList.init to i8*), i8* bitcast (i1 (%LinkedList*)* @LinkedList.isEmpty to i8*), i8* bitcast (i32 (%LinkedList*)* @LinkedList.removeFirst to i8*), i8* bitcast (i32 (%LinkedList*)* @LinkedList.size to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Main = global [8 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%IO*)* @IO.in_int to i8*), i8* bitcast (i8* (%IO*)* @IO.in_string to i8*), i8* bitcast (%Object* (%Main*)* @Main.main to i8*), i8* bitcast (%IO* (%IO*, i32)* @IO.out_int to i8*), i8* bitcast (%IO* (%IO*, i8*)* @IO.out_string to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@.str9 = internal constant [41 x i8] c"Error: Cannot remove from an empty list\0A\00"
@.str10 = internal constant [45 x i8] c"Error: Index out of bounds (negative index)\0A\00"
@.str11 = internal constant [46 x i8] c"Error: Index out of bounds (index too large)\0A\00"
@.str12 = internal constant [26 x i8] c"Created a new LinkedList\0A\00"
@.str13 = internal constant [12 x i8] c"List size: \00"
@.str14 = internal constant [2 x i8] c"\0A\00"
@.str15 = internal constant [16 x i8] c"First element: \00"
@.str16 = internal constant [2 x i8] c"\0A\00"
@.str17 = internal constant [17 x i8] c"Second element: \00"
@.str18 = internal constant [2 x i8] c"\0A\00"
@.str19 = internal constant [16 x i8] c"Third element: \00"
@.str20 = internal constant [2 x i8] c"\0A\00"
@.str21 = internal constant [24 x i8] c"Removed first element: \00"
@.str22 = internal constant [2 x i8] c"\0A\00"
@.str23 = internal constant [16 x i8] c"New list size: \00"
@.str24 = internal constant [2 x i8] c"\0A\00"
@.str25 = internal constant [20 x i8] c"New first element: \00"
@.str26 = internal constant [2 x i8] c"\0A\00"
@.str27 = internal constant [15 x i8] c"List is empty\0A\00"
@.str28 = internal constant [19 x i8] c"List is not empty\0A\00"
@.str29 = internal constant [15 x i8] c"List is empty\0A\00"
@.str30 = internal constant [19 x i8] c"List is not empty\0A\00"
@.str.Object = constant [7 x i8] c"Object\00"
@.str.fmt = global [3 x i8] c"%s\00"
@.str.fmt.int = constant [3 x i8] c"%d\00"
@.str.scanf_s_fmt = constant [3 x i8] c"%s\00"
@.str.scanf_d_fmt = constant [3 x i8] c"%d\00"
@.str.substr_error = constant [39 x i8] c"Runtime error: substring out of range\0A\00"

define %Object* @Object.abort(%Object* %self) {
entry:
	call void @exit(i32 1)
	unreachable
}

define i8* @Object.type_name(%Object* %self) {
entry:
	ret i8* getelementptr ([7 x i8], [7 x i8]* @.str.Object, i32 0, i32 0)
}

define %Object* @Object.copy(%Object* %self) {
entry:
	ret %Object* %self
}

define %IO* @IO.out_string(%IO* %self, i8* %x) {
entry:
	%0 = call i32 (i8*, ...) @printf(i8* getelementptr ([3 x i8], [3 x i8]* @.str.fmt, i32 0, i32 0), i8* %x)
	ret %IO* %self
}

define %IO* @IO.out_int(%IO* %self, i32 %x) {
entry:
	%0 = call i32 (i8*, ...) @printf(i8* getelementptr ([3 x i8], [3 x i8]* @.str.fmt.int, i32 0, i32 0), i32 %x)
	ret %IO* %self
}

define i8* @IO.in_string(%IO* %self) {
entry:
	%0 = alloca [1024 x i8]
	%1 = bitcast [3 x i8]* @.str.scanf_s_fmt to i8*
	%2 = bitcast [1024 x i8]* %0 to i8*
	%3 = call i32 (i8*, ...) @scanf(i8* %1, i8* %2)
	%4 = bitcast [1024 x i8]* %0 to i8*
	%5 = call i32 @strlen(i8* %4)
	%6 = add i32 %5, 1
	%7 = zext i32 %6 to i64
	%8 = call i8* @malloc(i64 %7)
	%9 = call i8* @strcpy(i8* %8, i8* %4)
	ret i8* %8
}

define i32 @IO.in_int(%IO* %self) {
entry:
	%0 = alloca i32
	%1 = bitcast [3 x i8]* @.str.scanf_d_fmt to i8*
	%2 = call i32 (i8*, ...) @scanf(i8* %1, i32* %0)
	%3 = load i32, i32* %0
	ret i32 %3
}

define i32 @String.length(%String* %self) {
entry:
	%0 = ptrtoint %String* %self to i64
	%1 = bitcast %String* %self to i8*
	%2 = ptrtoint i8* %1 to i64
	%3 = icmp eq i64 %0, %2
	br i1 %3, label %raw_string, label %struct_string

raw_string:
	%4 = bitcast %String* %self to i8*
	%5 = call i32 @strlen(i8* %4)
	ret i32 %5

struct_string:
	%6 = bitcast %String* %self to i8**
	%7 = load i8*, i8** %6
	%8 = call i32 @strlen(i8* %7)
	ret i32 %8
}

define i8* @String.concat(%String* %self, i8* %s) {
entry:
	%0 = ptrtoint %String* %self to i64
	%1 = bitcast %String* %self to i8*
	%2 = ptrtoint i8* %1 to i64
	%3 = icmp eq i64 %0, %2
	br i1 %3, label %raw_string, label %struct_string

raw_string:
	%4 = bitcast %String* %self to i8*
	%5 = call i32 @strlen(i8* %4)
	%6 = call i32 @strlen(i8* %s)
	%7 = add i32 %5, %6
	%8 = add i32 %7, 1
	%9 = zext i32 %8 to i64
	%10 = call i8* @malloc(i64 %9)
	%11 = call i8* @strcpy(i8* %10, i8* %4)
	%12 = call i8* @strcat(i8* %10, i8* %s)
	ret i8* %10

struct_string:
	%13 = bitcast %String* %self to i8**
	%14 = load i8*, i8** %13
	%15 = call i32 @strlen(i8* %14)
	%16 = call i32 @strlen(i8* %s)
	%17 = add i32 %15, %16
	%18 = add i32 %17, 1
	%19 = zext i32 %18 to i64
	%20 = call i8* @malloc(i64 %19)
	%21 = call i8* @strcpy(i8* %20, i8* %14)
	%22 = call i8* @strcat(i8* %20, i8* %s)
	ret i8* %20
}

define i8* @String.substr(%String* %self, i32 %i, i32 %l) {
entry:
	%0 = ptrtoint %String* %self to i64
	%1 = bitcast %String* %self to i8*
	%2 = ptrtoint i8* %1 to i64
	%3 = icmp eq i64 %0, %2
	br i1 %3, label %raw_string, label %struct_string

raw_string:
	%4 = bitcast %String* %self to i8*
	%5 = call i32 @strlen(i8* %4)
	br label %bounds_check_raw

struct_string:
	%6 = bitcast %String* %self to i8**
	%7 = load i8*, i8** %6
	%8 = call i32 @strlen(i8* %7)
	br label %bounds_check_struct

bounds_check_raw:
	%9 = icmp slt i32 %i, 0
	%10 = icmp sge i32 %i, %5
	%11 = or i1 %9, %10
	%12 = icmp slt i32 %l, 0
	%13 = or i1 %11, %12
	br i1 %13, label %error, label %alloc_raw

error:
	%14 = bitcast [39 x i8]* @.str.substr_error to i8*
	%15 = call i32 (i8*, ...) @printf(i8* %14)
	call void @exit(i32 1)
	unreachable

alloc_raw:
	%16 = add i32 %l, 1
	%17 = zext i32 %16 to i64
	%18 = call i8* @malloc(i64 %17)
	%19 = getelementptr i8, i8* %4, i32 %i
	%20 = call i8* @strncpy(i8* %18, i8* %19, i32 %l)
	%21 = getelementptr i8, i8* %18, i32 %l
	store i8 0, i8* %21
	ret i8* %18

bounds_check_struct:
	%22 = icmp slt i32 %i, 0
	%23 = icmp sge i32 %i, %8
	%24 = or i1 %22, %23
	%25 = icmp slt i32 %l, 0
	%26 = or i1 %24, %25
	br i1 %26, label %error, label %alloc_struct

alloc_struct:
	%27 = add i32 %l, 1
	%28 = zext i32 %27 to i64
	%29 = call i8* @malloc(i64 %28)
	%30 = getelementptr i8, i8* %7, i32 %i
	%31 = call i8* @strncpy(i8* %29, i8* %30, i32 %l)
	%32 = getelementptr i8, i8* %29, i32 %l
	store i8 0, i8* %32
	ret i8* %29
}

declare i8* @malloc(i64 %size)

declare void @free(i8* %ptr)

declare void @exit(i32 %status)

declare i32 @printf(i8* %format, ...)

declare i32 @scanf(i8* %format, ...)

declare i32 @strlen(i8* %str)

declare i8* @strcpy(i8* %dest, i8* %src)

declare i8* @strcat(i8* %dest, i8* %src)

declare i8* @strncpy(i8* %dest, i8* %src, i32 %n)

define %Node* @Node.init(%Node* %self, i32 %value, %Node* %nextNode) {
entry:
	%0 = getelementptr %Node, %Node* %self, i32 0, i32 1
	store i32 %value, i32* %0
	%1 = getelementptr %Node, %Node* %self, i32 0, i32 2
	store %Node* %nextNode, %Node** %1
	ret %Node* %self
}

define i32 @Node.getData(%Node* %self) {
entry:
	%0 = getelementptr %Node, %Node* %self, i32 0, i32 1
	%1 = load i32, i32* %0
	ret i32 %1
}

define %Node* @Node.getNext(%Node* %self) {
entry:
	%0 = getelementptr %Node, %Node* %self, i32 0, i32 2
	%1 = load %Node*, %Node** %0
	%2 = bitcast %Node* %1 to %Node*
	ret %Node* %2
}

define %Node* @Node.setNext(%Node* %self, %Node* %nextNode) {
entry:
	%0 = getelementptr %Node, %Node* %self, i32 0, i32 2
	store %Node* %nextNode, %Node** %0
	ret %Node* %self
}

define %LinkedList* @LinkedList.init(%LinkedList* %self) {
entry:
	%0 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	store i32 0, i32* %0
	ret %LinkedList* %self
}

define %LinkedList* @LinkedList.addFirst(%LinkedList* %self, i32 %value) {
entry:
	%0 = alloca %Node*
	%1 = getelementptr %Node, %Node* null, i32 1
	%2 = ptrtoint %Node* %1 to i64
	%3 = call i8* @malloc(i64 %2)
	%4 = bitcast i8* %3 to %Node*
	%5 = getelementptr %Node, %Node* %4, i32 0, i32 0
	%6 = bitcast [7 x i8*]* @vtable.Node to i8*
	store i8* %6, i8** %5
	%7 = getelementptr %Node, %Node* %4, i32 0, i32 1
	store i32 0, i32* %7
	%8 = getelementptr %Node, %Node* %4, i32 0, i32 2
	store %Node* null, %Node** %8
	store %Node* %4, %Node** %0
	%9 = load %Node*, %Node** %0
	%10 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 1
	%11 = load %Node*, %Node** %10
	%12 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 4
	%13 = load i8*, i8** %12
	%14 = bitcast i8* %13 to %Node* (%Node*, i32, %Node*)*
	%15 = bitcast %Node* %9 to %Node*
	%16 = bitcast %Node* %11 to %Node*
	%17 = call %Node* %14(%Node* %15, i32 %value, %Node* %16)
	%18 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 6
	%19 = load i8*, i8** %18
	%20 = bitcast i8* %19 to i1 (%LinkedList*)*
	%21 = bitcast %LinkedList* %self to %LinkedList*
	%22 = call i1 %20(%LinkedList* %21)
	br i1 %22, label %if.then.1, label %if.else.1

if.then.1:
	%23 = load %Node*, %Node** %0
	%24 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 2
	store %Node* %23, %Node** %24
	br label %if.end.1

if.else.1:
	br label %if.end.1

if.end.1:
	%25 = phi %Node* [ %23, %if.then.1 ], [ null, %if.else.1 ]
	%26 = load %Node*, %Node** %0
	%27 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 1
	store %Node* %26, %Node** %27
	%28 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	%29 = load i32, i32* %28
	%30 = add i32 %29, 1
	%31 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	store i32 %30, i32* %31
	ret %LinkedList* %self
}

define %LinkedList* @LinkedList.addLast(%LinkedList* %self, i32 %value) {
entry:
	%0 = alloca %Node*
	%1 = getelementptr %Node, %Node* null, i32 1
	%2 = ptrtoint %Node* %1 to i64
	%3 = call i8* @malloc(i64 %2)
	%4 = bitcast i8* %3 to %Node*
	%5 = getelementptr %Node, %Node* %4, i32 0, i32 0
	%6 = bitcast [7 x i8*]* @vtable.Node to i8*
	store i8* %6, i8** %5
	%7 = getelementptr %Node, %Node* %4, i32 0, i32 1
	store i32 0, i32* %7
	%8 = getelementptr %Node, %Node* %4, i32 0, i32 2
	store %Node* null, %Node** %8
	store %Node* %4, %Node** %0
	%9 = load %Node*, %Node** %0
	%10 = load %Node*, %Node** %0
	%11 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 4
	%12 = load i8*, i8** %11
	%13 = bitcast i8* %12 to %Node* (%Node*, i32, %Node*)*
	%14 = bitcast %Node* %9 to %Node*
	%15 = bitcast %Node* %10 to %Node*
	%16 = call %Node* %13(%Node* %14, i32 %value, %Node* %15)
	%17 = alloca %Node*
	store %Node* null, %Node** %17
	%18 = load %Node*, %Node** %0
	%19 = load %Node*, %Node** %17
	%20 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 5
	%21 = load i8*, i8** %20
	%22 = bitcast i8* %21 to %Node* (%Node*, %Node*)*
	%23 = bitcast %Node* %18 to %Node*
	%24 = bitcast %Node* %19 to %Node*
	%25 = call %Node* %22(%Node* %23, %Node* %24)
	%26 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 6
	%27 = load i8*, i8** %26
	%28 = bitcast i8* %27 to i1 (%LinkedList*)*
	%29 = bitcast %LinkedList* %self to %LinkedList*
	%30 = call i1 %28(%LinkedList* %29)
	br i1 %30, label %if.then.2, label %if.else.2

if.then.2:
	%31 = load %Node*, %Node** %0
	%32 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 1
	store %Node* %31, %Node** %32
	br label %if.end.2

if.else.2:
	%33 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 2
	%34 = load %Node*, %Node** %33
	%35 = load %Node*, %Node** %0
	%36 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 5
	%37 = load i8*, i8** %36
	%38 = bitcast i8* %37 to %Node* (%Node*, %Node*)*
	%39 = bitcast %Node* %34 to %Node*
	%40 = bitcast %Node* %35 to %Node*
	%41 = call %Node* %38(%Node* %39, %Node* %40)
	br label %if.end.2

if.end.2:
	%42 = phi %Node* [ %31, %if.then.2 ], [ %41, %if.else.2 ]
	%43 = load %Node*, %Node** %0
	%44 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 2
	store %Node* %43, %Node** %44
	%45 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	%46 = load i32, i32* %45
	%47 = add i32 %46, 1
	%48 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	store i32 %47, i32* %48
	ret %LinkedList* %self
}

define i32 @LinkedList.removeFirst(%LinkedList* %self) {
entry:
	%0 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 6
	%1 = load i8*, i8** %0
	%2 = bitcast i8* %1 to i1 (%LinkedList*)*
	%3 = bitcast %LinkedList* %self to %LinkedList*
	%4 = call i1 %2(%LinkedList* %3)
	br i1 %4, label %if.then.3, label %if.else.3

if.then.3:
	%5 = getelementptr %IO, %IO* null, i32 1
	%6 = ptrtoint %IO* %5 to i64
	%7 = call i8* @malloc(i64 %6)
	%8 = bitcast i8* %7 to %IO*
	%9 = getelementptr %IO, %IO* %8, i32 0, i32 0
	%10 = bitcast [7 x i8*]* @vtable.IO to i8*
	store i8* %10, i8** %9
	%11 = bitcast i8* getelementptr ([41 x i8], [41 x i8]* @.str9, i32 0, i32 0) to i8*
	%12 = call %IO* @IO.out_string(%IO* %8, i8* %11)
	%13 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 0
	%14 = load i8*, i8** %13
	%15 = bitcast i8* %14 to %Object* (%Object*)*
	%16 = bitcast %LinkedList* %self to %Object*
	%17 = call %Object* %15(%Object* %16)
	br label %if.end.3

if.else.3:
	%18 = alloca i32
	%19 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 1
	%20 = load %Node*, %Node** %19
	%21 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 2
	%22 = load i8*, i8** %21
	%23 = bitcast i8* %22 to i32 (%Node*)*
	%24 = bitcast %Node* %20 to %Node*
	%25 = call i32 %23(%Node* %24)
	store i32 %25, i32* %18
	%26 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 1
	%27 = load %Node*, %Node** %26
	%28 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 3
	%29 = load i8*, i8** %28
	%30 = bitcast i8* %29 to %Node* (%Node*)*
	%31 = bitcast %Node* %27 to %Node*
	%32 = call %Node* %30(%Node* %31)
	%33 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 1
	store %Node* %32, %Node** %33
	%34 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	%35 = load i32, i32* %34
	%36 = sub i32 %35, 1
	%37 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	store i32 %36, i32* %37
	%38 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 6
	%39 = load i8*, i8** %38
	%40 = bitcast i8* %39 to i1 (%LinkedList*)*
	%41 = bitcast %LinkedList* %self to %LinkedList*
	%42 = call i1 %40(%LinkedList* %41)
	br i1 %42, label %if.then.4, label %if.else.4

if.end.3:
	%43 = phi i32 [ 0, %if.then.3 ], [ %48, %if.end.4 ]
	ret i32 %43

if.then.4:
	%44 = alloca %Node*
	store %Node* null, %Node** %44
	%45 = load %Node*, %Node** %44
	%46 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 2
	store %Node* %45, %Node** %46
	br label %if.end.4

if.else.4:
	br label %if.end.4

if.end.4:
	%47 = phi %Node* [ %45, %if.then.4 ], [ null, %if.else.4 ]
	%48 = load i32, i32* %18
	br label %if.end.3
}

define i32 @LinkedList.get(%LinkedList* %self, i32 %index) {
entry:
	%0 = icmp slt i32 %index, 0
	br i1 %0, label %if.then.5, label %if.else.5

if.then.5:
	%1 = getelementptr %IO, %IO* null, i32 1
	%2 = ptrtoint %IO* %1 to i64
	%3 = call i8* @malloc(i64 %2)
	%4 = bitcast i8* %3 to %IO*
	%5 = getelementptr %IO, %IO* %4, i32 0, i32 0
	%6 = bitcast [7 x i8*]* @vtable.IO to i8*
	store i8* %6, i8** %5
	%7 = bitcast i8* getelementptr ([45 x i8], [45 x i8]* @.str10, i32 0, i32 0) to i8*
	%8 = call %IO* @IO.out_string(%IO* %4, i8* %7)
	%9 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 0
	%10 = load i8*, i8** %9
	%11 = bitcast i8* %10 to %Object* (%Object*)*
	%12 = bitcast %LinkedList* %self to %Object*
	%13 = call %Object* %11(%Object* %12)
	br label %if.end.5

if.else.5:
	%14 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	%15 = load i32, i32* %14
	%16 = icmp sle i32 %15, %index
	br i1 %16, label %if.then.6, label %if.else.6

if.end.5:
	%17 = phi i32 [ 0, %if.then.5 ], [ %35, %if.end.6 ]
	ret i32 %17

if.then.6:
	%18 = getelementptr %IO, %IO* null, i32 1
	%19 = ptrtoint %IO* %18 to i64
	%20 = call i8* @malloc(i64 %19)
	%21 = bitcast i8* %20 to %IO*
	%22 = getelementptr %IO, %IO* %21, i32 0, i32 0
	%23 = bitcast [7 x i8*]* @vtable.IO to i8*
	store i8* %23, i8** %22
	%24 = bitcast i8* getelementptr ([46 x i8], [46 x i8]* @.str11, i32 0, i32 0) to i8*
	%25 = call %IO* @IO.out_string(%IO* %21, i8* %24)
	%26 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 0
	%27 = load i8*, i8** %26
	%28 = bitcast i8* %27 to %Object* (%Object*)*
	%29 = bitcast %LinkedList* %self to %Object*
	%30 = call %Object* %28(%Object* %29)
	br label %if.end.6

if.else.6:
	%31 = alloca %Node*
	%32 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 1
	%33 = load %Node*, %Node** %32
	store %Node* %33, %Node** %31
	%34 = alloca i32
	store i32 0, i32* %34
	br label %while.cond.1

if.end.6:
	%35 = phi i32 [ 0, %if.then.6 ], [ %51, %while.exit.1 ]
	br label %if.end.5

while.cond.1:
	%36 = load i32, i32* %34
	%37 = icmp slt i32 %36, %index
	br i1 %37, label %while.body.1, label %while.exit.1

while.body.1:
	%38 = load %Node*, %Node** %31
	%39 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 3
	%40 = load i8*, i8** %39
	%41 = bitcast i8* %40 to %Node* (%Node*)*
	%42 = bitcast %Node* %38 to %Node*
	%43 = call %Node* %41(%Node* %42)
	store %Node* %43, %Node** %31
	%44 = load i32, i32* %34
	%45 = add i32 %44, 1
	store i32 %45, i32* %34
	br label %while.cond.1

while.exit.1:
	%46 = load %Node*, %Node** %31
	%47 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 2
	%48 = load i8*, i8** %47
	%49 = bitcast i8* %48 to i32 (%Node*)*
	%50 = bitcast %Node* %46 to %Node*
	%51 = call i32 %49(%Node* %50)
	br label %if.end.6
}

define i32 @LinkedList.size(%LinkedList* %self) {
entry:
	%0 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	%1 = load i32, i32* %0
	ret i32 %1
}

define i1 @LinkedList.isEmpty(%LinkedList* %self) {
entry:
	%0 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 1
	%1 = load %Node*, %Node** %0
	%2 = icmp eq %Node* %1, null
	ret i1 %2
}

define %Object* @Main.main(%Main* %self) {
entry:
	%0 = alloca %LinkedList*
	%1 = getelementptr %LinkedList, %LinkedList* null, i32 1
	%2 = ptrtoint %LinkedList* %1 to i64
	%3 = call i8* @malloc(i64 %2)
	%4 = bitcast i8* %3 to %LinkedList*
	%5 = getelementptr %LinkedList, %LinkedList* %4, i32 0, i32 0
	%6 = bitcast [10 x i8*]* @vtable.LinkedList to i8*
	store i8* %6, i8** %5
	%7 = getelementptr %LinkedList, %LinkedList* %4, i32 0, i32 1
	store %Node* null, %Node** %7
	%8 = getelementptr %LinkedList, %LinkedList* %4, i32 0, i32 2
	store %Node* null, %Node** %8
	%9 = getelementptr %LinkedList, %LinkedList* %4, i32 0, i32 3
	store i32 0, i32* %9
	store %LinkedList* %4, %LinkedList** %0
	%10 = load %LinkedList*, %LinkedList** %0
	%11 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 5
	%12 = load i8*, i8** %11
	%13 = bitcast i8* %12 to %LinkedList* (%LinkedList*)*
	%14 = bitcast %LinkedList* %10 to %LinkedList*
	%15 = call %LinkedList* %13(%LinkedList* %14)
	%16 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 6
	%17 = load i8*, i8** %16
	%18 = bitcast i8* %17 to %IO* (%IO*, i8*)*
	%19 = bitcast %Main* %self to %IO*
	%20 = bitcast i8* getelementptr ([26 x i8], [26 x i8]* @.str12, i32 0, i32 0) to i8*
	%21 = call %IO* %18(%IO* %19, i8* %20)
	%22 = load %LinkedList*, %LinkedList** %0
	%23 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 1
	%24 = load i8*, i8** %23
	%25 = bitcast i8* %24 to %LinkedList* (%LinkedList*, i32)*
	%26 = bitcast %LinkedList* %22 to %LinkedList*
	%27 = call %LinkedList* %25(%LinkedList* %26, i32 100)
	%28 = load %LinkedList*, %LinkedList** %0
	%29 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 2
	%30 = load i8*, i8** %29
	%31 = bitcast i8* %30 to %LinkedList* (%LinkedList*, i32)*
	%32 = bitcast %LinkedList* %28 to %LinkedList*
	%33 = call %LinkedList* %31(%LinkedList* %32, i32 200)
	%34 = load %LinkedList*, %LinkedList** %0
	%35 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 1
	%36 = load i8*, i8** %35
	%37 = bitcast i8* %36 to %LinkedList* (%LinkedList*, i32)*
	%38 = bitcast %LinkedList* %34 to %LinkedList*
	%39 = call %LinkedList* %37(%LinkedList* %38, i32 50)
	%40 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 6
	%41 = load i8*, i8** %40
	%42 = bitcast i8* %41 to %IO* (%IO*, i8*)*
	%43 = bitcast %Main* %self to %IO*
	%44 = bitcast i8* getelementptr ([12 x i8], [12 x i8]* @.str13, i32 0, i32 0) to i8*
	%45 = call %IO* %42(%IO* %43, i8* %44)
	%46 = load %LinkedList*, %LinkedList** %0
	%47 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 8
	%48 = load i8*, i8** %47
	%49 = bitcast i8* %48 to i32 (%LinkedList*)*
	%50 = bitcast %LinkedList* %46 to %LinkedList*
	%51 = call i32 %49(%LinkedList* %50)
	%52 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 5
	%53 = load i8*, i8** %52
	%54 = bitcast i8* %53 to %IO* (%IO*, i32)*
	%55 = bitcast %Main* %self to %IO*
	%56 = call %IO* %54(%IO* %55, i32 %51)
	%57 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 6
	%58 = load i8*, i8** %57
	%59 = bitcast i8* %58 to %IO* (%IO*, i8*)*
	%60 = bitcast %Main* %self to %IO*
	%61 = bitcast i8* getelementptr ([2 x i8], [2 x i8]* @.str14, i32 0, i32 0) to i8*
	%62 = call %IO* %59(%IO* %60, i8* %61)
	%63 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 6
	%64 = load i8*, i8** %63
	%65 = bitcast i8* %64 to %IO* (%IO*, i8*)*
	%66 = bitcast %Main* %self to %IO*
	%67 = bitcast i8* getelementptr ([16 x i8], [16 x i8]* @.str15, i32 0, i32 0) to i8*
	%68 = call %IO* %65(%IO* %66, i8* %67)
	%69 = load %LinkedList*, %LinkedList** %0
	%70 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 4
	%71 = load i8*, i8** %70
	%72 = bitcast i8* %71 to i32 (%LinkedList*, i32)*
	%73 = bitcast %LinkedList* %69 to %LinkedList*
	%74 = call i32 %72(%LinkedList* %73, i32 0)
	%75 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 5
	%76 = load i8*, i8** %75
	%77 = bitcast i8* %76 to %IO* (%IO*, i32)*
	%78 = bitcast %Main* %self to %IO*
	%79 = call %IO* %77(%IO* %78, i32 %74)
	%80 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 6
	%81 = load i8*, i8** %80
	%82 = bitcast i8* %81 to %IO* (%IO*, i8*)*
	%83 = bitcast %Main* %self to %IO*
	%84 = bitcast i8* getelementptr ([2 x i8], [2 x i8]* @.str16, i32 0, i32 0) to i8*
	%85 = call %IO* %82(%IO* %83, i8* %84)
	%86 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 6
	%87 = load i8*, i8** %86
	%88 = bitcast i8* %87 to %IO* (%IO*, i8*)*
	%89 = bitcast %Main* %self to %IO*
	%90 = bitcast i8* getelementptr ([17 x i8], [17 x i8]* @.str17, i32 0, i32 0) to i8*
	%91 = call %IO* %88(%IO* %89, i8* %90)
	%92 = load %LinkedList*, %LinkedList** %0
	%93 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 4
	%94 = load i8*, i8** %93
	%95 = bitcast i8* %94 to i32 (%LinkedList*, i32)*
	%96 = bitcast %LinkedList* %92 to %LinkedList*
	%97 = call i32 %95(%LinkedList* %96, i32 1)
	%98 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 5
	%99 = load i8*, i8** %98
	%100 = bitcast i8* %99 to %IO* (%IO*, i32)*
	%101 = bitcast %Main* %self to %IO*
	%102 = call %IO* %100(%IO* %101, i32 %97)
	%103 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 6
	%104 = load i8*, i8** %103
	%105 = bitcast i8* %104 to %IO* (%IO*, i8*)*
	%106 = bitcast %Main* %self to %IO*
	%107 = bitcast i8* getelementptr ([2 x i8], [2 x i8]* @.str18, i32 0, i32 0) to i8*
	%108 = call %IO* %105(%IO* %106, i8* %107)
	%109 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 6
	%110 = load i8*, i8** %109
	%111 = bitcast i8* %110 to %IO* (%IO*, i8*)*
	%112 = bitcast %Main* %self to %IO*
	%113 = bitcast i8* getelementptr ([16 x i8], [16 x i8]* @.str19, i32 0, i32 0) to i8*
	%114 = call %IO* %111(%IO* %112, i8* %113)
	%115 = load %LinkedList*, %LinkedList** %0
	%116 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 4
	%117 = load i8*, i8** %116
	%118 = bitcast i8* %117 to i32 (%LinkedList*, i32)*
	%119 = bitcast %LinkedList* %115 to %LinkedList*
	%120 = call i32 %118(%LinkedList* %119, i32 2)
	%121 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 5
	%122 = load i8*, i8** %121
	%123 = bitcast i8* %122 to %IO* (%IO*, i32)*
	%124 = bitcast %Main* %self to %IO*
	%125 = call %IO* %123(%IO* %124, i32 %120)
	%126 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 6
	%127 = load i8*, i8** %126
	%128 = bitcast i8* %127 to %IO* (%IO*, i8*)*
	%129 = bitcast %Main* %self to %IO*
	%130 = bitcast i8* getelementptr ([2 x i8], [2 x i8]* @.str20, i32 0, i32 0) to i8*
	%131 = call %IO* %128(%IO* %129, i8* %130)
	%132 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 6
	%133 = load i8*, i8** %132
	%134 = bitcast i8* %133 to %IO* (%IO*, i8*)*
	%135 = bitcast %Main* %self to %IO*
	%136 = bitcast i8* getelementptr ([24 x i8], [24 x i8]* @.str21, i32 0, i32 0) to i8*
	%137 = call %IO* %134(%IO* %135, i8* %136)
	%138 = load %LinkedList*, %LinkedList** %0
	%139 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 7
	%140 = load i8*, i8** %139
	%141 = bitcast i8* %140 to i32 (%LinkedList*)*
	%142 = bitcast %LinkedList* %138 to %LinkedList*
	%143 = call i32 %141(%LinkedList* %142)
	%144 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 5
	%145 = load i8*, i8** %144
	%146 = bitcast i8* %145 to %IO* (%IO*, i32)*
	%147 = bitcast %Main* %self to %IO*
	%148 = call %IO* %146(%IO* %147, i32 %143)
	%149 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 6
	%150 = load i8*, i8** %149
	%151 = bitcast i8* %150 to %IO* (%IO*, i8*)*
	%152 = bitcast %Main* %self to %IO*
	%153 = bitcast i8* getelementptr ([2 x i8], [2 x i8]* @.str22, i32 0, i32 0) to i8*
	%154 = call %IO* %151(%IO* %152, i8* %153)
	%155 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 6
	%156 = load i8*, i8** %155
	%157 = bitcast i8* %156 to %IO* (%IO*, i8*)*
	%158 = bitcast %Main* %self to %IO*
	%159 = bitcast i8* getelementptr ([16 x i8], [16 x i8]* @.str23, i32 0, i32 0) to i8*
	%160 = call %IO* %157(%IO* %158, i8* %159)
	%161 = load %LinkedList*, %LinkedList** %0
	%162 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 8
	%163 = load i8*, i8** %162
	%164 = bitcast i8* %163 to i32 (%LinkedList*)*
	%165 = bitcast %LinkedList* %161 to %LinkedList*
	%166 = call i32 %164(%LinkedList* %165)
	%167 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 5
	%168 = load i8*, i8** %167
	%169 = bitcast i8* %168 to %IO* (%IO*, i32)*
	%170 = bitcast %Main* %self to %IO*
	%171 = call %IO* %169(%IO* %170, i32 %166)
	%172 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 6
	%173 = load i8*, i8** %172
	%174 = bitcast i8* %173 to %IO* (%IO*, i8*)*
	%175 = bitcast %Main* %self to %IO*
	%176 = bitcast i8* getelementptr ([2 x i8], [2 x i8]* @.str24, i32 0, i32 0) to i8*
	%177 = call %IO* %174(%IO* %175, i8* %176)
	%178 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 6
	%179 = load i8*, i8** %178
	%180 = bitcast i8* %179 to %IO* (%IO*, i8*)*
	%181 = bitcast %Main* %self to %IO*
	%182 = bitcast i8* getelementptr ([20 x i8], [20 x i8]* @.str25, i32 0, i32 0) to i8*
	%183 = call %IO* %180(%IO* %181, i8* %182)
	%184 = load %LinkedList*, %LinkedList** %0
	%185 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 4
	%186 = load i8*, i8** %185
	%187 = bitcast i8* %186 to i32 (%LinkedList*, i32)*
	%188 = bitcast %LinkedList* %184 to %LinkedList*
	%189 = call i32 %187(%LinkedList* %188, i32 0)
	%190 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 5
	%191 = load i8*, i8** %190
	%192 = bitcast i8* %191 to %IO* (%IO*, i32)*
	%193 = bitcast %Main* %self to %IO*
	%194 = call %IO* %192(%IO* %193, i32 %189)
	%195 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 6
	%196 = load i8*, i8** %195
	%197 = bitcast i8* %196 to %IO* (%IO*, i8*)*
	%198 = bitcast %Main* %self to %IO*
	%199 = bitcast i8* getelementptr ([2 x i8], [2 x i8]* @.str26, i32 0, i32 0) to i8*
	%200 = call %IO* %197(%IO* %198, i8* %199)
	%201 = load %LinkedList*, %LinkedList** %0
	%202 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 6
	%203 = load i8*, i8** %202
	%204 = bitcast i8* %203 to i1 (%LinkedList*)*
	%205 = bitcast %LinkedList* %201 to %LinkedList*
	%206 = call i1 %204(%LinkedList* %205)
	br i1 %206, label %if.then.7, label %if.else.7

if.then.7:
	%207 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 6
	%208 = load i8*, i8** %207
	%209 = bitcast i8* %208 to %IO* (%IO*, i8*)*
	%210 = bitcast %Main* %self to %IO*
	%211 = bitcast i8* getelementptr ([15 x i8], [15 x i8]* @.str27, i32 0, i32 0) to i8*
	%212 = call %IO* %209(%IO* %210, i8* %211)
	br label %if.end.7

if.else.7:
	%213 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 6
	%214 = load i8*, i8** %213
	%215 = bitcast i8* %214 to %IO* (%IO*, i8*)*
	%216 = bitcast %Main* %self to %IO*
	%217 = bitcast i8* getelementptr ([19 x i8], [19 x i8]* @.str28, i32 0, i32 0) to i8*
	%218 = call %IO* %215(%IO* %216, i8* %217)
	br label %if.end.7

if.end.7:
	%219 = phi %IO* [ %212, %if.then.7 ], [ %218, %if.else.7 ]
	%220 = load %LinkedList*, %LinkedList** %0
	%221 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 7
	%222 = load i8*, i8** %221
	%223 = bitcast i8* %222 to i32 (%LinkedList*)*
	%224 = bitcast %LinkedList* %220 to %LinkedList*
	%225 = call i32 %223(%LinkedList* %224)
	%226 = load %LinkedList*, %LinkedList** %0
	%227 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 7
	%228 = load i8*, i8** %227
	%229 = bitcast i8* %228 to i32 (%LinkedList*)*
	%230 = bitcast %LinkedList* %226 to %LinkedList*
	%231 = call i32 %229(%LinkedList* %230)
	%232 = load %LinkedList*, %LinkedList** %0
	%233 = getelementptr [10 x i8*], [10 x i8*]* @vtable.LinkedList, i32 0, i32 6
	%234 = load i8*, i8** %233
	%235 = bitcast i8* %234 to i1 (%LinkedList*)*
	%236 = bitcast %LinkedList* %232 to %LinkedList*
	%237 = call i1 %235(%LinkedList* %236)
	br i1 %237, label %if.then.8, label %if.else.8

if.then.8:
	%238 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 6
	%239 = load i8*, i8** %238
	%240 = bitcast i8* %239 to %IO* (%IO*, i8*)*
	%241 = bitcast %Main* %self to %IO*
	%242 = bitcast i8* getelementptr ([15 x i8], [15 x i8]* @.str29, i32 0, i32 0) to i8*
	%243 = call %IO* %240(%IO* %241, i8* %242)
	br label %if.end.8

if.else.8:
	%244 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 6
	%245 = load i8*, i8** %244
	%246 = bitcast i8* %245 to %IO* (%IO*, i8*)*
	%247 = bitcast %Main* %self to %IO*
	%248 = bitcast i8* getelementptr ([19 x i8], [19 x i8]* @.str30, i32 0, i32 0) to i8*
	%249 = call %IO* %246(%IO* %247, i8* %248)
	br label %if.end.8

if.end.8:
	%250 = phi %IO* [ %243, %if.then.8 ], [ %249, %if.else.8 ]
	%251 = bitcast %IO* %250 to %Object*
	ret %Object* %251
}

define i32 @main() {
entry:
	%0 = getelementptr %Main, %Main* null, i32 1
	%1 = ptrtoint %Main* %0 to i64
	%2 = call i8* @malloc(i64 %1)
	%3 = bitcast i8* %2 to %Main*
	%4 = getelementptr %Main, %Main* %3, i32 0, i32 0
	%5 = bitcast [8 x i8*]* @vtable.Main to i8*
	store i8* %5, i8** %4
	%6 = alloca %Main*
	store %Main* %3, %Main** %6
	%7 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 4
	%8 = load i8*, i8** %7
	%9 = bitcast i8* %8 to i8* (%Main*)*
	%10 = call i8* %9(%Main* %3)
	ret i32 0
}
