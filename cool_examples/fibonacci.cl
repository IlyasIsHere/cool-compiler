class Fibonacci {
};

class Main {
    io : IO <- new IO;
    
    main() : Object {
        {
            io.out_string("Fibonacci Sequence Demo\n");
            
            -- Test recursive Fibonacci implementation
            io.out_string("Recursive Fibonacci(10): ");
            io.out_int(fibonacci_recursive(10));
            io.out_string("\n");
            
            -- Test iterative Fibonacci implementation
            io.out_string("Iterative Fibonacci(10): ");
            io.out_int(fibonacci_iterative(10));
            io.out_string("\n");
            
            -- Print the first 10 Fibonacci numbers using recursion
            io.out_string("\nFirst 10 Fibonacci numbers (recursive):\n");
            print_fibonacci_sequence(10);
            
            -- Compare performance with a larger number (not too large to avoid stack overflow)
            io.out_string("\nRecursive Fibonacci(20): ");
            io.out_int(fibonacci_recursive(20));
            io.out_string("\n");
            
            io.out_string("Iterative Fibonacci(20): ");
            io.out_int(fibonacci_iterative(20));
            io.out_string("\n");
        }
    };

    -- Recursive implementation of Fibonacci
    fibonacci_recursive(n : Int) : Int {
        if n <= 1 then 
            n 
        else 
            fibonacci_recursive(n-1) + fibonacci_recursive(n-2)
        fi
    };
    
    -- Iterative implementation of Fibonacci
    fibonacci_iterative(n : Int) : Int {
        let a : Int <- 0,
            b : Int <- 1,
            c : Int <- 0,
            i : Int <- 2 in {
            
            if n <= 0 then 
                0
            else if n = 1 then
                1
            else {
                while i <= n loop {
                    c <- a + b;
                    a <- b;
                    b <- c;
                    i <- i + 1;
                } pool;
                b;
            } fi fi;
        }
    };
    
    -- Function to print the first n Fibonacci numbers
    print_fibonacci_sequence(n : Int) : Object {
        let i : Int <- 0 in {
            while i < n loop {
                io.out_int(fibonacci_recursive(i));
                if i < n - 1 then 
                    io.out_string(", ")
                else 
                    io.out_string("\n")
                fi;
                i <- i + 1;
            } pool;
        }
    };
};
