class Animal inherits IO {
    name : String <- "Animal";
    
    speak() : Object {
        out_string(name).out_string(" makes a sound\n")
    };
    
    eat() : Object {
        out_string(name).out_string(" eats food\n")
    };
    
    setName(n : String) : SELF_TYPE {
        {
            name <- n;
            self;
        }
    };
};

class Dog inherits Animal {
    breed : String;
    
    init(b : String) : SELF_TYPE {
        {
            breed <- b;
            name <- "Dog";
            self;
        }
    };
    
    -- Override parent method
    speak() : Object {
        out_string(name).out_string(" barks!\n")
    };
    
    specialEat() : Object {
        {
            out_string(breed).out_string(" ");
            eat();  -- Calls parent's eat method
        }
    };
};

class Cat inherits Animal {
    -- Override parent method
    speak() : Object {
        out_string(name).out_string(" meows!\n")
    };
    
    -- New method specific to Cat
    purr() : Object {
        out_string(name).out_string(" purrs softly\n")
    };
};

class Main {
    io : IO <- new IO;
    
    main() : Object {
        let
            animal : Animal <- new Animal,
            dog : Dog <- (new Dog).init("Labrador").setName("Rex"),
            cat : Cat <- (new Cat).setName("Whiskers")
        in {
            animal.speak();
            dog.speak();
            cat.speak();
            
            dog.specialEat();
            cat.purr();
        }
    };
};
