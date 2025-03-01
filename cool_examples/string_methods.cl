class Main inherits IO {
    main() : Object {
        {
            -- Demonstrate String.length() method
            let s1 : String <- "Hello, Cool!" in {
                out_string("String: ").out_string(s1).out_string("\n");
                out_string("Length: ").out_int(s1.length()).out_string("\n\n");
            };
            
            -- Demonstrate String.concat() method
            let s2 : String <- "Hello", 
                s3 : String <- "World!" in {
                out_string("String 1: ").out_string(s2).out_string("\n");
                out_string("String 2: ").out_string(s3).out_string("\n");
                out_string("Concatenated: ").out_string(s2.concat(s3)).out_string("\n\n");
            };
            
            -- Demonstrate String.substr() method
            let s4 : String <- "The Cool Programming Language" in {
                out_string("String: ").out_string(s4).out_string("\n");
                out_string("Substring (4,4): ").out_string(s4.substr(4,4)).out_string("\n");
                out_string("Substring (0,3): ").out_string(s4.substr(0,3)).out_string("\n");
                out_string("Substring (9,11): ").out_string(s4.substr(9,7)).out_string("\n\n");
            };
        }
    };
};
