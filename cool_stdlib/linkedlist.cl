-- LinkedList.cl - A linked list implementation for Cool language
-- Provides an integer-only linked list data structure with common operations

class Node {
    -- Node in a linked list containing an integer value and a reference to the next node
    data : Int;       -- The integer data stored in this node
    next : Node;       -- Reference to the next node (void if end of list)
    
    -- Initialize a new node with data and optional next node
    init(value : Int, nextNode : Node) : SELF_TYPE {
        {
            data <- value;
            next <- nextNode;
            self;
        }
    };
    
    -- Get the data stored in this node
    getData() : Int { data };
    
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
    head : Node; 
    tail : Node; 
    count : Int <- 0; 
    
    -- Initialize an empty linked list
    init() : SELF_TYPE {
        {
            -- Head and tail are initialized as void implicitly
            count <- 0;
            self;
        }
    };
    
    -- Add an integer to the beginning of the list
    addFirst(value : Int) : SELF_TYPE {
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
    
    -- Add an integer to the end of the list
    addLast(value : Int) : SELF_TYPE {
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
    removeFirst() : Int {
        if isEmpty() then {
            (new IO).out_string("Error: Cannot remove from an empty list\n");
            abort();
            0; -- Dummy return value, never reached due to abort
        } else {
            let value : Int <- head.getData() in {
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
    
    -- Get element at the specified index
    get(index : Int) : Int {
        if index < 0 then {
            (new IO).out_string("Error: Index out of bounds (negative index)\n");
            abort();
            0; -- Dummy return value, never reached due to abort
        } else if count <= index then {
            (new IO).out_string("Error: Index out of bounds (index too large)\n");
            abort();
            0; -- Dummy return value, never reached due to abort
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