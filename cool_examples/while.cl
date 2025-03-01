class Main {
    counter : Int <- 0;
    limit : Int <- 5;
    io : IO <- new IO;
    
    main() : Object {
        {
            io.out_string("Starting while loop demonstration\n");
            
            while counter < limit loop
            {
                io.out_string("Counter value: ");
                io.out_int(counter);
                io.out_string("\n");
                
                counter <- counter + 1;
            }
            pool;
            
            io.out_string("While loop completed. Final counter value: ");
            io.out_int(counter);
            io.out_string("\n");
        }
    };
};
