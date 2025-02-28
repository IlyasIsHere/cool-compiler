class Main inherits IO {
   main() : IO {
      {
         let io : IO, s : Int
         in {
            io.out_string("Hello, COOL Ilyas!\n");
            io.out_int(6 / 3).out_string("\n");

            -- s <- io.in_int();
            io.out_int(s);
         };


         --(new IO).in_int();
         --(new IO).out_string("finished");
      }
   };

   (*
   myfunc() : IO {
      (new IO).out_string("ana mn myfunc hhhh\n")
   };
   *)
}; 