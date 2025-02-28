class Sum {
    num1 : Int;
    num2 : Int;

    init(a : Int, b : Int) : Sum {
        {
            (new IO).out_string("Hello from init\n");
            num1 <- a;
            num2 <- b;
            (new IO).out_int(num1).out_string("\n").out_int(num2).out_string("\n");
            self;
        }
    };

    sum() : Int {
        {
            (new IO).out_string("Hello from sum\n");
            num1 + num2;
        }
    };

    another() : String {
        {
            (new IO).out_string("hhhhhhhh");
            "Ilyas";
        }
    };
};

class Main {
   main() : Object {
      let 
         s : Sum <- (new Sum).init(5, 20),
         result : Int, 
         test_str : String <- new String,
         io : IO <- new IO,
         obj : Object <- new Object
      in {
         result <- s.sum();
         io.out_string("The sum is: ");
         io.out_int(result);
         io.out_string("\n");
         test_str <- "test STRING has changed haha";
         io.out_string(s.another()).out_string("\n");
         io.out_string(test_str);
      }
   };
};


