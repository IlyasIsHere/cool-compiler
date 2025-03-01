class GCD {
};

class Main {
    io : IO <- new IO;
    
    main() : Object {
        {
            io.out_string("Greatest Common Divisor (GCD) Demo\n\n");
            
            -- Test some GCD examples
            io.out_string("GCD of 48 and 18: ");
            io.out_int(gcd_recursive(48, 18));
            io.out_string("\n");
            
            io.out_string("GCD of 101 and 103: ");
            io.out_int(gcd_recursive(101, 103));
            io.out_string("\n");
            
            io.out_string("GCD of 1071 and 462: ");
            io.out_int(gcd_recursive(1071, 462));
            io.out_string("\n");
            
            -- Show a step-by-step GCD calculation
            io.out_string("\nStep-by-step GCD calculation for 1071 and 462:\n");
            gcd_recursive_steps(1071, 462, 1);
            
            -- Compare with iterative implementation
            io.out_string("\nComparing with iterative implementation:\n");
            io.out_string("GCD of 1071 and 462 (iterative): ");
            io.out_int(gcd_iterative(1071, 462));
            io.out_string("\n");
        }
    };

    -- Recursive implementation of Euclidean GCD algorithm
    gcd_recursive(a : Int, b : Int) : Int {
        if b = 0 then 
            a  -- Base case: GCD(a, 0) = a
        else 
            gcd_recursive(b, a - (a / b) * b)  -- Recursive case: GCD(a, b) = GCD(b, a mod b)
        fi
    };
    
    -- Step-by-step recursive GCD with level tracking
    gcd_recursive_steps(a : Int, b : Int, level : Int) : Int {
        let indent_str : String <- self.get_indent(level) in {
            io.out_string(indent_str);
            io.out_string("Step ");
            io.out_int(level);
            io.out_string(": gcd(");
            io.out_int(a);
            io.out_string(", ");
            io.out_int(b);
            io.out_string(")\n");
            
            if b = 0 then {
                io.out_string(indent_str);
                io.out_string("Result = ");
                io.out_int(a);
                io.out_string(" (Base case reached)\n");
                a;  -- Base case: GCD(a, 0) = a
            } else {
                let remainder : Int <- a - (a / b) * b in {
                    io.out_string(indent_str);
                    io.out_string("  ");
                    io.out_int(a);
                    io.out_string(" mod ");
                    io.out_int(b);
                    io.out_string(" = ");
                    io.out_int(remainder);
                    io.out_string("\n");
                    
                    -- Recursive call
                    gcd_recursive_steps(b, remainder, level + 1);
                };
            }
            fi;
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
    
    -- Iterative implementation of Euclidean GCD algorithm
    gcd_iterative(a : Int, b : Int) : Int {
        let a_val : Int <- a,
            b_val : Int <- b,
            temp : Int in {
            while not (b_val = 0) loop {
                temp <- b_val;
                b_val <- a_val - (a_val / b_val) * b_val;  -- b = a mod b
                a_val <- temp;
            } pool;
            a_val;
        }
    };
}; 