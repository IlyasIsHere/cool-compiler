%Object = type { i8* }
%IO = type { i8* }
%Int = type { i8* }
%String = type { i8* }
%Bool = type { i8* }
%A = type { i8* }
%B = type { i8* }
%C = type { i8* }
%Main = type { i8* }

@.str.empty = constant [1 x i8] c"\00"
@vtable.Object = global [3 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.IO = global [7 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%IO*)* @IO.in_int to i8*), i8* bitcast (i8* (%IO*)* @IO.in_string to i8*), i8* bitcast (%IO* (%IO*, i32)* @IO.out_int to i8*), i8* bitcast (%IO* (%IO*, i8*)* @IO.out_string to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Int = global [3 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.String = global [6 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (i8* (%String*, i8*)* @String.concat to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%String*)* @String.length to i8*), i8* bitcast (i8* (%String*, i32, i32)* @String.substr to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Bool = global [3 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.A = global [6 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%A*)* @A.method1 to i8*), i8* bitcast (i32 (%A*)* @A.method2 to i8*), i8* bitcast (i32 (%A*)* @A.method3 to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.B = global [6 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%A*)* @A.method1 to i8*), i8* bitcast (i32 (%B*)* @B.method2 to i8*), i8* bitcast (i32 (%B*)* @B.method3 to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.C = global [6 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%A*)* @A.method1 to i8*), i8* bitcast (i32 (%B*)* @B.method2 to i8*), i8* bitcast (i32 (%C*)* @C.method3 to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Main = global [8 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%IO*)* @IO.in_int to i8*), i8* bitcast (i8* (%IO*)* @IO.in_string to i8*), i8* bitcast (%Object* (%Main*)* @Main.main to i8*), i8* bitcast (%IO* (%IO*, i32)* @IO.out_int to i8*), i8* bitcast (%IO* (%IO*, i8*)* @IO.out_string to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@.str10 = internal constant [2 x i8] c" \00"
@.str11 = internal constant [2 x i8] c"\0A\00"
@.str12 = internal constant [2 x i8] c" \00"
@.str13 = internal constant [2 x i8] c"\0A\00"
@.str14 = internal constant [2 x i8] c" \00"
@.str15 = internal constant [2 x i8] c"\0A\00"
@.str16 = internal constant [2 x i8] c" \00"
@.str17 = internal constant [2 x i8] c" \00"
@.str18 = internal constant [2 x i8] c"\0A\00"
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

define i32 @A.method1(%A* %self) {
entry:
	ret i32 1
}

define i32 @A.method2(%A* %self) {
entry:
	ret i32 2
}

define i32 @A.method3(%A* %self) {
entry:
	ret i32 3
}

define i32 @B.method2(%B* %self) {
entry:
	ret i32 22
}

define i32 @B.method3(%B* %self) {
entry:
	ret i32 35
}

define i32 @C.method3(%C* %self) {
entry:
	ret i32 33
}

define %Object* @Main.main(%Main* %self) {
entry:
	%0 = alloca %A*
	%1 = call i8* @malloc(%A* getelementptr (%A, %A* null, i32 1))
	%2 = bitcast i8* %1 to %A*
	%3 = getelementptr %A, %A* %2, i32 0, i32 0
	%4 = bitcast [6 x i8*]* @vtable.A to i8*
	store i8* %4, i8** %3
	store %A* %2, %A** %0
	%5 = alloca %B*
	%6 = call i8* @malloc(%B* getelementptr (%B, %B* null, i32 1))
	%7 = bitcast i8* %6 to %B*
	%8 = getelementptr %B, %B* %7, i32 0, i32 0
	%9 = bitcast [6 x i8*]* @vtable.B to i8*
	store i8* %9, i8** %8
	store %B* %7, %B** %5
	%10 = alloca %C*
	%11 = call i8* @malloc(%C* getelementptr (%C, %C* null, i32 1))
	%12 = bitcast i8* %11 to %C*
	%13 = getelementptr %C, %C* %12, i32 0, i32 0
	%14 = bitcast [6 x i8*]* @vtable.C to i8*
	store i8* %14, i8** %13
	store %C* %12, %C** %10
	%15 = load %C*, %C** %10
	%16 = call i32 @A.method1(%C* %15)
	%17 = call %IO* @IO.out_int(%Main* %self, i32 %16)
	%18 = call %IO* @IO.out_string(%IO* %17, i8* getelementptr ([2 x i8], [2 x i8]* @.str10, i32 0, i32 0))
	%19 = load %C*, %C** %10
	%20 = call i32 @A.method1(%C* %19)
	%21 = call %IO* @IO.out_int(%IO* %17, i32 %20)
	%22 = call %IO* @IO.out_string(%Main* %self, i8* getelementptr ([2 x i8], [2 x i8]* @.str11, i32 0, i32 0))
	%23 = load %C*, %C** %10
	%24 = call i32 @B.method2(%C* %23)
	%25 = call %IO* @IO.out_int(%Main* %self, i32 %24)
	%26 = call %IO* @IO.out_string(%IO* %25, i8* getelementptr ([2 x i8], [2 x i8]* @.str12, i32 0, i32 0))
	%27 = load %C*, %C** %10
	%28 = call i32 @B.method2(%C* %27)
	%29 = call %IO* @IO.out_int(%IO* %25, i32 %28)
	%30 = call %IO* @IO.out_string(%Main* %self, i8* getelementptr ([2 x i8], [2 x i8]* @.str13, i32 0, i32 0))
	%31 = load %B*, %B** %5
	%32 = call i32 @B.method2(%B* %31)
	%33 = call %IO* @IO.out_int(%Main* %self, i32 %32)
	%34 = call %IO* @IO.out_string(%IO* %33, i8* getelementptr ([2 x i8], [2 x i8]* @.str14, i32 0, i32 0))
	%35 = load %B*, %B** %5
	%36 = call i32 @A.method2(%B* %35)
	%37 = call %IO* @IO.out_int(%IO* %33, i32 %36)
	%38 = call %IO* @IO.out_string(%Main* %self, i8* getelementptr ([2 x i8], [2 x i8]* @.str15, i32 0, i32 0))
	%39 = load %C*, %C** %10
	%40 = call i32 @C.method3(%C* %39)
	%41 = call %IO* @IO.out_int(%Main* %self, i32 %40)
	%42 = call %IO* @IO.out_string(%IO* %41, i8* getelementptr ([2 x i8], [2 x i8]* @.str16, i32 0, i32 0))
	%43 = load %C*, %C** %10
	%44 = call i32 @A.method3(%C* %43)
	%45 = call %IO* @IO.out_int(%IO* %41, i32 %44)
	%46 = call %IO* @IO.out_string(%IO* %41, i8* getelementptr ([2 x i8], [2 x i8]* @.str17, i32 0, i32 0))
	%47 = load %C*, %C** %10
	%48 = call i32 @B.method3(%C* %47)
	%49 = call %IO* @IO.out_int(%IO* %41, i32 %48)
	%50 = call %IO* @IO.out_string(%Main* %self, i8* getelementptr ([2 x i8], [2 x i8]* @.str18, i32 0, i32 0))
	%51 = bitcast %IO* %50 to %Object*
	ret %Object* %51
}

define i32 @main() {
entry:
	%0 = call i8* @malloc(%Main* getelementptr (%Main, %Main* null, i32 1))
	%1 = bitcast i8* %0 to %Main*
	%2 = getelementptr %Main, %Main* %1, i32 0, i32 0
	%3 = bitcast [8 x i8*]* @vtable.Main to i8*
	store i8* %3, i8** %2
	%4 = alloca %Main*
	store %Main* %1, %Main** %4
	%5 = getelementptr [8 x i8*], [8 x i8*]* @vtable.Main, i32 0, i32 4
	%6 = load i8*, i8** %5
	%7 = bitcast i8* %6 to i8* (%Main*)*
	%8 = call i8* %7(%Main* %1)
	ret i32 0
}
