class MutualRecursion {
};

class Main {
    io : IO <- new IO;
    
    main() : Object {
        {
            io.out_string("Mutual Recursion Demo (Even/Odd functions)\n\n");
            
            -- Test even numbers
            io.out_string("Testing even numbers:\n");
            test_number(0);
            test_number(2);
            test_number(4);
            test_number(6);
            test_number(10);
            test_number(20);
            
            -- Test odd numbers
            io.out_string("\nTesting odd numbers:\n");
            test_number(1);
            test_number(3);
            test_number(5);
            test_number(7);
            test_number(11);
            test_number(21);
            
            -- Test a larger number
            io.out_string("\nTesting a larger number:\n");
            test_number(42);
            test_number(99);
            
            -- Test recursive trace for a number
            io.out_string("\nRecursive trace for number 5:\n");
            io.out_string("is_even(5):\n");
            is_even_trace(5, 1);
        }
    };
    
    -- Helper method to test a number
    test_number(n : Int) : Object {
        {
            io.out_int(n);
            io.out_string(" is ");
            if is_even(n) then
                io.out_string("even")
            else
                io.out_string("odd")
            fi;
            io.out_string("\n");
        }
    };

    -- Recursive implementation for determining if a number is even
    -- Even numbers: 0, 2, 4, 6, etc.
    is_even(n : Int) : Bool {
        if n = 0 then
            true   -- Base case: 0 is even
        else if n = 1 then
            false  -- Base case: 1 is odd
        else
            is_odd(n-1)  -- If n-1 is odd, then n is even
        fi fi
    };
    
    -- Recursive implementation for determining if a number is odd
    -- Odd numbers: 1, 3, 5, 7, etc.
    is_odd(n : Int) : Bool {
        if n = 0 then
            false  -- Base case: 0 is not odd
        else if n = 1 then
            true   -- Base case: 1 is odd
        else
            is_even(n-1)  -- If n-1 is even, then n is odd
        fi fi
    };
    
    -- Traced version of is_even for visualization
    is_even_trace(n : Int, level : Int) : Bool {
        let indent_str : String <- self.get_indent(level) in {
            io.out_string(indent_str);
            io.out_string("is_even(");
            io.out_int(n);
            io.out_string(")\n");
            
            if n = 0 then {
                io.out_string(indent_str);
                io.out_string("  return true (Base case: 0 is even)\n");
                true;
            } else if n = 1 then {
                io.out_string(indent_str);
                io.out_string("  return false (Base case: 1 is odd)\n");
                false;
            } else {
                io.out_string(indent_str);
                io.out_string("  return is_odd(");
                io.out_int(n-1);
                io.out_string(")\n");
                let result : Bool <- is_odd_trace(n-1, level+1) in {
                    io.out_string(indent_str);
                    io.out_string("  is_even(");
                    io.out_int(n);
                    io.out_string(") returns ");
                    if result then
                        io.out_string("true\n")
                    else
                        io.out_string("false\n")
                    fi;
                    result;
                };
            } fi fi;
        }
    };
    
    -- Traced version of is_odd for visualization
    is_odd_trace(n : Int, level : Int) : Bool {
        let indent_str : String <- self.get_indent(level) in {
            io.out_string(indent_str);
            io.out_string("is_odd(");
            io.out_int(n);
            io.out_string(")\n");
            
            if n = 0 then {
                io.out_string(indent_str);
                io.out_string("  return false (Base case: 0 is not odd)\n");
                false;
            } else if n = 1 then {
                io.out_string(indent_str);
                io.out_string("  return true (Base case: 1 is odd)\n");
                true;
            } else {
                io.out_string(indent_str);
                io.out_string("  return is_even(");
                io.out_int(n-1);
                io.out_string(")\n");
                let result : Bool <- is_even_trace(n-1, level+1) in {
                    io.out_string(indent_str);
                    io.out_string("  is_odd(");
                    io.out_int(n);
                    io.out_string(") returns ");
                    if result then
                        io.out_string("true\n")
                    else
                        io.out_string("false\n")
                    fi;
                    result;
                };
            } fi fi;
        }
    };
    
    -- Helper method to generate indentation based on recursion level
    get_indent(level : Int) : String {
        let result : String <- "",
            i : Int <- 0 in {
            while i < level loop {
                result <- result.concat("  ");
                i <- i + 1;
            } pool;
            result;
        }
    };
}; 