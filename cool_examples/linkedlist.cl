-- test_linkedlist.cl - Demonstrates the built-in LinkedList

class Main inherits IO {
    main() : Object {
        let 
            list : LinkedList <- new LinkedList
        in {
            -- Initialize and add some integers
            list.init();
            out_string("Created a new LinkedList\n");
            
            -- Add elements to the list
            list.addFirst(100);
            list.addLast(200);
            list.addFirst(50);
            
            -- Print the list size
            out_string("List size: ");
            out_int(list.size());
            out_string("\n");
            
            -- Print first element
            out_string("First element: ");
            out_int(list.get(0));
            out_string("\n");
            
            -- Print second element
            out_string("Second element: ");
            out_int(list.get(1));
            out_string("\n");
            
            -- Print third element
            out_string("Third element: ");
            out_int(list.get(2));
            out_string("\n");
            
            -- Remove first element and print it
            out_string("Removed first element: ");
            out_int(list.removeFirst());
            out_string("\n");
            
            -- Print the new list size
            out_string("New list size: ");
            out_int(list.size());
            out_string("\n");
            
            -- Print the new first element
            out_string("New first element: ");
            out_int(list.get(0));
            out_string("\n");
            
            -- Check if list is empty
            if list.isEmpty() then
                out_string("List is empty\n")
            else
                out_string("List is not empty\n")
            fi;

            list.removeFirst();
            list.removeFirst();

            if list.isEmpty() then
                out_string("List is empty\n")
            else 
                out_string("List is not empty\n")
            fi;
        }
    };
}; 