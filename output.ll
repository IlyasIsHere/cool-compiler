%Object = type { i8* }
%IO = type { i8* }
%Int = type { i8* }
%String = type { i8* }
%Bool = type { i8* }
%Sum = type { i8*, i32, i32, i8* }
%Main = type { i8* }

@vtable.Object = global [0 x i8*] zeroinitializer
@vtable.IO = global [4 x i8*] [i8* bitcast (i32 (%IO*)* @IO.in_int to i8*), i8* bitcast (i8* (%IO*)* @IO.in_string to i8*), i8* bitcast (%IO* (%IO*, i32)* @IO.out_int to i8*), i8* bitcast (%IO* (%IO*, i8*)* @IO.out_string to i8*)]
@vtable.Int = global [0 x i8*] zeroinitializer
@vtable.String = global [0 x i8*] zeroinitializer
@vtable.Bool = global [0 x i8*] zeroinitializer
@vtable.Sum = global [5 x i8*] [i8* bitcast (i8* (%Sum*)* @Sum.getName to i8*), i8* bitcast (%Sum* (%Sum*, i32, i32)* @Sum.init to i8*), i8* bitcast (%IO* (%Sum*)* @Sum.printName to i8*), i8* bitcast (i8* (%Sum*, i8*)* @Sum.setName to i8*), i8* bitcast (i32 (%Sum*)* @Sum.sum to i8*)]
@vtable.Main = global [1 x i8*] [i8* bitcast (%Object* (%Main*)* @Main.main to i8*)]
@.str7 = internal constant [17 x i8] c"Hello from init\0A\00"
@.str8 = internal constant [2 x i8] c"\0A\00"
@.str9 = internal constant [2 x i8] c"\0A\00"
@.str10 = internal constant [2 x i8] c"\0A\00"
@.str11 = internal constant [6 x i8] c"Ilyas\00"
@.str12 = internal constant [13 x i8] c"The sum is: \00"
@.str13 = internal constant [2 x i8] c"\0A\00"
@.str14 = internal constant [13 x i8] c"another name\00"
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
	%0 = call i8* @malloc({ i8*, i8* }* getelementptr ({ i8*, i8* }, { i8*, i8* }* null, i32 1))
	%1 = bitcast i8* %0 to { i8*, i8* }*
	%2 = getelementptr { i8*, i8* }, { i8*, i8* }* %1, i32 0, i32 0
	%3 = bitcast [0 x i8*]* @vtable.String to i8*
	store i8* %3, i8** %2
	%4 = getelementptr { i8*, i8* }, { i8*, i8* }* %self, i32 0, i32 1
	%5 = load i8*, i8** %4
	%6 = getelementptr { i8*, i8* }, { i8*, i8* }* %other, i32 0, i32 1
	%7 = load i8*, i8** %6
	%8 = call i32 @strlen(i8* %5)
	%9 = call i32 @strlen(i8* %7)
	%10 = add i32 %8, %9
	%11 = add i32 %10, 1
	%12 = zext i32 %11 to i64
	%13 = call i8* @malloc(i64 %12)
	%14 = call i8* @strcpy(i8* %13, i8* %5)
	%15 = call i8* @strcat(i8* %13, i8* %7)
	%16 = getelementptr { i8*, i8* }, { i8*, i8* }* %1, i32 0, i32 1
	store i8* %13, i8** %16
	ret { i8*, i8* }* %1
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

define %Sum* @Sum.init(%Sum* %self, i32 %a, i32 %b) {
entry:
	%0 = call i8* @malloc(%IO* getelementptr (%IO, %IO* null, i32 1))
	%1 = bitcast i8* %0 to %IO*
	%2 = getelementptr %IO, %IO* %1, i32 0, i32 0
	%3 = bitcast [4 x i8*]* @vtable.IO to i8*
	store i8* %3, i8** %2
	%4 = call %IO* @IO.out_string(%IO* %1, i8* getelementptr ([17 x i8], [17 x i8]* @.str7, i32 0, i32 0))
	%5 = getelementptr %Sum, %Sum* %self, i32 0, i32 1
	store i32 %a, i32* %5
	%6 = getelementptr %Sum, %Sum* %self, i32 0, i32 2
	store i32 %b, i32* %6
	%7 = call i8* @malloc(%IO* getelementptr (%IO, %IO* null, i32 1))
	%8 = bitcast i8* %7 to %IO*
	%9 = getelementptr %IO, %IO* %8, i32 0, i32 0
	%10 = bitcast [4 x i8*]* @vtable.IO to i8*
	store i8* %10, i8** %9
	%11 = getelementptr %Sum, %Sum* %self, i32 0, i32 1
	%12 = load i32, i32* %11
	%13 = call %IO* @IO.out_int(%IO* %8, i32 %12)
	%14 = call %IO* @IO.out_string(%IO* %8, i8* getelementptr ([2 x i8], [2 x i8]* @.str8, i32 0, i32 0))
	%15 = getelementptr %Sum, %Sum* %self, i32 0, i32 2
	%16 = load i32, i32* %15
	%17 = call %IO* @IO.out_int(%IO* %8, i32 %16)
	%18 = call %IO* @IO.out_string(%IO* %8, i8* getelementptr ([2 x i8], [2 x i8]* @.str9, i32 0, i32 0))
	%19 = bitcast %Sum* %self to %Sum*
	ret %Sum* %19
}

define i32 @Sum.sum(%Sum* %self) {
entry:
	%0 = getelementptr %Sum, %Sum* %self, i32 0, i32 1
	%1 = load i32, i32* %0
	%2 = getelementptr %Sum, %Sum* %self, i32 0, i32 2
	%3 = load i32, i32* %2
	%4 = add i32 %1, %3
	ret i32 %4
}

define %IO* @Sum.printName(%Sum* %self) {
entry:
	%0 = call i8* @malloc(%IO* getelementptr (%IO, %IO* null, i32 1))
	%1 = bitcast i8* %0 to %IO*
	%2 = getelementptr %IO, %IO* %1, i32 0, i32 0
	%3 = bitcast [4 x i8*]* @vtable.IO to i8*
	store i8* %3, i8** %2
	%4 = getelementptr %Sum, %Sum* %self, i32 0, i32 3
	%5 = load i8*, i8** %4
	%6 = call %IO* @IO.out_string(%IO* %1, i8* %5)
	%7 = call %IO* @IO.out_string(%IO* %1, i8* getelementptr ([2 x i8], [2 x i8]* @.str10, i32 0, i32 0))
	%8 = bitcast %IO* %1 to %IO*
	ret %IO* %8
}

define i8* @Sum.setName(%Sum* %self, i8* %s) {
entry:
	%0 = getelementptr %Sum, %Sum* %self, i32 0, i32 3
	store i8* %s, i8** %0
	%1 = bitcast i8* %s to i8*
	ret i8* %1
}

define i8* @Sum.getName(%Sum* %self) {
entry:
	%0 = getelementptr %Sum, %Sum* %self, i32 0, i32 3
	%1 = load i8*, i8** %0
	%2 = bitcast i8* %1 to i8*
	ret i8* %2
}

define %Object* @Main.main(%Main* %self) {
entry:
	%0 = alloca %Sum*
	%1 = call i8* @malloc(%Sum* getelementptr (%Sum, %Sum* null, i32 1))
	%2 = bitcast i8* %1 to %Sum*
	%3 = getelementptr %Sum, %Sum* %2, i32 0, i32 0
	%4 = bitcast [5 x i8*]* @vtable.Sum to i8*
	store i8* %4, i8** %3
	%5 = getelementptr %Sum, %Sum* %2, i32 0, i32 3
	store i8* getelementptr ([6 x i8], [6 x i8]* @.str11, i32 0, i32 0), i8** %5
	%6 = getelementptr %Sum, %Sum* %2, i32 0, i32 0
	%7 = load i8*, i8** %6
	%8 = bitcast i8* %7 to [0 x i8*]*
	%9 = getelementptr [0 x i8*], [0 x i8*]* %8, i32 0, i32 1
	%10 = load i8*, i8** %9
	%11 = bitcast i8* %10 to %Sum* (%Sum*, i32, i32)*
	%12 = call %Sum* %11(%Sum* %2, i32 5, i32 20)
	store %Sum* %12, %Sum** %0
	%13 = alloca i32
	store i32 0, i32* %13
	%14 = alloca %IO*
	%15 = call i8* @malloc(%IO* getelementptr (%IO, %IO* null, i32 1))
	%16 = bitcast i8* %15 to %IO*
	%17 = getelementptr %IO, %IO* %16, i32 0, i32 0
	%18 = bitcast [4 x i8*]* @vtable.IO to i8*
	store i8* %18, i8** %17
	store %IO* %16, %IO** %14
	%19 = load %Sum*, %Sum** %0
	%20 = getelementptr %Sum, %Sum* %19, i32 0, i32 0
	%21 = load i8*, i8** %20
	%22 = bitcast i8* %21 to [0 x i8*]*
	%23 = getelementptr [0 x i8*], [0 x i8*]* %22, i32 0, i32 4
	%24 = load i8*, i8** %23
	%25 = bitcast i8* %24 to i32 (%Sum*)*
	%26 = call i32 %25(%Sum* %19)
	store i32 %26, i32* %13
	%27 = load %IO*, %IO** %14
	%28 = call %IO* @IO.out_string(%IO* %27, i8* getelementptr ([13 x i8], [13 x i8]* @.str12, i32 0, i32 0))
	%29 = load %IO*, %IO** %14
	%30 = load i32, i32* %13
	%31 = call %IO* @IO.out_int(%IO* %29, i32 %30)
	%32 = load %IO*, %IO** %14
	%33 = call %IO* @IO.out_string(%IO* %32, i8* getelementptr ([2 x i8], [2 x i8]* @.str13, i32 0, i32 0))
	%34 = load %Sum*, %Sum** %0
	%35 = getelementptr %Sum, %Sum* %34, i32 0, i32 0
	%36 = load i8*, i8** %35
	%37 = bitcast i8* %36 to [0 x i8*]*
	%38 = getelementptr [0 x i8*], [0 x i8*]* %37, i32 0, i32 2
	%39 = load i8*, i8** %38
	%40 = bitcast i8* %39 to %IO* (%Sum*)*
	%41 = call %IO* %40(%Sum* %34)
	%42 = load %Sum*, %Sum** %0
	%43 = getelementptr %Sum, %Sum* %42, i32 0, i32 0
	%44 = load i8*, i8** %43
	%45 = bitcast i8* %44 to [0 x i8*]*
	%46 = getelementptr [0 x i8*], [0 x i8*]* %45, i32 0, i32 3
	%47 = load i8*, i8** %46
	%48 = bitcast i8* %47 to i8* (%Sum*, i8*)*
	%49 = call i8* %48(%Sum* %42, i8* getelementptr ([13 x i8], [13 x i8]* @.str14, i32 0, i32 0))
	%50 = load %Sum*, %Sum** %0
	%51 = getelementptr %Sum, %Sum* %50, i32 0, i32 0
	%52 = load i8*, i8** %51
	%53 = bitcast i8* %52 to [0 x i8*]*
	%54 = getelementptr [0 x i8*], [0 x i8*]* %53, i32 0, i32 2
	%55 = load i8*, i8** %54
	%56 = bitcast i8* %55 to %IO* (%Sum*)*
	%57 = call %IO* %56(%Sum* %50)
	%58 = load %IO*, %IO** %14
	%59 = load %Sum*, %Sum** %0
	%60 = getelementptr %Sum, %Sum* %59, i32 0, i32 0
	%61 = load i8*, i8** %60
	%62 = bitcast i8* %61 to [0 x i8*]*
	%63 = getelementptr [0 x i8*], [0 x i8*]* %62, i32 0, i32 0
	%64 = load i8*, i8** %63
	%65 = bitcast i8* %64 to i8* (%Sum*)*
	%66 = call i8* %65(%Sum* %59)
	%67 = call %IO* @IO.out_string(%IO* %58, i8* %66)
	%68 = bitcast %IO* %58 to %Object*
	ret %Object* %68
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
	%4 = alloca %Main*
	store %Main* %1, %Main** %4
	%5 = getelementptr [1 x i8*], [1 x i8*]* @vtable.Main, i32 0, i32 0
	%6 = load i8*, i8** %5
	%7 = bitcast i8* %6 to i8* (%Main*)*
	%8 = call i8* %7(%Main* %1)
	ret i32 0
}
