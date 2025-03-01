class Main {
   main() : Object {
      let io : IO <- new IO,
          message : String,
          number : Int
      in {
         io.out_string("Enter a number: ");
         number <- io.in_int();
         io.out_string("You entered the number: ").out_int(number).out_string("\n");
         io.out_string("Double of your number: ").out_int(number * 2);
         
         io;
      }
   };
}; 