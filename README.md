# SNAKE IN CONSOLE
Based on a well known game that mock a snake moving in the screen that getting longer each time it eats food.

This is just a practice to understand more about concurency. EXPECT BUGS!!

## Concept
- Uses a linked list to record the every snake's chunk position.
- Uses a map to index the chunks for faster search.
- The chunk has next node that goes to the head, and prev node that goes to the tail
- Movement is done in 2 mode:
    - Growing-head, where new chunk is added to the next head, replacing the old head. 
    - Shrink-tail, where the old tail is removed as the count reaches the maximum count. Happens when the chunks count reaches chunk's maximum count.
- When the head's next move hits the other chunks in the map, the game ends
- When the head hits food it increase the chunk maximum count

## Caveats
- The basic packages provide by Go are lacking of tools to handle the keyboard and screen interactions, therefore additinal libraries are used
- A minimalist libraries from atomicgo are being used, see https://github.com/atomicgo
- The cursor package that is being used has no control over the coordinate, the calculation is done manually
- There are still bugs regarding routine with timer, where there is some part where in the routine that waits the after timer to be executed. I used a tricky solution for now by using counter.