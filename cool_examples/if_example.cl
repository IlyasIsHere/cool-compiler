class Main {
    num : Int <- 10;
    main() : Object {
      {
        if num < 20 then 
            (new IO).out_string("The number ").out_int(num).out_string(" is less than 20\n")
        else
            (new IO).out_string("The number ").out_int(num).out_string(" is greater than or equal to 20\n")
        fi;

        if num = 10 then
            (new IO).out_string("The number is exactly 10!\n")
        else
            (new IO).out_string("The number is not 10\n")
        fi;

        (new IO).out_string("This is after the if hahaha");
      }
    };
};