class Sum {
    num1 : Int;
    num2 : Int;
    name : String <- "Ilyas";

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
            num1 + num2
    };

    printName() : IO {
        {
            (new IO).out_string(name).out_string("\n");
        }
    };

    setName(s : String) : String {
        name <- s
    };

    getName() : String {
        name
    };

};

class Main {
   main() : Object {
      let 
         s : Sum <- (new Sum).init(5, 20),
         result : Int, 
         io : IO <- new IO
      in {
         result <- s.sum();
         io.out_string("The sum is: ");
         io.out_int(result);
         io.out_string("\n");
         
         s.printName();
         s.setName("another name");
         s.printName();

         io.out_string(s.getName());
      }
   };
};


