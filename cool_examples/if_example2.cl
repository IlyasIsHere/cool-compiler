class Main {
    num : Int <- 15;
    io : IO <- new IO;
    main() : Object {
      {
        if num < 20 then {
            io.out_string("Outer if: num < 20\n");
            if num < 10 then {
                io.out_string("The number ");
                io.out_int(num);
                io.out_string(" is less than 10\n");
            } else {
                io.out_string("The number ");
                io.out_int(num);
                io.out_string(" is between 10 and 19\n");
            }
            fi;
        }
        else {
            io.out_string("Outer if: num >= 20\n");
            if num < 30 then {
                io.out_string("The number ");
                io.out_int(num);
                io.out_string(" is between 20 and 29\n");
            } else {
                io.out_string("The number ");
                io.out_int(num);
                io.out_string(" is 30 or greater\n");
            }
            fi;
        }
        fi;

        io.out_string("---------------------------");
        io.out_string("\nNow checking equality:\n");
        
        if num = 15 then {
            io.out_string("The number is exactly 15!\n");
        } else {
            io.out_string("The number is not 15\n");
        }
        fi;
        
        if num = 20 then {
            io.out_string("The number is exactly 20!\n");
        } else {
            io.out_string("The number is not 20\n");
        }
        fi;

        io.out_string("This is after the nested if statements\n");
      }
    };
};
