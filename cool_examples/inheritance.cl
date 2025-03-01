-- Inheritance example in COOL demonstrating:
-- 1. Method overriding
-- 2. Attribute inheritance
-- 3. Dynamic dispatch
-- 4. SELF_TYPE usage

class Shape {
    name : String <- "Shape";
    color : String <- "white";
    
    identify() : Object {
        (new IO).out_string("I am a ").out_string(name)
            .out_string(" with color ").out_string(color)
            .out_string("\n")
    };
    
    getArea() : Int { 0 };
    
    setColor(c : String) : SELF_TYPE {
        {
            color <- c;
            self;
        }
    };
};

class Circle inherits Shape {
    radius : Int <- 5;
    
    init(r : Int) : SELF_TYPE {
        {
            radius <- r;
            name <- "Circle";
            self;
        }
    };
    
    -- Override parent method
    getArea() : Int { 3 * radius * radius };
    
    -- Child-specific method
    getCircumference() : Int { 2 * 3 * radius };
};

class Rectangle inherits Shape {
    width : Int <- 10;
    height : Int <- 5;
    
    init(w : Int, h : Int) : SELF_TYPE {
        {
            width <- w;
            height <- h;
            name <- "Rectangle";
            self;
        }
    };
    
    -- Override parent method
    getArea() : Int { width * height };
    
    -- Child-specific method
    getPerimeter() : Int { 2 * (width + height) };
};

class Main {
    io : IO <- new IO;
    
    main() : Object {
        let
            shape : Shape <- new Shape,
            circle : Circle <- (new Circle).init(7).setColor("red"),
            rectangle : Rectangle <- (new Rectangle).init(4, 8).setColor("blue")
        in {
            io.out_string("===== Shape Inheritance Demo =====\n");
            
            -- Parent class
            io.out_string("Base Shape: ");
            shape.identify();
            io.out_string("Area: ").out_int(shape.getArea()).out_string("\n\n");
            
            -- Circle subclass
            io.out_string("Circle: ");
            circle.identify();
            io.out_string("Area: ").out_int(circle.getArea()).out_string("\n");
            io.out_string("Circumference: ").out_int(circle.getCircumference()).out_string("\n\n");
            
            -- Rectangle subclass
            io.out_string("Rectangle: ");
            rectangle.identify();
            io.out_string("Area: ").out_int(rectangle.getArea()).out_string("\n");
            io.out_string("Perimeter: ").out_int(rectangle.getPerimeter()).out_string("\n\n");
            
        }
    };
};
