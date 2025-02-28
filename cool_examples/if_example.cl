class Main {
   main() : Object {
      let io : IO <- new IO,
          num : Int <- 10
      in {
         io.out_string("Testing if condition:\n");
         
         if num < 20 then
            io.out_string("The number ").out_int(num).out_string(" is less than 20\n")
         else
            io.out_string("The number ").out_int(num).out_string(" is greater than or equal to 20\n")
         fi;
         
         if num = 10 then {
            io.out_string("The number is exactly 10!\n");
            io.out_string("Let's add 5 to it: ").out_int(num + 5).out_string("\n");
         } else
            io.out_string("The number is not 10\n")
         fi;
         
         io;
      }
   };
}; 