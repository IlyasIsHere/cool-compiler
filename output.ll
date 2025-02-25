%Object = type { i8* }
%IO = type { i8* }
%Int = type { i8* }
%String = type { i8* }
%Bool = type { i8* }
%Main = type { i8* }

@vtable.Object = global [3 x i8*] zeroinitializer
@vtable.IO = global [4 x i8*] zeroinitializer
@vtable.Int = global [0 x i8*] zeroinitializer
@vtable.String = global [3 x i8*] zeroinitializer
@vtable.Bool = global [0 x i8*] zeroinitializer
@vtable.Main = global [1 x i8*] zeroinitializer
@.str6 = internal constant [20 x i8] c"Hello, COOL World!\0A\00"

declare i8* @malloc(i64 %size)

declare void @free(i8* %ptr)

declare i32 @out_string(i8* %str)

declare i32 @out_int(i32 %num)

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

declare %Object* @Object.abort(%Object* %self)

declare i8* @Object.type_name(%Object* %self)

declare %Object* @Object.copy(%Object* %self)

declare %IO* @IO.out_string(%IO* %self, i8* %x)

declare %IO* @IO.out_int(%IO* %self, i32 %x)

declare i8* @IO.in_string(%IO* %self)

declare i32 @IO.in_int(%IO* %self)

declare i32 @String.length(%String* %self)

declare i8* @String.concat(%String* %self, i8* %s)

declare i8* @String.substr(%String* %self, i32 %i, i32 %l)

define %Object* @Main.main(%Main* %self) {
entry:
	%0 = alloca i8*
	store i8* getelementptr ([20 x i8], [20 x i8]* @.str6, i32 0, i32 0), i8** %0
	%1 = call i8* @malloc(i64 ptrtoint (%IO* getelementptr (%IO, %IO* null, i32 1) to i64))
	%2 = bitcast i8* %1 to %IO*
	%3 = getelementptr %IO, %IO* %2, i32 0, i32 0
	%4 = bitcast [4 x i8*]* @vtable.IO to i8*
	store i8* %4, i8** %3
	%5 = load i8*, i8** %0
	%6 = getelementptr %IO, %IO* %2, i32 0, i32 0
	%7 = load i8*, i8** %6
	%8 = bitcast i8* %7 to [0 x i8*]*
	%9 = getelementptr [0 x i8*], [0 x i8*]* %8, i32 0, i32 0
	%10 = load i8*, i8** %9
	%11 = bitcast i8* %10 to void ()*
	call void %11(%IO* %2, i8* %5)
	ret void %0
}

define i32 @main() {
entry:
	%0 = call i8* @malloc(i64 ptrtoint (%Main* getelementptr (%Main, %Main* null, i32 1) to i64))
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
