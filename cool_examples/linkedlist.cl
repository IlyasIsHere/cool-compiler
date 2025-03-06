-- LinkedList.cl - A linked list implementation for Cool language
-- Provides a standard linked list data structure with common operations

class Node {
    -- Node in a linked list containing data and a reference to the next node
    data : Object;       -- The data stored in this node
    next : Node;         -- Reference to the next node (void if end of list)
    
    -- Initialize a new node with data and optional next node
    init(value : Object, nextNode : Node) : SELF_TYPE {
        {
            data <- value;
            next <- nextNode;
            self;
        }
    };
    
    -- Get the data stored in this node
    getData() : Object { data };
    
    -- Get the next node
    getNext() : Node { next };
    
    -- Set the next node reference
    setNext(nextNode : Node) : SELF_TYPE {
        {
            next <- nextNode;
            self;
        }
    };
};

class LinkedList {
    head : Node;         -- First node in the list
    tail : Node;         -- Last node in the list
    count : Int <- 0;    -- Number of elements in the list
    io : IO <- new IO;   -- IO object for printing
    
    -- Initialize an empty linked list
    init() : SELF_TYPE {
        {
            -- Head and tail are initialized as void implicitly
            count <- 0;
            self;
        }
    };
    
    -- Add an element to the beginning of the list
    addFirst(value : Object) : SELF_TYPE {
        let newNode : Node <- new Node in {
            newNode.init(value, head);
            
            -- If list was empty, also update tail
            if isEmpty() then
                tail <- newNode
            else 
                false
            fi;
            
            head <- newNode;
            count <- count + 1;
            self;
        }
    };
    
    -- Add an element to the end of the list
    addLast(value : Object) : SELF_TYPE {
        let newNode : Node <- new Node in {
            -- Initialize with next as void
            newNode.init(value, newNode);
            -- Clear the circular reference to make next void
            let temp : Node in newNode.setNext(temp);
            
            -- If list was empty, also update head
            if isEmpty() then
                head <- newNode
            else
                tail.setNext(newNode)
            fi;
            
            tail <- newNode;
            count <- count + 1;
            self;
        }
    };
    
    -- Remove the first element from the list
    removeFirst() : Object {
        if isEmpty() then {
            io.out_string("Error: Cannot remove from an empty list\n");
            abort();
            new Object;
        } else {
            let value : Object <- head.getData() in {
                head <- head.getNext();
                count <- count - 1;
                
                -- If list is now empty, also clear tail reference
                if isEmpty() then {
                    let temp : Node in tail <- temp; -- Make tail void
                } else 
                    false
                fi;
                
                value;
            };
        }
        fi
    };
    
    -- Remove the last element from the list
    removeLast() : Object {
        if isEmpty() then {
            io.out_string("Error: Cannot remove from an empty list\n");
            abort();
            new Object;
        } else if count = 1 then {
            -- Special case: only one element in the list
            let 
                value : Object <- head.getData(),
                temp : Node
            in {
                head <- temp; -- Make head void
                tail <- temp; -- Make tail void
                count <- 0;
                value;
            };
        } else {
            -- Find the second-to-last node
            let 
                current : Node <- head,
                index : Int <- 0,
                value : Object,
                temp : Node
            in {
                -- Navigate to second-to-last node
                while index < count - 2 loop {
                    current <- current.getNext();
                    index <- index + 1;
                }
                pool;
                
                -- Get value from last node before removing it
                value <- current.getNext().getData();
                
                -- Update tail and remove reference to last node
                tail <- current;
                current.setNext(temp); -- Make next reference void
                
                count <- count - 1;
                value;
            };
        }
        fi fi
    };
    
    -- Get element at the specified index
    get(index : Int) : Object {
        if index < 0 then {
            io.out_string("Error: Index out of bounds (negative index)\n");
            abort();
            new Object;
        } else if count <= index then {
            io.out_string("Error: Index out of bounds (index too large)\n");
            abort();
            new Object;
        } else {
            let 
                current : Node <- head,
                currentIndex : Int <- 0
            in {
                while currentIndex < index loop {
                    current <- current.getNext();
                    currentIndex <- currentIndex + 1;
                }
                pool;
                
                current.getData();
            };
        }
        fi fi
    };
    
    -- Get the number of elements in the list
    size() : Int { count };
    
    -- Check if the list is empty
    isEmpty() : Bool { isvoid head };
};

-- Simple example usage
class Main {
    main() : Object {
        let 
            list : LinkedList <- new LinkedList,
            io : IO <- new IO
        in {
            -- Initialize and add elements
            list.init();
            list.addFirst(100);
            list.addLast(200);
            list.addFirst(50);
            
            -- Print size and elements
            io.out_string("List size: ");
            io.out_int(list.size());
            io.out_string("\n");
            
            -- Print first element
            io.out_string("First element: ");
            io.out_int(case list.get(0) of i : Int => i; esac);
            io.out_string("\n");
            
            -- Remove first element and print it
            io.out_string("Removed: ");
            io.out_int(case list.removeFirst() of i : Int => i; esac);
            io.out_string("\n");
            
            -- Final size
            io.out_string("New size: ");
            io.out_int(list.size());
            io.out_string("\n");
        }
    };
}; 