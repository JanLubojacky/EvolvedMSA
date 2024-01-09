package pkg

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func NewSequence(data string) *Sequence {
	seq := &Sequence{}
	seq.IterPointer = seq.Head
  seq.Length = 0

	for _, char := range data {
    seq.InsertTail(char)
	}
  seq.ResetIterPointer()
	return seq
}

func ParseFASTA(filename string) (MSA, error) {
	file, err := os.Open(filename)
	if err != nil {
		return MSA{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	msa := MSA{}

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, ">target sequence") {
			scanner.Scan()
			targetData := scanner.Text()
      // capitalize all the letters
      targetData = strings.ToUpper(targetData)
			msa.TargetSequence = *NewSequence(targetData)
		} else if strings.HasPrefix(line, ">sequence") {
			scanner.Scan()
			otherData := scanner.Text()
      // capitalize all the letters
      otherData = strings.ToUpper(otherData)
			otherSeq := *NewSequence(otherData)
			msa.OtherSequences = append(msa.OtherSequences, otherSeq)
		}
	}

	if err := scanner.Err(); err != nil {
		return MSA{}, err
	}

	return msa, nil
}

func DisplaySequence(sequence *Sequence) {
	current := sequence.Head
	for current != nil {
		fmt.Printf("%c", current.Value)
		current = current.Next
	}
	fmt.Println()
}
