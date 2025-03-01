class Main inherits IO {
   main() : Object {
      let x : Int <- 10 in
      {
         out_string("Testing case with Int (10):\n");
         case x of
            s : String => out_string("It's a string: ").out_string(s).out_string("\n");
            i : Int => out_string("It's an integer: ").out_int(i).out_string("\n");
            b : Bool => out_string("It's a boolean\n");
            o : Object => out_string("It's some other object\n");
         esac;
      }
   };
}; 