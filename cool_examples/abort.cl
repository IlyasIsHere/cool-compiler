class Main inherits IO {
    main() : Object {
        {
            out_string("Testing Object.abort()\n");
            
            -- Call abort method which will terminate the program
            out_string("About to call abort...\n");
            abort();
            
            -- This line should never be executed
            out_string("This line should not be printed\n");
        }
    };
};
