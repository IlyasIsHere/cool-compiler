%Object = type { i8* }
%IO = type { i8* }
%Int = type { i8* }
%String = type { i8* }
%Bool = type { i8* }
%Main = type { i8* }

@vtable.Object = global [0 x i8*] zeroinitializer
@vtable.IO = global [1 x i8*] [i8* bitcast (%IO* (%IO*, i8*)* @IO.out_string to i8*)]
@vtable.Int = global [0 x i8*] zeroinitializer
@vtable.String = global [0 x i8*] zeroinitializer
@vtable.Bool = global [0 x i8*] zeroinitializer
@vtable.Main = global [1 x i8*] [i8* bitcast (%Object* (%Main*)* @Main.main to i8*)]
@.str6 = internal constant [20 x i8] c"Hello, COOL World!\0A\00"
@.str.Object = constant [7 x i8] c"Object\00"
@.str.fmt = global [3 x i8] c"%s\00"
@.str.fmt.int = constant [3 x i8] c"%d\00"

define %IO* @IO.out_string(%IO* %self, i8* %str) {
entry:
	%0 = call i32 (i8*, ...) @printf(i8* getelementptr ([3 x i8], [3 x i8]* @.str.fmt, i32 0, i32 0), i8* %str)
	ret %IO* %self
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
	%8 = bitcast i8* %5 to i8*
	%9 = bitcast i8* %7 to i8*
	%10 = call i8* @string_concat(i8* %8, i8* %9)
	%11 = getelementptr { i8*, i8* }, { i8*, i8* }* %1, i32 0, i32 1
	store i8* %10, i8** %11
	ret { i8*, i8* }* %1
}

define { i8*, i8* }* @String.substr({ i8*, i8* }* %self, i32 %start, i32 %length) {
entry:
	%0 = getelementptr { i8*, i8* }, { i8*, i8* }* %self, i32 0, i32 1
	%1 = load i8*, i8** %0
	%2 = call i8* @string_substr(i8* %1, i32 %start, i32 %length)
	%3 = call i8* @malloc({ i8*, i8* }* getelementptr ({ i8*, i8* }, { i8*, i8* }* null, i32 1))
	%4 = bitcast i8* %3 to { i8*, i8* }*
	%5 = getelementptr { i8*, i8* }, { i8*, i8* }* %4, i32 0, i32 0
	%6 = bitcast [0 x i8*]* @vtable.String to i8*
	store i8* %6, i8** %5
	%7 = getelementptr { i8*, i8* }, { i8*, i8* }* %4, i32 0, i32 1
	store i8* %2, i8** %7
	ret { i8*, i8* }* %4
}

declare i8* @malloc(i64 %size)

declare void @free(i8* %ptr)

declare i32 @printf(i8* %format, ...)

declare i8* @in_string()

declare i32 @in_int()

declare i32 @string_length(i8* %str)

declare i8* @string_concat(i8* %str1, i8* %str2)

declare i8* @string_substr(i8* %str, i32 %start, i32 %length)

declare void @abort()

declare i8* @type_name(i8* %obj)

declare i8* @object_copy(i8* %obj)

declare void @case_abort()

declare void @dispatch_abort()

define %Object* @Main.main(%Main* %self) {
entry:
	%0 = alloca i8*
	store i8* getelementptr ([20 x i8], [20 x i8]* @.str6, i32 0, i32 0), i8** %0
	%1 = call i8* @malloc(%IO* getelementptr (%IO, %IO* null, i32 1))
	%2 = bitcast i8* %1 to %IO*
	%3 = getelementptr %IO, %IO* %2, i32 0, i32 0
	%4 = bitcast [1 x i8*]* @vtable.IO to i8*
	store i8* %4, i8** %3
	%5 = load i8*, i8** %0
	%6 = call %IO* @IO.out_string(%IO* %2, i8* %5)
	%7 = bitcast %IO* %2 to %Object*
	ret %Object* %7
}

define %Object* @Object.abort(%Object* %self) {
entry:
	call void @abort()
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

define %IO* @IO.out_int(%IO* %self, i32 %x) {
entry:
	%0 = call i32 (i8*, ...) @printf(i8* getelementptr ([3 x i8], [3 x i8]* @.str.fmt.int, i32 0, i32 0), i32 %x)
	ret %IO* %self
}

define i8* @IO.in_string(%IO* %self) {
entry:
	%0 = call i8* @in_string()
	ret i8* %0
}

define i32 @IO.in_int(%IO* %self) {
entry:
	%0 = call i32 @in_int()
	ret i32 %0
}

define i32 @String.length({ i8*, i8* }* %self) {
entry:
	%0 = getelementptr { i8*, i8* }, { i8*, i8* }* %self, i32 0, i32 1
	%1 = load i8*, i8** %0
	%2 = call i32 @string_length(i8* %1)
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
