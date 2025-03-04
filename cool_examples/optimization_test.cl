(*
 * Simple Optimization Test for Constant Propagation and Function Inlining
 * 
 * This program is deliberately simple but computationally intensive
 * to clearly demonstrate the impact of optimizations.
 *)

class Main inherits IO {
    -- Constants that can be optimized by constant propagation
    const1 : Int <- 10;
    const2 : Int <- 20;
    const3 : Int <- 5;
    
    -- Computed constant that could be determined at compile-time
    computed_const : Int <- const1 * const2; -- 200
    
    -- Simple function that is a good candidate for inlining
    add(a : Int, b : Int) : Int {
        a + b
    };
    
    -- Function with constant values that should be propagated
    constCalc(a : Int) : Int {
        a * const1 + computed_const
    };
    
    -- Function that combines multiple inlinable calls
    nestedCalc(a : Int) : Int {
        add(constCalc(a), add(a, const3))
    };
    
    -- Compute an expensive operation without using tail recursion
    -- This is deliberately inefficient as a baseline
    fib(n : Int) : Int {
        if n <= 1 then 
            1
        else 
            fib(n-1) + fib(n-2)
        fi
    };
    
    main() : Object {
        let i : Int <- 0,
            j : Int <- 0,
            result : Int <- 0,
            iterations : Int <- 10000000  -- Very high to make it run long
        in {
            out_string("Starting optimization test...\n");
            
            -- Warm-up with something slow regardless of optimization
            out_string("Warming up with Fibonacci(20)...\n");
            result <- fib(20);
            out_string("Fibonacci result: ").out_int(result).out_string("\n\n");
            
            -- Now the main loop that can benefit from optimizations
            out_string("Starting main computation loop (").out_int(iterations).out_string(" iterations)...\n");
            
            -- Reset for main test
            result <- 0;
            i <- 0;
            
            -- Simple loop structure that's easier to understand
            while i < iterations loop {
                -- These operations will be much faster with inlining and constant propagation
                result <- nestedCalc(i);
                
                -- Show progress periodically
                if i - (i/1000000)*1000000 = 0 then {
                    out_string("Progress: ").out_int(i / 100000).out_string("%\n");
                } else 0 fi;
                
                i <- i + 1;
            } pool;
            
            out_string("\nFinal result: ").out_int(result).out_string("\n");
            out_string("Test complete! If optimizations are working, this should be much faster.\n");
        }
    };
}; 