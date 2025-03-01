-- Multi-level inheritance example
class A {
    a : Int <- 10;
    
    getA() : Int { a };
    
    identify() : String { "Class A" };
};

class B inherits A {
    b : Int <- 20;
    
    getB() : Int { b };
    
    identify() : String { "Class B" };
};

class C inherits B {
    c : Int <- 30;
    
    getC() : Int { c };
    
    identify() : String { "Class C" };
};

class Main inherits IO {
    main() : Object {
        let 
            objA : A <- new A,
            objB : B <- new B,
            objC : C <- new C
        in {
            out_string(objA.identify()).out_string(": ").out_int(objA.getA()).out_string("\n");
            out_int(objA.getA()).out_string("\n");
            out_string(objB.identify()).out_string(": ").out_int(objB.getA() + objB.getB()).out_string("\n");
            out_int(objB.getA()).out_string("\n");
            out_int(objB.getB()).out_string("\n");
            out_string(objC.identify()).out_string(": ").out_int(objC.getA() + objC.getB() + objC.getC()).out_string("\n");
            out_int(objC.getA()).out_string("\n");
            out_int(objC.getB()).out_string("\n");
            out_int(objC.getC()).out_string("\n");
        }
    };
};
