%Object = type { i8* }
%IO = type { i8* }
%Int = type { i8* }
%String = type { i8* }
%Bool = type { i8* }
%Main = type { i8*, i32, %IO* }

@vtable.Object = global [0 x i8*] zeroinitializer
@vtable.IO = global [4 x i8*] [i8* bitcast (i32 (%IO*)* @IO.in_int to i8*), i8* bitcast (i8* (%IO*)* @IO.in_string to i8*), i8* bitcast (%IO* (%IO*, i32)* @IO.out_int to i8*), i8* bitcast (%IO* (%IO*, i8*)* @IO.out_string to i8*)]
@vtable.Int = global [0 x i8*] zeroinitializer
@vtable.String = global [0 x i8*] zeroinitializer
@vtable.Bool = global [0 x i8*] zeroinitializer
@vtable.Main = global [1 x i8*] [i8* bitcast (%Object* (%Main*)* @Main.main to i8*)]
@.str6 = internal constant [20 x i8] c"Outer if: num < 20\0A\00"
@.str7 = internal constant [12 x i8] c"The number \00"
@.str8 = internal constant [18 x i8] c" is less than 10\0A\00"
@.str9 = internal constant [12 x i8] c"The number \00"
@.str10 = internal constant [23 x i8] c" is between 10 and 19\0A\00"
@.str11 = internal constant [21 x i8] c"Outer if: num >= 20\0A\00"
@.str12 = internal constant [12 x i8] c"The number \00"
@.str13 = internal constant [23 x i8] c" is between 20 and 29\0A\00"
@.str14 = internal constant [12 x i8] c"The number \00"
@.str15 = internal constant [19 x i8] c" is 30 or greater\0A\00"
@.str16 = internal constant [28 x i8] c"---------------------------\00"
@.str17 = internal constant [25 x i8] c"\0ANow checking equality:\0A\00"
@.str18 = internal constant [27 x i8] c"The number is exactly 15!\0A\00"
@.str19 = internal constant [22 x i8] c"The number is not 15\0A\00"
@.str20 = internal constant [27 x i8] c"The number is exactly 20!\0A\00"
@.str21 = internal constant [22 x i8] c"The number is not 20\0A\00"
@.str22 = internal constant [40 x i8] c"This is after the nested if statements\0A\00"
@.str.abort_msg = constant [17 x i8] c"Program aborted\0A\00"
@.str.Object = constant [7 x i8] c"Object\00"
@.str.fmt = global [3 x i8] c"%s\00"
@.str.fmt.int = constant [3 x i8] c"%d\00"
@.str.scanf_s_fmt = constant [3 x i8] c"%s\00"
@.str.scanf_d_fmt = constant [3 x i8] c"%d\00"
@.str.substr_error = constant [39 x i8] c"Runtime error: substring out of range\0A\00"

define %IO* @IO.out_string(%IO* %self, i8* %str) {
entry:
	%0 = call i32 (i8*, ...) @printf(i8* getelementptr ([3 x i8], [3 x i8]* @.str.fmt, i32 0, i32 0), i8* %str)
	ret %IO* %self
}

define %IO* @IO.out_int(%IO* %self, i32 %i) {
entry:
	%0 = call i32 (i8*, ...) @printf(i8* getelementptr ([3 x i8], [3 x i8]* @.str.fmt.int, i32 0, i32 0), i32 %i)
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

define { i8*, i8* }* @String.concat({ i8*, i8* }* %self, { i8*, i8* }* %other) {
entry:
	%0 = getelementptr { i8*, i8* }, { i8*, i8* }* %self, i32 0, i32 1
	%1 = load i8*, i8** %0
	%2 = getelementptr { i8*, i8* }, { i8*, i8* }* %other, i32 0, i32 1
	%3 = load i8*, i8** %2
	%4 = call i32 @strlen(i8* %1)
	%5 = call i32 @strlen(i8* %3)
	%6 = add i32 %4, %5
	%7 = add i32 %6, 1
	%8 = zext i32 %7 to i64
	%9 = call i8* @malloc(i64 %8)
	%10 = call i8* @strcpy(i8* %9, i8* %1)
	%11 = call i8* @strcat(i8* %9, i8* %3)
	%12 = getelementptr { i8*, i8* }, { i8*, i8* }* %0, i32 0, i32 1
	store i8* %9, i8** %12
	ret { i8*, i8* }* %0
}

define { i8*, i8* }* @String.substr({ i8*, i8* }* %self, i32 %start, i32 %length) {
entry:
	%0 = getelementptr { i8*, i8* }, { i8*, i8* }* %self, i32 0, i32 1
	%1 = load i8*, i8** %0
	%2 = call i32 @strlen(i8* %1)
	br label %bounds_check

bounds_check:
	%3 = icmp slt i32 %start, 0
	%4 = icmp sge i32 %start, %2
	%5 = or i1 %3, %4
	%6 = icmp slt i32 %length, 0
	%7 = or i1 %5, %6
	br i1 %7, label %error, label %alloc

alloc:
	%8 = add i32 %length, 1
	%9 = zext i32 %8 to i64
	%10 = call i8* @malloc(i64 %9)
	%11 = getelementptr i8, i8* %1, i32 %start
	%12 = call i8* @strncpy(i8* %10, i8* %11, i32 %length)
	%13 = getelementptr i8, i8* %10, i32 %length
	store i8 0, i8* %13
	%14 = call i8* @malloc({ i8*, i8* }* getelementptr ({ i8*, i8* }, { i8*, i8* }* null, i32 1))
	%15 = bitcast i8* %14 to { i8*, i8* }*
	%16 = getelementptr { i8*, i8* }, { i8*, i8* }* %15, i32 0, i32 0
	%17 = bitcast [0 x i8*]* @vtable.String to i8*
	store i8* %17, i8** %16
	%18 = getelementptr { i8*, i8* }, { i8*, i8* }* %15, i32 0, i32 1
	store i8* %10, i8** %18
	ret { i8*, i8* }* %15

error:
	%19 = bitcast [39 x i8]* @.str.substr_error to i8*
	%20 = call i32 (i8*, ...) @printf(i8* %19)
	call void @exit(i32 1)
	unreachable
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
	%3 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%4 = load %IO*, %IO** %3
	%5 = call %IO* @IO.out_string(%IO* %4, i8* getelementptr ([20 x i8], [20 x i8]* @.str6, i32 0, i32 0))
	%6 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%7 = load i32, i32* %6
	%8 = icmp slt i32 %7, 10
	br i1 %8, label %if.then.2, label %if.else.2

if.else.1:
	%9 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%10 = load %IO*, %IO** %9
	%11 = call %IO* @IO.out_string(%IO* %10, i8* getelementptr ([21 x i8], [21 x i8]* @.str11, i32 0, i32 0))
	%12 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%13 = load i32, i32* %12
	%14 = icmp slt i32 %13, 30
	br i1 %14, label %if.then.3, label %if.else.3

if.end.1:
	%15 = phi %IO* [ %47, %if.then.1 ], [ %70, %if.else.1 ]
	%16 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%17 = load %IO*, %IO** %16
	%18 = call %IO* @IO.out_string(%IO* %17, i8* getelementptr ([28 x i8], [28 x i8]* @.str16, i32 0, i32 0))
	%19 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%20 = load %IO*, %IO** %19
	%21 = call %IO* @IO.out_string(%IO* %20, i8* getelementptr ([25 x i8], [25 x i8]* @.str17, i32 0, i32 0))
	%22 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%23 = load i32, i32* %22
	%24 = icmp eq i32 %23, 15
	br i1 %24, label %if.then.4, label %if.else.4

if.then.2:
	%25 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%26 = load %IO*, %IO** %25
	%27 = call %IO* @IO.out_string(%IO* %26, i8* getelementptr ([12 x i8], [12 x i8]* @.str7, i32 0, i32 0))
	%28 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%29 = load %IO*, %IO** %28
	%30 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%31 = load i32, i32* %30
	%32 = call %IO* @IO.out_int(%IO* %29, i32 %31)
	%33 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%34 = load %IO*, %IO** %33
	%35 = call %IO* @IO.out_string(%IO* %34, i8* getelementptr ([18 x i8], [18 x i8]* @.str8, i32 0, i32 0))
	br label %if.end.2

if.else.2:
	%36 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%37 = load %IO*, %IO** %36
	%38 = call %IO* @IO.out_string(%IO* %37, i8* getelementptr ([12 x i8], [12 x i8]* @.str9, i32 0, i32 0))
	%39 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%40 = load %IO*, %IO** %39
	%41 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%42 = load i32, i32* %41
	%43 = call %IO* @IO.out_int(%IO* %40, i32 %42)
	%44 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%45 = load %IO*, %IO** %44
	%46 = call %IO* @IO.out_string(%IO* %45, i8* getelementptr ([23 x i8], [23 x i8]* @.str10, i32 0, i32 0))
	br label %if.end.2

if.end.2:
	%47 = phi %IO* [ %34, %if.then.2 ], [ %45, %if.else.2 ]
	br label %if.end.1

if.then.3:
	%48 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%49 = load %IO*, %IO** %48
	%50 = call %IO* @IO.out_string(%IO* %49, i8* getelementptr ([12 x i8], [12 x i8]* @.str12, i32 0, i32 0))
	%51 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%52 = load %IO*, %IO** %51
	%53 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%54 = load i32, i32* %53
	%55 = call %IO* @IO.out_int(%IO* %52, i32 %54)
	%56 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%57 = load %IO*, %IO** %56
	%58 = call %IO* @IO.out_string(%IO* %57, i8* getelementptr ([23 x i8], [23 x i8]* @.str13, i32 0, i32 0))
	br label %if.end.3

if.else.3:
	%59 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%60 = load %IO*, %IO** %59
	%61 = call %IO* @IO.out_string(%IO* %60, i8* getelementptr ([12 x i8], [12 x i8]* @.str14, i32 0, i32 0))
	%62 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%63 = load %IO*, %IO** %62
	%64 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%65 = load i32, i32* %64
	%66 = call %IO* @IO.out_int(%IO* %63, i32 %65)
	%67 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%68 = load %IO*, %IO** %67
	%69 = call %IO* @IO.out_string(%IO* %68, i8* getelementptr ([19 x i8], [19 x i8]* @.str15, i32 0, i32 0))
	br label %if.end.3

if.end.3:
	%70 = phi %IO* [ %57, %if.then.3 ], [ %68, %if.else.3 ]
	br label %if.end.1

if.then.4:
	%71 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%72 = load %IO*, %IO** %71
	%73 = call %IO* @IO.out_string(%IO* %72, i8* getelementptr ([27 x i8], [27 x i8]* @.str18, i32 0, i32 0))
	br label %if.end.4

if.else.4:
	%74 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%75 = load %IO*, %IO** %74
	%76 = call %IO* @IO.out_string(%IO* %75, i8* getelementptr ([22 x i8], [22 x i8]* @.str19, i32 0, i32 0))
	br label %if.end.4

if.end.4:
	%77 = phi %IO* [ %72, %if.then.4 ], [ %75, %if.else.4 ]
	%78 = getelementptr %Main, %Main* %self, i32 0, i32 1
	%79 = load i32, i32* %78
	%80 = icmp eq i32 %79, 20
	br i1 %80, label %if.then.5, label %if.else.5

if.then.5:
	%81 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%82 = load %IO*, %IO** %81
	%83 = call %IO* @IO.out_string(%IO* %82, i8* getelementptr ([27 x i8], [27 x i8]* @.str20, i32 0, i32 0))
	br label %if.end.5

if.else.5:
	%84 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%85 = load %IO*, %IO** %84
	%86 = call %IO* @IO.out_string(%IO* %85, i8* getelementptr ([22 x i8], [22 x i8]* @.str21, i32 0, i32 0))
	br label %if.end.5

if.end.5:
	%87 = phi %IO* [ %82, %if.then.5 ], [ %85, %if.else.5 ]
	%88 = getelementptr %Main, %Main* %self, i32 0, i32 2
	%89 = load %IO*, %IO** %88
	%90 = call %IO* @IO.out_string(%IO* %89, i8* getelementptr ([40 x i8], [40 x i8]* @.str22, i32 0, i32 0))
	%91 = bitcast %IO* %89 to %Object*
	%92 = call i8* @malloc({ i8*, i8* }* getelementptr ({ i8*, i8* }, { i8*, i8* }* null, i32 1))
	%93 = bitcast i8* %92 to { i8*, i8* }*
	%94 = getelementptr { i8*, i8* }, { i8*, i8* }* %93, i32 0, i32 0
	%95 = bitcast [0 x i8*]* @vtable.String to i8*
	store i8* %95, i8** %94
	ret %Object* %91
}

define %Object* @Object.abort(%Object* %self) {
entry:
	%0 = call i32 (i8*, ...) @printf(i8* getelementptr ([17 x i8], [17 x i8]* @.str.abort_msg, i32 0, i32 0))
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

define i32 @String.length({ i8*, i8* }* %self) {
entry:
	%0 = getelementptr { i8*, i8* }, { i8*, i8* }* %self, i32 0, i32 1
	%1 = load i8*, i8** %0
	%2 = call i32 @strlen(i8* %1)
	ret i32 %2
}

define i32 @main() {
entry:
	%0 = call i8* @malloc(%Main* getelementptr (%Main, %Main* null, i32 1))
	%1 = bitcast i8* %0 to %Main*
	%2 = getelementptr %Main, %Main* %1, i32 0, i32 0
	%3 = bitcast [1 x i8*]* @vtable.Main to i8*
	store i8* %3, i8** %2
	%4 = getelementptr %Main, %Main* %1, i32 0, i32 1
	store i32 15, i32* %4
	%5 = call i8* @malloc(%IO* getelementptr (%IO, %IO* null, i32 1))
	%6 = bitcast i8* %5 to %IO*
	%7 = getelementptr %IO, %IO* %6, i32 0, i32 0
	%8 = bitcast [4 x i8*]* @vtable.IO to i8*
	store i8* %8, i8** %7
	%9 = getelementptr %Main, %Main* %1, i32 0, i32 2
	store %IO* %6, %IO** %9
	%10 = alloca %Main*
	store %Main* %1, %Main** %10
	%11 = getelementptr [1 x i8*], [1 x i8*]* @vtable.Main, i32 0, i32 0
	%12 = load i8*, i8** %11
	%13 = bitcast i8* %12 to i8* (%Main*)*
	%14 = call i8* %13(%Main* %1)
	ret i32 0
}
