%Object = type { i8* }
%IO = type { i8* }
%Int = type { i8* }
%String = type { i8* }
%Bool = type { i8* }
%Node = type { i8*, %Object*, %Node* }
%LinkedList = type { i8*, %Node*, %Node*, i32, %IO* }
%Main = type { i8* }

@.str.empty = constant [1 x i8] c"\00"
@vtable.Object = global [3 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.IO = global [7 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%IO*)* @IO.in_int to i8*), i8* bitcast (i8* (%IO*)* @IO.in_string to i8*), i8* bitcast (%IO* (%IO*, i32)* @IO.out_int to i8*), i8* bitcast (%IO* (%IO*, i8*)* @IO.out_string to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Int = global [3 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.String = global [6 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (i8* (%String*, i8*)* @String.concat to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%String*)* @String.length to i8*), i8* bitcast (i8* (%String*, i32, i32)* @String.substr to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Bool = global [3 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Node = global [7 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (%Object* (%Node*)* @Node.getData to i8*), i8* bitcast (%Node* (%Node*)* @Node.getNext to i8*), i8* bitcast (%Node* (%Node*, %Object*, %Node*)* @Node.init to i8*), i8* bitcast (%Node* (%Node*, %Node*)* @Node.setNext to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.LinkedList = global [11 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%LinkedList* (%LinkedList*, %Object*)* @LinkedList.addFirst to i8*), i8* bitcast (%LinkedList* (%LinkedList*, %Object*)* @LinkedList.addLast to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (%Object* (%LinkedList*, i32)* @LinkedList.get to i8*), i8* bitcast (%LinkedList* (%LinkedList*)* @LinkedList.init to i8*), i8* bitcast (i1 (%LinkedList*)* @LinkedList.isEmpty to i8*), i8* bitcast (%Object* (%LinkedList*)* @LinkedList.removeFirst to i8*), i8* bitcast (%Object* (%LinkedList*)* @LinkedList.removeLast to i8*), i8* bitcast (i32 (%LinkedList*)* @LinkedList.size to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Main = global [4 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (%Object* (%Main*)* @Main.main to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@.str9 = internal constant [41 x i8] c"Error: Cannot remove from an empty list\0A\00"
@.str10 = internal constant [41 x i8] c"Error: Cannot remove from an empty list\0A\00"
@.str11 = internal constant [45 x i8] c"Error: Index out of bounds (negative index)\0A\00"
@.str12 = internal constant [46 x i8] c"Error: Index out of bounds (index too large)\0A\00"
@.str13 = internal constant [12 x i8] c"List size: \00"
@.str14 = internal constant [2 x i8] c"\0A\00"
@.str15 = internal constant [16 x i8] c"First element: \00"
@.str16 = internal constant [2 x i8] c"\0A\00"
@.str17 = internal constant [10 x i8] c"Removed: \00"
@.str18 = internal constant [2 x i8] c"\0A\00"
@.str19 = internal constant [11 x i8] c"New size: \00"
@.str20 = internal constant [2 x i8] c"\0A\00"
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

define %Node* @Node.init(%Node* %self, %Object* %value, %Node* %nextNode) {
entry:
	%0 = getelementptr %Node, %Node* %self, i32 0, i32 1
	store %Object* %value, %Object** %0
	%1 = getelementptr %Node, %Node* %self, i32 0, i32 2
	store %Node* %nextNode, %Node** %1
	ret %Node* %self
}

define %Object* @Node.getData(%Node* %self) {
entry:
	%0 = getelementptr %Node, %Node* %self, i32 0, i32 1
	%1 = load %Object*, %Object** %0
	%2 = bitcast %Object* %1 to %Object*
	ret %Object* %2
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

define %LinkedList* @LinkedList.addFirst(%LinkedList* %self, %Object* %value) {
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
	store %Object* null, %Object** %7
	%8 = getelementptr %Node, %Node* %4, i32 0, i32 2
	store %Node* null, %Node** %8
	store %Node* %4, %Node** %0
	%9 = load %Node*, %Node** %0
	%10 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 1
	%11 = load %Node*, %Node** %10
	%12 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 4
	%13 = load i8*, i8** %12
	%14 = bitcast i8* %13 to %Node* (%Node*, %Object*, %Node*)*
	%15 = bitcast %Node* %9 to %Node*
	%16 = bitcast %Object* %value to %Object*
	%17 = bitcast %Node* %11 to %Node*
	%18 = call %Node* %14(%Node* %15, %Object* %16, %Node* %17)
	%19 = getelementptr [11 x i8*], [11 x i8*]* @vtable.LinkedList, i32 0, i32 6
	%20 = load i8*, i8** %19
	%21 = bitcast i8* %20 to i1 (%LinkedList*)*
	%22 = bitcast %LinkedList* %self to %LinkedList*
	%23 = call i1 %21(%LinkedList* %22)
	br i1 %23, label %if.then.1, label %if.else.1

if.then.1:
	%24 = load %Node*, %Node** %0
	%25 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 2
	store %Node* %24, %Node** %25
	br label %if.end.1

if.else.1:
	br label %if.end.1

if.end.1:
	%26 = phi %Node* [ %24, %if.then.1 ], [ null, %if.else.1 ]
	%27 = load %Node*, %Node** %0
	%28 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 1
	store %Node* %27, %Node** %28
	%29 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	%30 = load i32, i32* %29
	%31 = add i32 %30, 1
	%32 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	store i32 %31, i32* %32
	ret %LinkedList* %self
}

define %LinkedList* @LinkedList.addLast(%LinkedList* %self, %Object* %value) {
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
	store %Object* null, %Object** %7
	%8 = getelementptr %Node, %Node* %4, i32 0, i32 2
	store %Node* null, %Node** %8
	store %Node* %4, %Node** %0
	%9 = load %Node*, %Node** %0
	%10 = load %Node*, %Node** %0
	%11 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 4
	%12 = load i8*, i8** %11
	%13 = bitcast i8* %12 to %Node* (%Node*, %Object*, %Node*)*
	%14 = bitcast %Node* %9 to %Node*
	%15 = bitcast %Object* %value to %Object*
	%16 = bitcast %Node* %10 to %Node*
	%17 = call %Node* %13(%Node* %14, %Object* %15, %Node* %16)
	%18 = alloca %Node*
	store %Node* null, %Node** %18
	%19 = load %Node*, %Node** %0
	%20 = load %Node*, %Node** %18
	%21 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 5
	%22 = load i8*, i8** %21
	%23 = bitcast i8* %22 to %Node* (%Node*, %Node*)*
	%24 = bitcast %Node* %19 to %Node*
	%25 = bitcast %Node* %20 to %Node*
	%26 = call %Node* %23(%Node* %24, %Node* %25)
	%27 = getelementptr [11 x i8*], [11 x i8*]* @vtable.LinkedList, i32 0, i32 6
	%28 = load i8*, i8** %27
	%29 = bitcast i8* %28 to i1 (%LinkedList*)*
	%30 = bitcast %LinkedList* %self to %LinkedList*
	%31 = call i1 %29(%LinkedList* %30)
	br i1 %31, label %if.then.2, label %if.else.2

if.then.2:
	%32 = load %Node*, %Node** %0
	%33 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 1
	store %Node* %32, %Node** %33
	br label %if.end.2

if.else.2:
	%34 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 2
	%35 = load %Node*, %Node** %34
	%36 = load %Node*, %Node** %0
	%37 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 5
	%38 = load i8*, i8** %37
	%39 = bitcast i8* %38 to %Node* (%Node*, %Node*)*
	%40 = bitcast %Node* %35 to %Node*
	%41 = bitcast %Node* %36 to %Node*
	%42 = call %Node* %39(%Node* %40, %Node* %41)
	br label %if.end.2

if.end.2:
	%43 = phi %Node* [ %32, %if.then.2 ], [ %42, %if.else.2 ]
	%44 = load %Node*, %Node** %0
	%45 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 2
	store %Node* %44, %Node** %45
	%46 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	%47 = load i32, i32* %46
	%48 = add i32 %47, 1
	%49 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	store i32 %48, i32* %49
	ret %LinkedList* %self
}

define %Object* @LinkedList.removeFirst(%LinkedList* %self) {
entry:
	%0 = getelementptr [11 x i8*], [11 x i8*]* @vtable.LinkedList, i32 0, i32 6
	%1 = load i8*, i8** %0
	%2 = bitcast i8* %1 to i1 (%LinkedList*)*
	%3 = bitcast %LinkedList* %self to %LinkedList*
	%4 = call i1 %2(%LinkedList* %3)
	br i1 %4, label %if.then.3, label %if.else.3

if.then.3:
	%5 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 4
	%6 = load %IO*, %IO** %5
	%7 = bitcast i8* getelementptr ([41 x i8], [41 x i8]* @.str9, i32 0, i32 0) to i8*
	%8 = call %IO* @IO.out_string(%IO* %6, i8* %7)
	%9 = getelementptr [11 x i8*], [11 x i8*]* @vtable.LinkedList, i32 0, i32 0
	%10 = load i8*, i8** %9
	%11 = bitcast i8* %10 to %Object* (%Object*)*
	%12 = bitcast %LinkedList* %self to %Object*
	%13 = call %Object* %11(%Object* %12)
	%14 = getelementptr %Object, %Object* null, i32 1
	%15 = ptrtoint %Object* %14 to i64
	%16 = call i8* @malloc(i64 %15)
	%17 = bitcast i8* %16 to %Object*
	%18 = getelementptr %Object, %Object* %17, i32 0, i32 0
	%19 = bitcast [3 x i8*]* @vtable.Object to i8*
	store i8* %19, i8** %18
	br label %if.end.3

if.else.3:
	%20 = alloca %Object*
	%21 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 1
	%22 = load %Node*, %Node** %21
	%23 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 2
	%24 = load i8*, i8** %23
	%25 = bitcast i8* %24 to %Object* (%Node*)*
	%26 = bitcast %Node* %22 to %Node*
	%27 = call %Object* %25(%Node* %26)
	store %Object* %27, %Object** %20
	%28 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 1
	%29 = load %Node*, %Node** %28
	%30 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 3
	%31 = load i8*, i8** %30
	%32 = bitcast i8* %31 to %Node* (%Node*)*
	%33 = bitcast %Node* %29 to %Node*
	%34 = call %Node* %32(%Node* %33)
	%35 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 1
	store %Node* %34, %Node** %35
	%36 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	%37 = load i32, i32* %36
	%38 = sub i32 %37, 1
	%39 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	store i32 %38, i32* %39
	%40 = getelementptr [11 x i8*], [11 x i8*]* @vtable.LinkedList, i32 0, i32 6
	%41 = load i8*, i8** %40
	%42 = bitcast i8* %41 to i1 (%LinkedList*)*
	%43 = bitcast %LinkedList* %self to %LinkedList*
	%44 = call i1 %42(%LinkedList* %43)
	br i1 %44, label %if.then.4, label %if.else.4

if.end.3:
	%45 = phi %Object* [ %17, %if.then.3 ], [ %51, %if.end.4 ]
	%46 = bitcast %Object* %45 to %Object*
	ret %Object* %46

if.then.4:
	%47 = alloca %Node*
	store %Node* null, %Node** %47
	%48 = load %Node*, %Node** %47
	%49 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 2
	store %Node* %48, %Node** %49
	br label %if.end.4

if.else.4:
	br label %if.end.4

if.end.4:
	%50 = phi %Node* [ %48, %if.then.4 ], [ null, %if.else.4 ]
	%51 = load %Object*, %Object** %20
	br label %if.end.3
}

define %Object* @LinkedList.removeLast(%LinkedList* %self) {
entry:
	%0 = getelementptr [11 x i8*], [11 x i8*]* @vtable.LinkedList, i32 0, i32 6
	%1 = load i8*, i8** %0
	%2 = bitcast i8* %1 to i1 (%LinkedList*)*
	%3 = bitcast %LinkedList* %self to %LinkedList*
	%4 = call i1 %2(%LinkedList* %3)
	br i1 %4, label %if.then.5, label %if.else.5

if.then.5:
	%5 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 4
	%6 = load %IO*, %IO** %5
	%7 = bitcast i8* getelementptr ([41 x i8], [41 x i8]* @.str10, i32 0, i32 0) to i8*
	%8 = call %IO* @IO.out_string(%IO* %6, i8* %7)
	%9 = getelementptr [11 x i8*], [11 x i8*]* @vtable.LinkedList, i32 0, i32 0
	%10 = load i8*, i8** %9
	%11 = bitcast i8* %10 to %Object* (%Object*)*
	%12 = bitcast %LinkedList* %self to %Object*
	%13 = call %Object* %11(%Object* %12)
	%14 = getelementptr %Object, %Object* null, i32 1
	%15 = ptrtoint %Object* %14 to i64
	%16 = call i8* @malloc(i64 %15)
	%17 = bitcast i8* %16 to %Object*
	%18 = getelementptr %Object, %Object* %17, i32 0, i32 0
	%19 = bitcast [3 x i8*]* @vtable.Object to i8*
	store i8* %19, i8** %18
	br label %if.end.5

if.else.5:
	%20 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	%21 = load i32, i32* %20
	%22 = icmp eq i32 %21, 1
	br i1 %22, label %if.then.6, label %if.else.6

if.end.5:
	%23 = phi %Object* [ %17, %if.then.5 ], [ %46, %if.end.6 ]
	%24 = bitcast %Object* %23 to %Object*
	ret %Object* %24

if.then.6:
	%25 = alloca %Object*
	%26 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 1
	%27 = load %Node*, %Node** %26
	%28 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 2
	%29 = load i8*, i8** %28
	%30 = bitcast i8* %29 to %Object* (%Node*)*
	%31 = bitcast %Node* %27 to %Node*
	%32 = call %Object* %30(%Node* %31)
	store %Object* %32, %Object** %25
	%33 = alloca %Node*
	store %Node* null, %Node** %33
	%34 = load %Node*, %Node** %33
	%35 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 1
	store %Node* %34, %Node** %35
	%36 = load %Node*, %Node** %33
	%37 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 2
	store %Node* %36, %Node** %37
	%38 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	store i32 0, i32* %38
	%39 = load %Object*, %Object** %25
	br label %if.end.6

if.else.6:
	%40 = alloca %Node*
	%41 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 1
	%42 = load %Node*, %Node** %41
	store %Node* %42, %Node** %40
	%43 = alloca i32
	store i32 0, i32* %43
	%44 = alloca %Object*
	store %Object* null, %Object** %44
	%45 = alloca %Node*
	store %Node* null, %Node** %45
	br label %while.cond.1

if.end.6:
	%46 = phi %Object* [ %39, %if.then.6 ], [ %85, %while.exit.1 ]
	br label %if.end.5

while.cond.1:
	%47 = load i32, i32* %43
	%48 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	%49 = load i32, i32* %48
	%50 = sub i32 %49, 2
	%51 = icmp slt i32 %47, %50
	br i1 %51, label %while.body.1, label %while.exit.1

while.body.1:
	%52 = load %Node*, %Node** %40
	%53 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 3
	%54 = load i8*, i8** %53
	%55 = bitcast i8* %54 to %Node* (%Node*)*
	%56 = bitcast %Node* %52 to %Node*
	%57 = call %Node* %55(%Node* %56)
	store %Node* %57, %Node** %40
	%58 = load i32, i32* %43
	%59 = add i32 %58, 1
	store i32 %59, i32* %43
	br label %while.cond.1

while.exit.1:
	%60 = load %Node*, %Node** %40
	%61 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 3
	%62 = load i8*, i8** %61
	%63 = bitcast i8* %62 to %Node* (%Node*)*
	%64 = bitcast %Node* %60 to %Node*
	%65 = call %Node* %63(%Node* %64)
	%66 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 2
	%67 = load i8*, i8** %66
	%68 = bitcast i8* %67 to %Object* (%Node*)*
	%69 = bitcast %Node* %65 to %Node*
	%70 = call %Object* %68(%Node* %69)
	store %Object* %70, %Object** %44
	%71 = load %Node*, %Node** %40
	%72 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 2
	store %Node* %71, %Node** %72
	%73 = load %Node*, %Node** %40
	%74 = load %Node*, %Node** %45
	%75 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 5
	%76 = load i8*, i8** %75
	%77 = bitcast i8* %76 to %Node* (%Node*, %Node*)*
	%78 = bitcast %Node* %73 to %Node*
	%79 = bitcast %Node* %74 to %Node*
	%80 = call %Node* %77(%Node* %78, %Node* %79)
	%81 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	%82 = load i32, i32* %81
	%83 = sub i32 %82, 1
	%84 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	store i32 %83, i32* %84
	%85 = load %Object*, %Object** %44
	br label %if.end.6
}

define %Object* @LinkedList.get(%LinkedList* %self, i32 %index) {
entry:
	%0 = icmp slt i32 %index, 0
	br i1 %0, label %if.then.7, label %if.else.7

if.then.7:
	%1 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 4
	%2 = load %IO*, %IO** %1
	%3 = bitcast i8* getelementptr ([45 x i8], [45 x i8]* @.str11, i32 0, i32 0) to i8*
	%4 = call %IO* @IO.out_string(%IO* %2, i8* %3)
	%5 = getelementptr [11 x i8*], [11 x i8*]* @vtable.LinkedList, i32 0, i32 0
	%6 = load i8*, i8** %5
	%7 = bitcast i8* %6 to %Object* (%Object*)*
	%8 = bitcast %LinkedList* %self to %Object*
	%9 = call %Object* %7(%Object* %8)
	%10 = getelementptr %Object, %Object* null, i32 1
	%11 = ptrtoint %Object* %10 to i64
	%12 = call i8* @malloc(i64 %11)
	%13 = bitcast i8* %12 to %Object*
	%14 = getelementptr %Object, %Object* %13, i32 0, i32 0
	%15 = bitcast [3 x i8*]* @vtable.Object to i8*
	store i8* %15, i8** %14
	br label %if.end.7

if.else.7:
	%16 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 3
	%17 = load i32, i32* %16
	%18 = icmp sle i32 %17, %index
	br i1 %18, label %if.then.8, label %if.else.8

if.end.7:
	%19 = phi %Object* [ %13, %if.then.7 ], [ %40, %if.end.8 ]
	%20 = bitcast %Object* %19 to %Object*
	ret %Object* %20

if.then.8:
	%21 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 4
	%22 = load %IO*, %IO** %21
	%23 = bitcast i8* getelementptr ([46 x i8], [46 x i8]* @.str12, i32 0, i32 0) to i8*
	%24 = call %IO* @IO.out_string(%IO* %22, i8* %23)
	%25 = getelementptr [11 x i8*], [11 x i8*]* @vtable.LinkedList, i32 0, i32 0
	%26 = load i8*, i8** %25
	%27 = bitcast i8* %26 to %Object* (%Object*)*
	%28 = bitcast %LinkedList* %self to %Object*
	%29 = call %Object* %27(%Object* %28)
	%30 = getelementptr %Object, %Object* null, i32 1
	%31 = ptrtoint %Object* %30 to i64
	%32 = call i8* @malloc(i64 %31)
	%33 = bitcast i8* %32 to %Object*
	%34 = getelementptr %Object, %Object* %33, i32 0, i32 0
	%35 = bitcast [3 x i8*]* @vtable.Object to i8*
	store i8* %35, i8** %34
	br label %if.end.8

if.else.8:
	%36 = alloca %Node*
	%37 = getelementptr %LinkedList, %LinkedList* %self, i32 0, i32 1
	%38 = load %Node*, %Node** %37
	store %Node* %38, %Node** %36
	%39 = alloca i32
	store i32 0, i32* %39
	br label %while.cond.2

if.end.8:
	%40 = phi %Object* [ %33, %if.then.8 ], [ %56, %while.exit.2 ]
	br label %if.end.7

while.cond.2:
	%41 = load i32, i32* %39
	%42 = icmp slt i32 %41, %index
	br i1 %42, label %while.body.2, label %while.exit.2

while.body.2:
	%43 = load %Node*, %Node** %36
	%44 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 3
	%45 = load i8*, i8** %44
	%46 = bitcast i8* %45 to %Node* (%Node*)*
	%47 = bitcast %Node* %43 to %Node*
	%48 = call %Node* %46(%Node* %47)
	store %Node* %48, %Node** %36
	%49 = load i32, i32* %39
	%50 = add i32 %49, 1
	store i32 %50, i32* %39
	br label %while.cond.2

while.exit.2:
	%51 = load %Node*, %Node** %36
	%52 = getelementptr [7 x i8*], [7 x i8*]* @vtable.Node, i32 0, i32 2
	%53 = load i8*, i8** %52
	%54 = bitcast i8* %53 to %Object* (%Node*)*
	%55 = bitcast %Node* %51 to %Node*
	%56 = call %Object* %54(%Node* %55)
	br label %if.end.8
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
	%6 = bitcast [11 x i8*]* @vtable.LinkedList to i8*
	store i8* %6, i8** %5
	%7 = getelementptr %LinkedList, %LinkedList* %4, i32 0, i32 1
	store %Node* null, %Node** %7
	%8 = getelementptr %LinkedList, %LinkedList* %4, i32 0, i32 2
	store %Node* null, %Node** %8
	%9 = getelementptr %LinkedList, %LinkedList* %4, i32 0, i32 3
	store i32 0, i32* %9
	%10 = getelementptr %LinkedList, %LinkedList* %4, i32 0, i32 4
	%11 = getelementptr %IO, %IO* null, i32 1
	%12 = ptrtoint %IO* %11 to i64
	%13 = call i8* @malloc(i64 %12)
	%14 = bitcast i8* %13 to %IO*
	%15 = getelementptr %IO, %IO* %14, i32 0, i32 0
	%16 = bitcast [7 x i8*]* @vtable.IO to i8*
	store i8* %16, i8** %15
	store %IO* %14, %IO** %10
	store %LinkedList* %4, %LinkedList** %0
	%17 = alloca %IO*
	%18 = getelementptr %IO, %IO* null, i32 1
	%19 = ptrtoint %IO* %18 to i64
	%20 = call i8* @malloc(i64 %19)
	%21 = bitcast i8* %20 to %IO*
	%22 = getelementptr %IO, %IO* %21, i32 0, i32 0
	%23 = bitcast [7 x i8*]* @vtable.IO to i8*
	store i8* %23, i8** %22
	store %IO* %21, %IO** %17
	%24 = load %LinkedList*, %LinkedList** %0
	%25 = getelementptr [11 x i8*], [11 x i8*]* @vtable.LinkedList, i32 0, i32 5
	%26 = load i8*, i8** %25
	%27 = bitcast i8* %26 to %LinkedList* (%LinkedList*)*
	%28 = bitcast %LinkedList* %24 to %LinkedList*
	%29 = call %LinkedList* %27(%LinkedList* %28)
	%30 = load %LinkedList*, %LinkedList** %0
	%31 = getelementptr [11 x i8*], [11 x i8*]* @vtable.LinkedList, i32 0, i32 1
	%32 = load i8*, i8** %31
	%33 = bitcast i8* %32 to %LinkedList* (%LinkedList*, %Object*)*
	%34 = bitcast %LinkedList* %30 to %LinkedList*
	%35 = inttoptr i32 100 to i8*
	%36 = bitcast i8* %35 to %Object*
	%37 = call %LinkedList* %33(%LinkedList* %34, %Object* %36)
	%38 = load %LinkedList*, %LinkedList** %0
	%39 = getelementptr [11 x i8*], [11 x i8*]* @vtable.LinkedList, i32 0, i32 2
	%40 = load i8*, i8** %39
	%41 = bitcast i8* %40 to %LinkedList* (%LinkedList*, %Object*)*
	%42 = bitcast %LinkedList* %38 to %LinkedList*
	%43 = inttoptr i32 200 to i8*
	%44 = bitcast i8* %43 to %Object*
	%45 = call %LinkedList* %41(%LinkedList* %42, %Object* %44)
	%46 = load %LinkedList*, %LinkedList** %0
	%47 = getelementptr [11 x i8*], [11 x i8*]* @vtable.LinkedList, i32 0, i32 1
	%48 = load i8*, i8** %47
	%49 = bitcast i8* %48 to %LinkedList* (%LinkedList*, %Object*)*
	%50 = bitcast %LinkedList* %46 to %LinkedList*
	%51 = inttoptr i32 50 to i8*
	%52 = bitcast i8* %51 to %Object*
	%53 = call %LinkedList* %49(%LinkedList* %50, %Object* %52)
	%54 = load %IO*, %IO** %17
	%55 = bitcast i8* getelementptr ([12 x i8], [12 x i8]* @.str13, i32 0, i32 0) to i8*
	%56 = call %IO* @IO.out_string(%IO* %54, i8* %55)
	%57 = load %IO*, %IO** %17
	%58 = load %LinkedList*, %LinkedList** %0
	%59 = getelementptr [11 x i8*], [11 x i8*]* @vtable.LinkedList, i32 0, i32 9
	%60 = load i8*, i8** %59
	%61 = bitcast i8* %60 to i32 (%LinkedList*)*
	%62 = bitcast %LinkedList* %58 to %LinkedList*
	%63 = call i32 %61(%LinkedList* %62)
	%64 = call %IO* @IO.out_int(%IO* %57, i32 %63)
	%65 = load %IO*, %IO** %17
	%66 = bitcast i8* getelementptr ([2 x i8], [2 x i8]* @.str14, i32 0, i32 0) to i8*
	%67 = call %IO* @IO.out_string(%IO* %65, i8* %66)
	%68 = load %IO*, %IO** %17
	%69 = bitcast i8* getelementptr ([16 x i8], [16 x i8]* @.str15, i32 0, i32 0) to i8*
	%70 = call %IO* @IO.out_string(%IO* %68, i8* %69)
	%71 = load %IO*, %IO** %17
	%72 = load %LinkedList*, %LinkedList** %0
	%73 = getelementptr [11 x i8*], [11 x i8*]* @vtable.LinkedList, i32 0, i32 4
	%74 = load i8*, i8** %73
	%75 = bitcast i8* %74 to %Object* (%LinkedList*, i32)*
	%76 = bitcast %LinkedList* %72 to %LinkedList*
	%77 = call %Object* %75(%LinkedList* %76, i32 0)
	%78 = icmp eq %Object* %77, null
	%79 = call %IO* @IO.out_int(%IO* %71, i32 %81)
	br i1 %78, label %case.nomatch.1, label %case.notnull.1

case.end.1:
	%80 = phi %Int* [ %97, %case.branch.0.1 ]
	%81 = ptrtoint %Int* %80 to i32
	%82 = load %IO*, %IO** %17
	%83 = bitcast i8* getelementptr ([2 x i8], [2 x i8]* @.str16, i32 0, i32 0) to i8*
	%84 = call %IO* @IO.out_string(%IO* %82, i8* %83)
	%85 = load %IO*, %IO** %17
	%86 = bitcast i8* getelementptr ([10 x i8], [10 x i8]* @.str17, i32 0, i32 0) to i8*
	%87 = call %IO* @IO.out_string(%IO* %85, i8* %86)
	%88 = load %IO*, %IO** %17
	%89 = load %LinkedList*, %LinkedList** %0
	%90 = getelementptr [11 x i8*], [11 x i8*]* @vtable.LinkedList, i32 0, i32 7
	%91 = load i8*, i8** %90
	%92 = bitcast i8* %91 to %Object* (%LinkedList*)*
	%93 = bitcast %LinkedList* %89 to %LinkedList*
	%94 = call %Object* %92(%LinkedList* %93)
	%95 = icmp eq %Object* %94, null
	%96 = call %IO* @IO.out_int(%IO* %88, i32 %99)
	br i1 %95, label %case.nomatch.2, label %case.notnull.2

case.branch.0.1:
	%97 = bitcast %Object* %77 to %Int*
	br label %case.end.1

case.nomatch.1:
	call void @exit(i32 1)
	unreachable

case.notnull.1:
	br label %case.typecheck.1

case.typecheck.1:
	br label %case.decision.0.1

case.decision.0.1:
	br i1 false, label %case.branch.0.1, label %case.nomatch.1

case.end.2:
	%98 = phi %Int* [ %118, %case.branch.0.2 ]
	%99 = ptrtoint %Int* %98 to i32
	%100 = load %IO*, %IO** %17
	%101 = bitcast i8* getelementptr ([2 x i8], [2 x i8]* @.str18, i32 0, i32 0) to i8*
	%102 = call %IO* @IO.out_string(%IO* %100, i8* %101)
	%103 = load %IO*, %IO** %17
	%104 = bitcast i8* getelementptr ([11 x i8], [11 x i8]* @.str19, i32 0, i32 0) to i8*
	%105 = call %IO* @IO.out_string(%IO* %103, i8* %104)
	%106 = load %IO*, %IO** %17
	%107 = load %LinkedList*, %LinkedList** %0
	%108 = getelementptr [11 x i8*], [11 x i8*]* @vtable.LinkedList, i32 0, i32 9
	%109 = load i8*, i8** %108
	%110 = bitcast i8* %109 to i32 (%LinkedList*)*
	%111 = bitcast %LinkedList* %107 to %LinkedList*
	%112 = call i32 %110(%LinkedList* %111)
	%113 = call %IO* @IO.out_int(%IO* %106, i32 %112)
	%114 = load %IO*, %IO** %17
	%115 = bitcast i8* getelementptr ([2 x i8], [2 x i8]* @.str20, i32 0, i32 0) to i8*
	%116 = call %IO* @IO.out_string(%IO* %114, i8* %115)
	%117 = bitcast %IO* %114 to %Object*
	ret %Object* %117

case.branch.0.2:
	%118 = bitcast %Object* %94 to %Int*
	br label %case.end.2

case.nomatch.2:
	call void @exit(i32 1)
	unreachable

case.notnull.2:
	br label %case.typecheck.2

case.typecheck.2:
	br label %case.decision.0.2

case.decision.0.2:
	br i1 false, label %case.branch.0.2, label %case.nomatch.2
}

define i32 @main() {
entry:
	%0 = getelementptr %Main, %Main* null, i32 1
	%1 = ptrtoint %Main* %0 to i64
	%2 = call i8* @malloc(i64 %1)
	%3 = bitcast i8* %2 to %Main*
	%4 = getelementptr %Main, %Main* %3, i32 0, i32 0
	%5 = bitcast [4 x i8*]* @vtable.Main to i8*
	store i8* %5, i8** %4
	%6 = alloca %Main*
	store %Main* %3, %Main** %6
	%7 = getelementptr [4 x i8*], [4 x i8*]* @vtable.Main, i32 0, i32 2
	%8 = load i8*, i8** %7
	%9 = bitcast i8* %8 to i8* (%Main*)*
	%10 = call i8* %9(%Main* %3)
	ret i32 0
}
