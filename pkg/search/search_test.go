package search

import (
	"context"
	"log"
	"testing"
)

func TestAll_user(t *testing.T) {

	ch := All(context.Background(), "HTTP", []string{"../../data/test.txt"})

	s, ok := <-ch
	if !ok {
		t.Errorf("function All error => %v", ok)
	}

	log.Println("s = ", s)

}

func TestAny_user(t *testing.T) {

	res := Any(context.Background(), "HTTP", []string{"../../test.txt", "../../test copy.txt"})

	r, ok := <-res
	if !ok {
		log.Println("error Any ok =>", ok)
	}

	log.Println("---------------")
	log.Println("res.Phrase) => ", r.Phrase)
	log.Println("res.Line) => ", r.Line)
	log.Println("res.LineNum) => ", r.LineNum)
	log.Println("res.ColNum) => ", r.ColNum)
	log.Println("---------------")

}