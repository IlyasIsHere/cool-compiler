%Object = type { i8* }
%IO = type { i8* }
%Int = type { i8* }
%String = type { i8* }
%Bool = type { i8* }
%Main = type { i8* }

@.str.empty = constant [1 x i8] c"\00"
@vtable.Object = global [3 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.IO = global [7 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%IO*)* @IO.in_int to i8*), i8* bitcast (i8* (%IO*)* @IO.in_string to i8*), i8* bitcast (%IO* (%IO*, i32)* @IO.out_int to i8*), i8* bitcast (%IO* (%IO*, i8*)* @IO.out_string to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Int = global [3 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.String = global [6 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (i8* (%String*, i8*)* @String.concat to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%String*)* @String.length to i8*), i8* bitcast (i8* (%String*, i32, i32)* @String.substr to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Bool = global [3 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@vtable.Main = global [9 x i8*] [i8* bitcast (%Object* (%Object*)* @Object.abort to i8*), i8* bitcast (%Object* (%Object*)* @Object.copy to i8*), i8* bitcast (i32 (%IO*)* @IO.in_int to i8*), i8* bitcast (i8* (%IO*)* @IO.in_string to i8*), i8* bitcast (i1 (%Main*, i32)* @Main.isPrime to i8*), i8* bitcast (%Object* (%Main*)* @Main.main to i8*), i8* bitcast (%IO* (%IO*, i32)* @IO.out_int to i8*), i8* bitcast (%IO* (%IO*, i8*)* @IO.out_string to i8*), i8* bitcast (i8* (%Object*)* @Object.type_name to i8*)]
@.str7 = internal constant [26 x i8] c"Testing isPrime function\0A\00"
@.str8 = internal constant [13 x i8] c"Is 2 prime? \00"
@.str9 = internal constant [17 x i8] c"Yes, it's prime\0A\00"
@.str10 = internal constant [20 x i8] c"No, it's not prime\0A\00"
@.str11 = internal constant [13 x i8] c"Is 7 prime? \00"
@.str12 = internal constant [17 x i8] c"Yes, it's prime\0A\00"
@.str13 = internal constant [20 x i8] c"No, it's not prime\0A\00"
@.str14 = internal constant [14 x i8] c"Is 10 prime? \00"
@.str15 = internal constant [17 x i8] c"Yes, it's prime\0A\00"
@.str16 = internal constant [20 x i8] c"No, it's not prime\0A\00"
@.str17 = internal constant [14 x i8] c"Is 17 prime? \00"
@.str18 = internal constant [17 x i8] c"Yes, it's prime\0A\00"
@.str19 = internal constant [20 x i8] c"No, it's not prime\0A\00"
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

define i1 @Main.isPrime(%Main* %self, i32 %n) {
entry:
	%0 = alloca i32
	store i32 2, i32* %0
	%1 = alloca i1
	store i1 true, i1* %1
	%2 = icmp sle i32 %n, 1
	br i1 %2, label %if.then.1, label %if.else.1

if.then.1:
	store i1 false, i1* %1
	%3 = inttoptr i1 false to i8*
	br label %if.end.1

if.else.1:
	%4 = icmp eq i32 %n, 2
	br i1 %4, label %if.then.2, label %if.else.2

if.end.1:
	%5 = phi i8* [ %3, %if.then.1 ], [ %8, %if.end.2 ]
	%6 = load i1, i1* %1
	ret i1 %6

if.then.2:
	store i1 true, i1* %1
	%7 = inttoptr i1 true to i8*
	br label %if.end.2

if.else.2:
	br label %while.cond.1

if.end.2:
	%8 = phi i8* [ %7, %if.then.2 ], [ null, %while.exit.1 ]
	br label %if.end.1

while.cond.1:
	%9 = load i32, i32* %0
	%10 = load i32, i32* %0
	%11 = mul i32 %9, %10
	%12 = icmp sle i32 %11, %n
	br i1 %12, label %while.body.1, label %while.exit.1

while.body.1:
	%13 = load i32, i32* %0
	%14 = sdiv i32 %n, %13
	%15 = load i32, i32* %0
	%16 = mul i32 %14, %15
	%17 = sub i32 %n, %16
	%18 = icmp eq i32 %17, 0
	br i1 %18, label %if.then.3, label %if.else.3

while.exit.1:
	br label %if.end.2

if.then.3:
	store i1 false, i1* %1
	store i32 %n, i32* %0
	br label %if.end.3

if.else.3:
	%19 = load i32, i32* %0
	%20 = add i32 %19, 1
	store i32 %20, i32* %0
	br label %if.end.3

if.end.3:
	%21 = phi i32 [ %n, %if.then.3 ], [ %20, %if.else.3 ]
	br label %while.cond.1
}

define %Object* @Main.main(%Main* %self) {
entry:
	%0 = getelementptr [9 x i8*], [9 x i8*]* @vtable.Main, i32 0, i32 7
	%1 = load i8*, i8** %0
	%2 = bitcast i8* %1 to %Object* (%IO*, i8*)*
	%3 = bitcast %Main* %self to %IO*
	%4 = call %Object* %2(%IO* %3, i8* getelementptr ([26 x i8], [26 x i8]* @.str7, i32 0, i32 0))
	%5 = getelementptr [9 x i8*], [9 x i8*]* @vtable.Main, i32 0, i32 7
	%6 = load i8*, i8** %5
	%7 = bitcast i8* %6 to %Object* (%IO*, i8*)*
	%8 = bitcast %Main* %self to %IO*
	%9 = call %Object* %7(%IO* %8, i8* getelementptr ([13 x i8], [13 x i8]* @.str8, i32 0, i32 0))
	%10 = getelementptr [9 x i8*], [9 x i8*]* @vtable.Main, i32 0, i32 4
	%11 = load i8*, i8** %10
	%12 = bitcast i8* %11 to i1 (%Main*, i32)*
	%13 = bitcast %Main* %self to %Main*
	%14 = call i1 %12(%Main* %13, i32 2)
	br i1 %14, label %if.then.4, label %if.else.4

if.then.4:
	%15 = getelementptr [9 x i8*], [9 x i8*]* @vtable.Main, i32 0, i32 7
	%16 = load i8*, i8** %15
	%17 = bitcast i8* %16 to %Object* (%IO*, i8*)*
	%18 = bitcast %Main* %self to %IO*
	%19 = call %Object* %17(%IO* %18, i8* getelementptr ([17 x i8], [17 x i8]* @.str9, i32 0, i32 0))
	br label %if.end.4

if.else.4:
	%20 = getelementptr [9 x i8*], [9 x i8*]* @vtable.Main, i32 0, i32 7
	%21 = load i8*, i8** %20
	%22 = bitcast i8* %21 to %Object* (%IO*, i8*)*
	%23 = bitcast %Main* %self to %IO*
	%24 = call %Object* %22(%IO* %23, i8* getelementptr ([20 x i8], [20 x i8]* @.str10, i32 0, i32 0))
	br label %if.end.4

if.end.4:
	%25 = phi %Object* [ %19, %if.then.4 ], [ %24, %if.else.4 ]
	%26 = getelementptr [9 x i8*], [9 x i8*]* @vtable.Main, i32 0, i32 7
	%27 = load i8*, i8** %26
	%28 = bitcast i8* %27 to %Object* (%IO*, i8*)*
	%29 = bitcast %Main* %self to %IO*
	%30 = call %Object* %28(%IO* %29, i8* getelementptr ([13 x i8], [13 x i8]* @.str11, i32 0, i32 0))
	%31 = getelementptr [9 x i8*], [9 x i8*]* @vtable.Main, i32 0, i32 4
	%32 = load i8*, i8** %31
	%33 = bitcast i8* %32 to i1 (%Main*, i32)*
	%34 = bitcast %Main* %self to %Main*
	%35 = call i1 %33(%Main* %34, i32 7)
	br i1 %35, label %if.then.5, label %if.else.5

if.then.5:
	%36 = getelementptr [9 x i8*], [9 x i8*]* @vtable.Main, i32 0, i32 7
	%37 = load i8*, i8** %36
	%38 = bitcast i8* %37 to %Object* (%IO*, i8*)*
	%39 = bitcast %Main* %self to %IO*
	%40 = call %Object* %38(%IO* %39, i8* getelementptr ([17 x i8], [17 x i8]* @.str12, i32 0, i32 0))
	br label %if.end.5

if.else.5:
	%41 = getelementptr [9 x i8*], [9 x i8*]* @vtable.Main, i32 0, i32 7
	%42 = load i8*, i8** %41
	%43 = bitcast i8* %42 to %Object* (%IO*, i8*)*
	%44 = bitcast %Main* %self to %IO*
	%45 = call %Object* %43(%IO* %44, i8* getelementptr ([20 x i8], [20 x i8]* @.str13, i32 0, i32 0))
	br label %if.end.5

if.end.5:
	%46 = phi %Object* [ %40, %if.then.5 ], [ %45, %if.else.5 ]
	%47 = getelementptr [9 x i8*], [9 x i8*]* @vtable.Main, i32 0, i32 7
	%48 = load i8*, i8** %47
	%49 = bitcast i8* %48 to %Object* (%IO*, i8*)*
	%50 = bitcast %Main* %self to %IO*
	%51 = call %Object* %49(%IO* %50, i8* getelementptr ([14 x i8], [14 x i8]* @.str14, i32 0, i32 0))
	%52 = getelementptr [9 x i8*], [9 x i8*]* @vtable.Main, i32 0, i32 4
	%53 = load i8*, i8** %52
	%54 = bitcast i8* %53 to i1 (%Main*, i32)*
	%55 = bitcast %Main* %self to %Main*
	%56 = call i1 %54(%Main* %55, i32 10)
	br i1 %56, label %if.then.6, label %if.else.6

if.then.6:
	%57 = getelementptr [9 x i8*], [9 x i8*]* @vtable.Main, i32 0, i32 7
	%58 = load i8*, i8** %57
	%59 = bitcast i8* %58 to %Object* (%IO*, i8*)*
	%60 = bitcast %Main* %self to %IO*
	%61 = call %Object* %59(%IO* %60, i8* getelementptr ([17 x i8], [17 x i8]* @.str15, i32 0, i32 0))
	br label %if.end.6

if.else.6:
	%62 = getelementptr [9 x i8*], [9 x i8*]* @vtable.Main, i32 0, i32 7
	%63 = load i8*, i8** %62
	%64 = bitcast i8* %63 to %Object* (%IO*, i8*)*
	%65 = bitcast %Main* %self to %IO*
	%66 = call %Object* %64(%IO* %65, i8* getelementptr ([20 x i8], [20 x i8]* @.str16, i32 0, i32 0))
	br label %if.end.6

if.end.6:
	%67 = phi %Object* [ %61, %if.then.6 ], [ %66, %if.else.6 ]
	%68 = getelementptr [9 x i8*], [9 x i8*]* @vtable.Main, i32 0, i32 7
	%69 = load i8*, i8** %68
	%70 = bitcast i8* %69 to %Object* (%IO*, i8*)*
	%71 = bitcast %Main* %self to %IO*
	%72 = call %Object* %70(%IO* %71, i8* getelementptr ([14 x i8], [14 x i8]* @.str17, i32 0, i32 0))
	%73 = getelementptr [9 x i8*], [9 x i8*]* @vtable.Main, i32 0, i32 4
	%74 = load i8*, i8** %73
	%75 = bitcast i8* %74 to i1 (%Main*, i32)*
	%76 = bitcast %Main* %self to %Main*
	%77 = call i1 %75(%Main* %76, i32 17)
	br i1 %77, label %if.then.7, label %if.else.7

if.then.7:
	%78 = getelementptr [9 x i8*], [9 x i8*]* @vtable.Main, i32 0, i32 7
	%79 = load i8*, i8** %78
	%80 = bitcast i8* %79 to %Object* (%IO*, i8*)*
	%81 = bitcast %Main* %self to %IO*
	%82 = call %Object* %80(%IO* %81, i8* getelementptr ([17 x i8], [17 x i8]* @.str18, i32 0, i32 0))
	br label %if.end.7

if.else.7:
	%83 = getelementptr [9 x i8*], [9 x i8*]* @vtable.Main, i32 0, i32 7
	%84 = load i8*, i8** %83
	%85 = bitcast i8* %84 to %Object* (%IO*, i8*)*
	%86 = bitcast %Main* %self to %IO*
	%87 = call %Object* %85(%IO* %86, i8* getelementptr ([20 x i8], [20 x i8]* @.str19, i32 0, i32 0))
	br label %if.end.7

if.end.7:
	%88 = phi %Object* [ %82, %if.then.7 ], [ %87, %if.else.7 ]
	%89 = bitcast %Object* %88 to %Object*
	ret %Object* %89
}

define i32 @main() {
entry:
	%0 = getelementptr %Main, %Main* null, i32 1
	%1 = ptrtoint %Main* %0 to i64
	%2 = call i8* @malloc(i64 %1)
	%3 = bitcast i8* %2 to %Main*
	%4 = getelementptr %Main, %Main* %3, i32 0, i32 0
	%5 = bitcast [9 x i8*]* @vtable.Main to i8*
	store i8* %5, i8** %4
	%6 = alloca %Main*
	store %Main* %3, %Main** %6
	%7 = getelementptr [9 x i8*], [9 x i8*]* @vtable.Main, i32 0, i32 5
	%8 = load i8*, i8** %7
	%9 = bitcast i8* %8 to i8* (%Main*)*
	%10 = call i8* %9(%Main* %3)
	ret i32 0
}
