%Object = type { i8* }
%IO = type { i8* }
%Int = type { i8* }
%String = type { i8* }
%Bool = type { i8* }
%Main = type { i8*, i32 }

@.str.empty = constant [1 x i8] c"\00"
@vtable.Object = global [3 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.IO = global [7 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%IO*)* @IO.in_int to i8*), i8* bitcast (i8* (%IO*)* @IO.in_string to i8*), i8* bitcast (%IO* (%IO*, i32)* @IO.out_int to i8*), i8* bitcast (%IO* (%IO*, i8*)* @IO.out_string to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Int = global [3 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.String = global [6 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (i8* (%String*, i8*)* @String.concat to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%String*)* @String.length to i8*), i8* bitcast (i8* (%String*, i32, i32)* @String.substr to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Bool = global [3 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Main = global [4 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (%Object* (%Main*)* @Main.main to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@.str7 = internal constant [12 x i8] c"The number \00"
@.str8 = internal constant [18 x i8] c" is less than 20\0A\00"
@.str9 = internal constant [12 x i8] c"The number \00"
@.str10 = internal constant [33 x i8] c" is greater than or equal to 20\0A\00"
@.str11 = internal constant [27 x i8] c"The number is exactly 10!\0A\00"
@.str12 = internal constant [22 x i8] c"The number is not 10\0A\00"
@.str13 = internal constant [28 x i8] c"This is after the if hahaha\00"
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
	%1 = load i32, i32* %0
	%2 = icmp slt i32 %1, 20
	br i1 %2, label %if.then.1, label %if.else.1

if.then.1:
	%3 = call i8* @malloc(%IO* getelementptr (%IO, %IO* null, i32 1))
	%4 = bitcast i8* %3 to %IO*
	%5 = getelementptr %IO, %IO* %4, i32 0, i32 0
	%6 = bitcast [7 x i8*]* @vtable.IO to i8*
	store i8* %6, i8** %5
	%7 = call %IO* @IO.out_string(%IO* %4, i8* getelementptr ([12 x i8], [12 x i8]* @.str7, i32 0, i32 0))
	%8 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%9 = load i32, i32* %8
	%10 = call %IO* @IO.out_int(%IO* %4, i32 %9)
	%11 = call %IO* @IO.out_string(%IO* %4, i8* getelementptr ([18 x i8], [18 x i8]* @.str8, i32 0, i32 0))
	br label %if.end.1

if.else.1:
	%12 = call i8* @malloc(%IO* getelementptr (%IO, %IO* null, i32 1))
	%13 = bitcast i8* %12 to %IO*
	%14 = getelementptr %IO, %IO* %13, i32 0, i32 0
	%15 = bitcast [7 x i8*]* @vtable.IO to i8*
	store i8* %15, i8** %14
	%16 = call %IO* @IO.out_string(%IO* %13, i8* getelementptr ([12 x i8], [12 x i8]* @.str9, i32 0, i32 0))
	%17 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%18 = load i32, i32* %17
	%19 = call %IO* @IO.out_int(%IO* %13, i32 %18)
	%20 = call %IO* @IO.out_string(%IO* %13, i8* getelementptr ([33 x i8], [33 x i8]* @.str10, i32 0, i32 0))
	br label %if.end.1

if.end.1:
	%21 = phi %IO* [ %4, %if.then.1 ], [ %13, %if.else.1 ]
	%22 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%23 = load i32, i32* %22
	%24 = icmp eq i32 %23, 10
	br i1 %24, label %if.then.2, label %if.else.2

if.then.2:
	%25 = call i8* @malloc(%IO* getelementptr (%IO, %IO* null, i32 1))
	%26 = bitcast i8* %25 to %IO*
	%27 = getelementptr %IO, %IO* %26, i32 0, i32 0
	%28 = bitcast [7 x i8*]* @vtable.IO to i8*
	store i8* %28, i8** %27
	%29 = call %IO* @IO.out_string(%IO* %26, i8* getelementptr ([27 x i8], [27 x i8]* @.str11, i32 0, i32 0))
	br label %if.end.2

if.else.2:
	%30 = call i8* @malloc(%IO* getelementptr (%IO, %IO* null, i32 1))
	%31 = bitcast i8* %30 to %IO*
	%32 = getelementptr %IO, %IO* %31, i32 0, i32 0
	%33 = bitcast [7 x i8*]* @vtable.IO to i8*
	store i8* %33, i8** %32
	%34 = call %IO* @IO.out_string(%IO* %31, i8* getelementptr ([22 x i8], [22 x i8]* @.str12, i32 0, i32 0))
	br label %if.end.2

if.end.2:
	%35 = phi %IO* [ %26, %if.then.2 ], [ %31, %if.else.2 ]
	%36 = call i8* @malloc(%IO* getelementptr (%IO, %IO* null, i32 1))
	%37 = bitcast i8* %36 to %IO*
	%38 = getelementptr %IO, %IO* %37, i32 0, i32 0
	%39 = bitcast [7 x i8*]* @vtable.IO to i8*
	store i8* %39, i8** %38
	%40 = call %IO* @IO.out_string(%IO* %37, i8* getelementptr ([28 x i8], [28 x i8]* @.str13, i32 0, i32 0))
	%41 = bitcast %IO* %37 to %Object*
	ret %Object* %41
}

define i32 @main() {
entry:
	%0 = call i8* @malloc(%Main* getelementptr (%Main, %Main* null, i32 1))
	%1 = bitcast i8* %0 to %Main*
	%2 = getelementptr %Main, %Main* %1, i32 0, i32 0
	%3 = bitcast [4 x i8*]* @vtable.Main to i8*
	store i8* %3, i8** %2
	%4 = getelementptr %Main, %Main* %1, i32 0, i32 1
	store i32 10, i32* %4
	%5 = alloca %Main*
	store %Main* %1, %Main** %5
	%6 = getelementptr %Main, %Main* %1, i32 0, i32 1
	store i32 10, i32* %6
	%7 = getelementptr [4 x i8*], [4 x i8*]* @vtable.Main, i32 0, i32 2
	%8 = load i8*, i8** %7
	%9 = bitcast i8* %8 to i8* (%Main*)*
	%10 = call i8* %9(%Main* %1)
	ret i32 0
}
