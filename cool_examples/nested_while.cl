class Main {
    i : Int; -- default value is 0
    j : Int;
    io : IO <- new IO;
    
    main() : Object {
        {
            io.out_string("Nested while loops:\n");
            
            while i < 3 loop {
                io.out_string("Outer: ");
                io.out_int(i);
                io.out_string("\n");
                
                j <- 0;
                while j < 2 loop {
                    io.out_string("  Inner: ");
                    io.out_int(j);
                    io.out_string("\n");
                    j <- j + 1;
                }
                pool;
                
                i <- i + 1;
            }
            pool;
            
            io.out_string("Done.\n");
        }
    };
};

