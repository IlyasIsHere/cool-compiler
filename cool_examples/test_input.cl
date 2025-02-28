class Main {
   main() : Object {
      let io : IO <- new IO,
          message : String,
          number : Int
      in {
         io.out_string("Enter a message: ");
         message <- io.in_string();
         io.out_string("You entered: ").out_string(message).out_string("\n");
         
         io.out_string("Enter a number: ");
         number <- io.in_int();
         io.out_string("You entered the number: ").out_int(number).out_string("\n");
         
         io;
      }
   };
}; 