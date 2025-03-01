class Factorial {
};

class Main {
    io : IO <- new IO;
    
    main() : Object {
        {
           let num_recursive: Int,
               num_iterative: Int in {
                num_recursive <- factorial(5);
                io.out_string("Recursive factorial: ");
                io.out_int(num_recursive);
                io.out_string("\n");
                
                num_iterative <- self.factorial_while(5);
                io.out_string("Iterative factorial: ");
                io.out_int(num_iterative);
                io.out_string("\n");
           }; 
        }
    };

    factorial(n : Int) : Int {
        if n <= 1 then 1 else n * self.factorial(n-1) fi
    };
    
    factorial_while(n : Int) : Int {
        let result : Int <- 1,
            i : Int <- 1 in {
            while i <= n loop {
                result <- result * i;
                i <- i + 1;
            } pool;
            result;
        }
    };
};
