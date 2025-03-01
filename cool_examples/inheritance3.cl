class Shape inherits IO {
   identify() : Object { out_string("I am a shape\n") };
   area() : Int { 0 };
};

class Circle inherits Shape {
   identify() : Object { out_string("I am a circle\n") };
   area() : Int { 314 };  -- Simplified area calculation
};

class Rectangle inherits Shape {
   identify() : Object { out_string("I am a rectangle\n") };
   area() : Int { 200 };  -- Simplified area calculation
};

class Main inherits IO {
   main() : Object {
      let
         s : Shape <- new Shape,
         c : Shape <- new Circle,
         r : Shape <- new Rectangle
      in {
         out_string("Case Statement with Inheritance:\n");
         
         processShape(s);
         processShape(c);
         processShape(r);
      }
   };
   
   processShape(s : Shape) : Object {
      {
         case s of
            c : Circle => {
               c.identify();
               out_int(c.area()).out_string("\n");
            };
            r : Rectangle => {
               r.identify();
               out_int(r.area()).out_string("\n");
            };
            s : Shape => {
               s.identify();
               out_int(s.area()).out_string("\n");
            };
         esac;
      }
   };
}; 