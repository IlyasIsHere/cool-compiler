class A {
   method1() : Int { 1 };
   method2() : Int { 2 };
   method3() : Int { 3 };
};

class B inherits A {
   method2() : Int { 22 }; 
   method3() : Int { 35 }; 
};

class C inherits B {
   method3() : Int { 33 }; 
};

class Main inherits IO {
   main() : Object {
      let a : A <- new A,
          b : B <- new B,
          c : C <- new C
           in
      {
         out_int(c@A.method1()).out_string(" ").out_int(c.method1());
         out_string("\n");
         out_int(c@B.method2()).out_string(" ").out_int(c.method2());
         out_string("\n");
         out_int(b.method2()).out_string(" ").out_int(b@A.method2());
         out_string("\n");
         out_int(c.method3()).out_string(" ").out_int(c@A.method3()).out_string(" ").out_int(c@B.method3());
         out_string("\n");
      }
   };
};
