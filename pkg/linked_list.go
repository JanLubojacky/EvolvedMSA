// linked lists used to store the sequences for the MSA
package pkg

import "fmt"

type Node struct {
	Value rune
	Prev  *Node
	Next  *Node
}

type Sequence struct {
	Head *Node
  IterPointer *Node // used to iterate over the sequence
  Length    int       // length of the sequence
}

// insert at the head of the list
func (list *Sequence) InsertHead(char rune) {
	newNode := &Node{Value: char, Next: list.Head}
	if list.Head != nil {
		list.Head.Prev = newNode
	}
	list.Head = newNode
  list.IterPointer = list.Head

  list.Length++
}

func (s *Sequence) Copy() Sequence {
    // Create a new Sequence
    newSeq := Sequence{
        Length: s.Length,
    }

    // Iterate over the original sequence and copy each node
    currentOriginal := s.Head
    var prevNewNode *Node

    for currentOriginal != nil {
        newNode := &Node{
            Value: currentOriginal.Value,
        }

        if prevNewNode != nil {
            // Link the new node to the previous one
            prevNewNode.Next = newNode
            newNode.Prev = prevNewNode
        } else {
            // Set the head of the new sequence
            newSeq.Head = newNode
            newSeq.IterPointer = newNode
        }

        prevNewNode = newNode
        currentOriginal = currentOriginal.Next
    }

    return newSeq
}

func (list *Sequence) InsertTail(char rune) {
	if list.Head == nil {
		list.InsertHead(char)
		return
	}

	if list.IterPointer == nil {
		// If IterPointer is nil, set it to the head
		list.IterPointer = list.Head
	} else {
		// Insert after IterPointer and update IterPointer
		list.InsertAfter(list.IterPointer, char)
		list.IterPointer = list.IterPointer.Next
	}
}

func (list *Sequence) InsertAfter(nodeToInsertAfter *Node, char rune) {
	// create new node
	newNode := &Node{
		Value: char,
		Next:  nodeToInsertAfter.Next,
		Prev:  nodeToInsertAfter,
	}

	// Update the next pointer of the previous node
	nodeToInsertAfter.Next = newNode

	// Update the previous pointer of the next node
	if newNode.Next != nil {
		newNode.Next.Prev = newNode
	}

	list.Length++
}

func (list *Sequence) Print() {
	current := list.Head
	for current != nil {
    fmt.Printf("%c", current.Value)
		// fmt.Print(current.Value, " ")
		current = current.Next
	}

  fmt.Println("    Sequence length:", list.Length)
}

func (list *Sequence) Debug() {
  // print all nodes and their pointers
  current := list.Head
  for current != nil {
    fmt.Println("Node:", current)
    fmt.Println("Value:", string(current.Value))
    fmt.Println("Prev:", current.Prev)
    fmt.Println("Next:", current.Next)
    fmt.Println("============")
    current = current.Next
  }
}

func (list *Sequence) Delete(nodeToDelete *Node) {
	if nodeToDelete == nil {
		return
	}

	if nodeToDelete.Prev != nil {
		nodeToDelete.Prev.Next = nodeToDelete.Next
	} else {
		list.Head = nodeToDelete.Next
	}

	if nodeToDelete.Next != nil {
		nodeToDelete.Next.Prev = nodeToDelete.Prev
	}

  list.Length--
}

// return the value of the current node and
// move the IterPointer to the next node
// if the end of the list is reached,
// return null charuner
func (list *Sequence) Yield() *Node {

  if list.IterPointer.Next != nil {
    list.IterPointer = list.IterPointer.Next
  } else {
    return nil
  }

  return list.IterPointer.Prev
}

func (list *Sequence) ResetIterPointer() {
  list.IterPointer = list.Head
}
