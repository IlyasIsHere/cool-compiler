%Object = type { i8* }
%IO = type { i8* }
%Int = type { i8* }
%String = type { i8* }
%Bool = type { i8* }
%MutualRecursion = type { i8* }
%Main = type { i8*, %IO* }

@.str.empty = constant [1 x i8] c"\00"
@vtable.Object = global [3 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.IO = global [7 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%IO*)* @IO.in_int to i8*), i8* bitcast (i8* (%IO*)* @IO.in_string to i8*), i8* bitcast (%IO* (%IO*, i32)* @IO.out_int to i8*), i8* bitcast (%IO* (%IO*, i8*)* @IO.out_string to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Int = global [3 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.String = global [6 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (i8* (%String*, i8*)* @String.concat to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%String*)* @String.length to i8*), i8* bitcast (i8* (%String*, i32, i32)* @String.substr to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Bool = global [3 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.MutualRecursion = global [3 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Main = global [10 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Main*, i32)* @Main.get_indent to i8*), i8* bitcast (i1 (%Main*, i32)* @Main.is_even to i8*), i8* bitcast (i1 (%Main*, i32, i32)* @Main.is_even_trace to i8*), i8* bitcast (i1 (%Main*, i32)* @Main.is_odd to i8*), i8* bitcast (i1 (%Main*, i32, i32)* @Main.is_odd_trace to i8*), i8* bitcast (%Object* (%Main*)* @Main.main to i8*), i8* bitcast (%Object* (%Main*, i32)* @Main.test_number to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@.str8 = internal constant [45 x i8] c"Mutual Recursion Demo (Even/Odd functions)\0A\0A\00"
@.str9 = internal constant [23 x i8] c"Testing even numbers:\0A\00"
@.str10 = internal constant [23 x i8] c"\0ATesting odd numbers:\0A\00"
@.str11 = internal constant [27 x i8] c"\0ATesting a larger number:\0A\00"
@.str12 = internal constant [32 x i8] c"\0ARecursive trace for number 5:\0A\00"
@.str13 = internal constant [13 x i8] c"is_even(5):\0A\00"
@.str14 = internal constant [5 x i8] c" is \00"
@.str15 = internal constant [5 x i8] c"even\00"
@.str16 = internal constant [4 x i8] c"odd\00"
@.str17 = internal constant [2 x i8] c"\0A\00"
@.str18 = internal constant [9 x i8] c"is_even(\00"
@.str19 = internal constant [3 x i8] c")\0A\00"
@.str20 = internal constant [38 x i8] c"  return true (Base case: 0 is even)\0A\00"
@.str21 = internal constant [38 x i8] c"  return false (Base case: 1 is odd)\0A\00"
@.str22 = internal constant [17 x i8] c"  return is_odd(\00"
@.str23 = internal constant [3 x i8] c")\0A\00"
@.str24 = internal constant [11 x i8] c"  is_even(\00"
@.str25 = internal constant [11 x i8] c") returns \00"
@.str26 = internal constant [6 x i8] c"true\0A\00"
@.str27 = internal constant [7 x i8] c"false\0A\00"
@.str28 = internal constant [8 x i8] c"is_odd(\00"
@.str29 = internal constant [3 x i8] c")\0A\00"
@.str30 = internal constant [42 x i8] c"  return false (Base case: 0 is not odd)\0A\00"
@.str31 = internal constant [37 x i8] c"  return true (Base case: 1 is odd)\0A\00"
@.str32 = internal constant [18 x i8] c"  return is_even(\00"
@.str33 = internal constant [3 x i8] c")\0A\00"
@.str34 = internal constant [10 x i8] c"  is_odd(\00"
@.str35 = internal constant [11 x i8] c") returns \00"
@.str36 = internal constant [6 x i8] c"true\0A\00"
@.str37 = internal constant [7 x i8] c"false\0A\00"
@.str38 = internal constant [1 x i8] c"\00"
@.str39 = internal constant [3 x i8] c"  \00"
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
	%6 = getelementptr { i8*, i8* }, %String* %self, i32 0, i32 1
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
	%13 = getelementptr { i8*, i8* }, %String* %self, i32 0, i32 1
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
	%6 = getelementptr { i8*, i8* }, %String* %self, i32 0, i32 1
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

define %Object* @Main.main(%Main* %self) {
entry:
	%0 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%1 = load %IO*, %IO** %0
	%2 = call %IO* @IO.out_string(%IO* %1, i8* getelementptr ([45 x i8], [45 x i8]* @.str8, i32 0, i32 0))
	%3 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%4 = load %IO*, %IO** %3
	%5 = call %IO* @IO.out_string(%IO* %4, i8* getelementptr ([23 x i8], [23 x i8]* @.str9, i32 0, i32 0))
	%6 = call %Object* @Main.test_number(%Main* %self, i32 0)
	%7 = call %Object* @Main.test_number(%Main* %self, i32 2)
	%8 = call %Object* @Main.test_number(%Main* %self, i32 4)
	%9 = call %Object* @Main.test_number(%Main* %self, i32 6)
	%10 = call %Object* @Main.test_number(%Main* %self, i32 10)
	%11 = call %Object* @Main.test_number(%Main* %self, i32 20)
	%12 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%13 = load %IO*, %IO** %12
	%14 = call %IO* @IO.out_string(%IO* %13, i8* getelementptr ([23 x i8], [23 x i8]* @.str10, i32 0, i32 0))
	%15 = call %Object* @Main.test_number(%Main* %self, i32 1)
	%16 = call %Object* @Main.test_number(%Main* %self, i32 3)
	%17 = call %Object* @Main.test_number(%Main* %self, i32 5)
	%18 = call %Object* @Main.test_number(%Main* %self, i32 7)
	%19 = call %Object* @Main.test_number(%Main* %self, i32 11)
	%20 = call %Object* @Main.test_number(%Main* %self, i32 21)
	%21 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%22 = load %IO*, %IO** %21
	%23 = call %IO* @IO.out_string(%IO* %22, i8* getelementptr ([27 x i8], [27 x i8]* @.str11, i32 0, i32 0))
	%24 = call %Object* @Main.test_number(%Main* %self, i32 42)
	%25 = call %Object* @Main.test_number(%Main* %self, i32 99)
	%26 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%27 = load %IO*, %IO** %26
	%28 = call %IO* @IO.out_string(%IO* %27, i8* getelementptr ([32 x i8], [32 x i8]* @.str12, i32 0, i32 0))
	%29 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%30 = load %IO*, %IO** %29
	%31 = call %IO* @IO.out_string(%IO* %30, i8* getelementptr ([13 x i8], [13 x i8]* @.str13, i32 0, i32 0))
	%32 = call i1 @Main.is_even_trace(%Main* %self, i32 5, i32 1)
	%33 = inttoptr i1 %32 to %Object*
	ret %Object* %33
}

define %Object* @Main.test_number(%Main* %self, i32 %n) {
entry:
	%0 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%1 = load %IO*, %IO** %0
	%2 = call %IO* @IO.out_int(%IO* %1, i32 %n)
	%3 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%4 = load %IO*, %IO** %3
	%5 = call %IO* @IO.out_string(%IO* %4, i8* getelementptr ([5 x i8], [5 x i8]* @.str14, i32 0, i32 0))
	%6 = call i1 @Main.is_even(%Main* %self, i32 %n)
	br i1 %6, label %if.then.1, label %if.else.1

if.then.1:
	%7 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%8 = load %IO*, %IO** %7
	%9 = call %IO* @IO.out_string(%IO* %8, i8* getelementptr ([5 x i8], [5 x i8]* @.str15, i32 0, i32 0))
	br label %if.end.1

if.else.1:
	%10 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%11 = load %IO*, %IO** %10
	%12 = call %IO* @IO.out_string(%IO* %11, i8* getelementptr ([4 x i8], [4 x i8]* @.str16, i32 0, i32 0))
	br label %if.end.1

if.end.1:
	%13 = phi %IO* [ %8, %if.then.1 ], [ %11, %if.else.1 ]
	%14 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%15 = load %IO*, %IO** %14
	%16 = call %IO* @IO.out_string(%IO* %15, i8* getelementptr ([2 x i8], [2 x i8]* @.str17, i32 0, i32 0))
	%17 = bitcast %IO* %15 to %Object*
	ret %Object* %17
}

define i1 @Main.is_even(%Main* %self, i32 %n) {
entry:
	%0 = icmp eq i32 %n, 0
	br i1 %0, label %if.then.2, label %if.else.2

if.then.2:
	br label %if.end.2

if.else.2:
	%1 = icmp eq i32 %n, 1
	br i1 %1, label %if.then.3, label %if.else.3

if.end.2:
	%2 = phi i1 [ true, %if.then.2 ], [ %5, %if.end.3 ]
	ret i1 %2

if.then.3:
	br label %if.end.3

if.else.3:
	%3 = sub i32 %n, 1
	%4 = call i1 @Main.is_odd(%Main* %self, i32 %3)
	br label %if.end.3

if.end.3:
	%5 = phi i1 [ false, %if.then.3 ], [ %4, %if.else.3 ]
	br label %if.end.2
}

define i1 @Main.is_odd(%Main* %self, i32 %n) {
entry:
	%0 = icmp eq i32 %n, 0
	br i1 %0, label %if.then.4, label %if.else.4

if.then.4:
	br label %if.end.4

if.else.4:
	%1 = icmp eq i32 %n, 1
	br i1 %1, label %if.then.5, label %if.else.5

if.end.4:
	%2 = phi i1 [ false, %if.then.4 ], [ %5, %if.end.5 ]
	ret i1 %2

if.then.5:
	br label %if.end.5

if.else.5:
	%3 = sub i32 %n, 1
	%4 = call i1 @Main.is_even(%Main* %self, i32 %3)
	br label %if.end.5

if.end.5:
	%5 = phi i1 [ true, %if.then.5 ], [ %4, %if.else.5 ]
	br label %if.end.4
}

define i1 @Main.is_even_trace(%Main* %self, i32 %n, i32 %level) {
entry:
	%0 = alloca i8*
	%1 = call i8* @Main.get_indent(%Main* %self, i32 %level)
	store i8* %1, i8** %0
	%2 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%3 = load %IO*, %IO** %2
	%4 = load i8*, i8** %0
	%5 = call %IO* @IO.out_string(%IO* %3, i8* %4)
	%6 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%7 = load %IO*, %IO** %6
	%8 = call %IO* @IO.out_string(%IO* %7, i8* getelementptr ([9 x i8], [9 x i8]* @.str18, i32 0, i32 0))
	%9 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%10 = load %IO*, %IO** %9
	%11 = call %IO* @IO.out_int(%IO* %10, i32 %n)
	%12 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%13 = load %IO*, %IO** %12
	%14 = call %IO* @IO.out_string(%IO* %13, i8* getelementptr ([3 x i8], [3 x i8]* @.str19, i32 0, i32 0))
	%15 = icmp eq i32 %n, 0
	br i1 %15, label %if.then.6, label %if.else.6

if.then.6:
	%16 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%17 = load %IO*, %IO** %16
	%18 = load i8*, i8** %0
	%19 = call %IO* @IO.out_string(%IO* %17, i8* %18)
	%20 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%21 = load %IO*, %IO** %20
	%22 = call %IO* @IO.out_string(%IO* %21, i8* getelementptr ([38 x i8], [38 x i8]* @.str20, i32 0, i32 0))
	br label %if.end.6

if.else.6:
	%23 = icmp eq i32 %n, 1
	br i1 %23, label %if.then.7, label %if.else.7

if.end.6:
	%24 = phi i1 [ true, %if.then.6 ], [ %64, %if.end.7 ]
	ret i1 %24

if.then.7:
	%25 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%26 = load %IO*, %IO** %25
	%27 = load i8*, i8** %0
	%28 = call %IO* @IO.out_string(%IO* %26, i8* %27)
	%29 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%30 = load %IO*, %IO** %29
	%31 = call %IO* @IO.out_string(%IO* %30, i8* getelementptr ([38 x i8], [38 x i8]* @.str21, i32 0, i32 0))
	br label %if.end.7

if.else.7:
	%32 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%33 = load %IO*, %IO** %32
	%34 = load i8*, i8** %0
	%35 = call %IO* @IO.out_string(%IO* %33, i8* %34)
	%36 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%37 = load %IO*, %IO** %36
	%38 = call %IO* @IO.out_string(%IO* %37, i8* getelementptr ([17 x i8], [17 x i8]* @.str22, i32 0, i32 0))
	%39 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%40 = load %IO*, %IO** %39
	%41 = sub i32 %n, 1
	%42 = call %IO* @IO.out_int(%IO* %40, i32 %41)
	%43 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%44 = load %IO*, %IO** %43
	%45 = call %IO* @IO.out_string(%IO* %44, i8* getelementptr ([3 x i8], [3 x i8]* @.str23, i32 0, i32 0))
	%46 = alloca i1
	%47 = sub i32 %n, 1
	%48 = add i32 %level, 1
	%49 = call i1 @Main.is_odd_trace(%Main* %self, i32 %47, i32 %48)
	store i1 %49, i1* %46
	%50 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%51 = load %IO*, %IO** %50
	%52 = load i8*, i8** %0
	%53 = call %IO* @IO.out_string(%IO* %51, i8* %52)
	%54 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%55 = load %IO*, %IO** %54
	%56 = call %IO* @IO.out_string(%IO* %55, i8* getelementptr ([11 x i8], [11 x i8]* @.str24, i32 0, i32 0))
	%57 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%58 = load %IO*, %IO** %57
	%59 = call %IO* @IO.out_int(%IO* %58, i32 %n)
	%60 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%61 = load %IO*, %IO** %60
	%62 = call %IO* @IO.out_string(%IO* %61, i8* getelementptr ([11 x i8], [11 x i8]* @.str25, i32 0, i32 0))
	%63 = load i1, i1* %46
	br i1 %63, label %if.then.8, label %if.else.8

if.end.7:
	%64 = phi i1 [ false, %if.then.7 ], [ %72, %if.end.8 ]
	br label %if.end.6

if.then.8:
	%65 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%66 = load %IO*, %IO** %65
	%67 = call %IO* @IO.out_string(%IO* %66, i8* getelementptr ([6 x i8], [6 x i8]* @.str26, i32 0, i32 0))
	br label %if.end.8

if.else.8:
	%68 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%69 = load %IO*, %IO** %68
	%70 = call %IO* @IO.out_string(%IO* %69, i8* getelementptr ([7 x i8], [7 x i8]* @.str27, i32 0, i32 0))
	br label %if.end.8

if.end.8:
	%71 = phi %IO* [ %66, %if.then.8 ], [ %69, %if.else.8 ]
	%72 = load i1, i1* %46
	br label %if.end.7
}

define i1 @Main.is_odd_trace(%Main* %self, i32 %n, i32 %level) {
entry:
	%0 = alloca i8*
	%1 = call i8* @Main.get_indent(%Main* %self, i32 %level)
	store i8* %1, i8** %0
	%2 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%3 = load %IO*, %IO** %2
	%4 = load i8*, i8** %0
	%5 = call %IO* @IO.out_string(%IO* %3, i8* %4)
	%6 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%7 = load %IO*, %IO** %6
	%8 = call %IO* @IO.out_string(%IO* %7, i8* getelementptr ([8 x i8], [8 x i8]* @.str28, i32 0, i32 0))
	%9 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%10 = load %IO*, %IO** %9
	%11 = call %IO* @IO.out_int(%IO* %10, i32 %n)
	%12 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%13 = load %IO*, %IO** %12
	%14 = call %IO* @IO.out_string(%IO* %13, i8* getelementptr ([3 x i8], [3 x i8]* @.str29, i32 0, i32 0))
	%15 = icmp eq i32 %n, 0
	br i1 %15, label %if.then.9, label %if.else.9

if.then.9:
	%16 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%17 = load %IO*, %IO** %16
	%18 = load i8*, i8** %0
	%19 = call %IO* @IO.out_string(%IO* %17, i8* %18)
	%20 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%21 = load %IO*, %IO** %20
	%22 = call %IO* @IO.out_string(%IO* %21, i8* getelementptr ([42 x i8], [42 x i8]* @.str30, i32 0, i32 0))
	br label %if.end.9

if.else.9:
	%23 = icmp eq i32 %n, 1
	br i1 %23, label %if.then.10, label %if.else.10

if.end.9:
	%24 = phi i1 [ false, %if.then.9 ], [ %64, %if.end.10 ]
	ret i1 %24

if.then.10:
	%25 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%26 = load %IO*, %IO** %25
	%27 = load i8*, i8** %0
	%28 = call %IO* @IO.out_string(%IO* %26, i8* %27)
	%29 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%30 = load %IO*, %IO** %29
	%31 = call %IO* @IO.out_string(%IO* %30, i8* getelementptr ([37 x i8], [37 x i8]* @.str31, i32 0, i32 0))
	br label %if.end.10

if.else.10:
	%32 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%33 = load %IO*, %IO** %32
	%34 = load i8*, i8** %0
	%35 = call %IO* @IO.out_string(%IO* %33, i8* %34)
	%36 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%37 = load %IO*, %IO** %36
	%38 = call %IO* @IO.out_string(%IO* %37, i8* getelementptr ([18 x i8], [18 x i8]* @.str32, i32 0, i32 0))
	%39 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%40 = load %IO*, %IO** %39
	%41 = sub i32 %n, 1
	%42 = call %IO* @IO.out_int(%IO* %40, i32 %41)
	%43 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%44 = load %IO*, %IO** %43
	%45 = call %IO* @IO.out_string(%IO* %44, i8* getelementptr ([3 x i8], [3 x i8]* @.str33, i32 0, i32 0))
	%46 = alloca i1
	%47 = sub i32 %n, 1
	%48 = add i32 %level, 1
	%49 = call i1 @Main.is_even_trace(%Main* %self, i32 %47, i32 %48)
	store i1 %49, i1* %46
	%50 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%51 = load %IO*, %IO** %50
	%52 = load i8*, i8** %0
	%53 = call %IO* @IO.out_string(%IO* %51, i8* %52)
	%54 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%55 = load %IO*, %IO** %54
	%56 = call %IO* @IO.out_string(%IO* %55, i8* getelementptr ([10 x i8], [10 x i8]* @.str34, i32 0, i32 0))
	%57 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%58 = load %IO*, %IO** %57
	%59 = call %IO* @IO.out_int(%IO* %58, i32 %n)
	%60 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%61 = load %IO*, %IO** %60
	%62 = call %IO* @IO.out_string(%IO* %61, i8* getelementptr ([11 x i8], [11 x i8]* @.str35, i32 0, i32 0))
	%63 = load i1, i1* %46
	br i1 %63, label %if.then.11, label %if.else.11

if.end.10:
	%64 = phi i1 [ true, %if.then.10 ], [ %72, %if.end.11 ]
	br label %if.end.9

if.then.11:
	%65 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%66 = load %IO*, %IO** %65
	%67 = call %IO* @IO.out_string(%IO* %66, i8* getelementptr ([6 x i8], [6 x i8]* @.str36, i32 0, i32 0))
	br label %if.end.11

if.else.11:
	%68 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%69 = load %IO*, %IO** %68
	%70 = call %IO* @IO.out_string(%IO* %69, i8* getelementptr ([7 x i8], [7 x i8]* @.str37, i32 0, i32 0))
	br label %if.end.11

if.end.11:
	%71 = phi %IO* [ %66, %if.then.11 ], [ %69, %if.else.11 ]
	%72 = load i1, i1* %46
	br label %if.end.10
}

define i8* @Main.get_indent(%Main* %self, i32 %level) {
entry:
	%0 = alloca i8*
	store i8* getelementptr ([1 x i8], [1 x i8]* @.str38, i32 0, i32 0), i8** %0
	%1 = alloca i32
	store i32 0, i32* %1
	br label %while.cond.1

while.cond.1:
	%2 = load i32, i32* %1
	%3 = icmp slt i32 %2, %level
	br i1 %3, label %while.body.1, label %while.exit.1

while.body.1:
	%4 = load i8*, i8** %0
	%5 = call i8* @String.concat(i8* %4, i8* getelementptr ([3 x i8], [3 x i8]* @.str39, i32 0, i32 0))
	store i8* %5, i8** %0
	%6 = load i32, i32* %1
	%7 = add i32 %6, 1
	store i32 %7, i32* %1
	br label %while.cond.1

while.exit.1:
	%8 = load i8*, i8** %0
	%9 = bitcast i8* %8 to i8*
	ret i8* %9
}

define i32 @main() {
entry:
	%0 = call i8* @malloc(%Main* getelementptr (%Main, %Main* null, i32 1))
	%1 = bitcast i8* %0 to %Main*
	%2 = getelementptr %Main, %Main* %1, i32 0, i32 0
	%3 = bitcast [10 x i8*]* @vtable.Main to i8*
	store i8* %3, i8** %2
	%4 = getelementptr %Main, %Main* %1, i32 0, i32 1
	%5 = call i8* @malloc(%IO* getelementptr (%IO, %IO* null, i32 1))
	%6 = bitcast i8* %5 to %IO*
	%7 = getelementptr %IO, %IO* %6, i32 0, i32 0
	%8 = bitcast [7 x i8*]* @vtable.IO to i8*
	store i8* %8, i8** %7
	store %IO* %6, %IO** %4
	%9 = alloca %Main*
	store %Main* %1, %Main** %9
	%10 = getelementptr %Main, %Main* %1, i32 0, i32 1
	%11 = call i8* @malloc(%IO* getelementptr (%IO, %IO* null, i32 1))
	%12 = bitcast i8* %11 to %IO*
	%13 = getelementptr %IO, %IO* %12, i32 0, i32 0
	%14 = bitcast [7 x i8*]* @vtable.IO to i8*
	store i8* %14, i8** %13
	store %IO* %12, %IO** %10
	%15 = getelementptr [10 x i8*], [10 x i8*]* @vtable.Main, i32 0, i32 7
	%16 = load i8*, i8** %15
	%17 = bitcast i8* %16 to i8* (%Main*)*
	%18 = call i8* %17(%Main* %1)
	ret i32 0
}
