class Main inherits IO {
   -- Method to check if a number is prime
   isPrime(n : Int) : Bool {
      let i : Int <- 2,
          prime : Bool <- true in {
         if n <= 1 then prime <- false
         else if n = 2 then prime <- true
         else {
            while i * i <= n loop {
               if (n - (n/i)*i) = 0 then {
                  prime <- false;
                  i <- n;  -- Exit loop
               } else
                  i <- i + 1
               fi;
            } pool;
         }
         fi fi;
         prime;
      }
   };

   main() : Object { {
      out_string("Testing isPrime function\n");
      
      out_string("Is 2 prime? ");
      if isPrime(2) then out_string("Yes, it's prime\n")
      else out_string("No, it's not prime\n") fi;
      
      out_string("Is 7 prime? ");
      if isPrime(7) then out_string("Yes, it's prime\n")
      else out_string("No, it's not prime\n") fi;
      
      out_string("Is 10 prime? ");
      if isPrime(10) then out_string("Yes, it's prime\n")
      else out_string("No, it's not prime\n") fi;
      
      out_string("Is 17 prime? ");
      if isPrime(17) then out_string("Yes, it's prime\n")
      else out_string("No, it's not prime\n") fi;
   }
   };
};
