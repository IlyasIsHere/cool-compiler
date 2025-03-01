class Main {
    counter : Int <- 0;
    limit : Int <- 5;
    io : IO <- new IO;
    
    main() : Object {
        {
            io.out_string("Starting while loop demonstration\n");
            
            while counter < limit loop {
                io.out_string("Counter value: ");
                io.out_int(counter);
                io.out_string("\n");
                
                if counter = 2 then {
                    io.out_string("Found the middle value!\n");
                    io.out_string("Executing nested while loop...\n");
                    
                    let nestedCounter : Int <- 0 in {
                        while nestedCounter < 3 loop {
                            io.out_string("  Nested counter: ");
                            io.out_int(nestedCounter);
                            io.out_string("\n");
                            nestedCounter <- nestedCounter + 1;
                        } pool;
                        
                        io.out_string("Nested while loop completed\n");
                    };
                } else {
                    io.out_string("Regular iteration\n");
                }
                fi;
                
                counter <- counter + 1;
            } pool;
            
            io.out_string("While loop completed. Final counter value: ");
            io.out_int(counter);
            io.out_string("\n");
        }
    };
};
