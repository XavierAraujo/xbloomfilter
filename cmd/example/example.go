package main

import (
	"log"
	"xbloomfilter/xbloomfilter"
)

type Person struct {
	ID   string
	Name string
}

func main() {
	bloomfilter := xbloomfilter.NewBloomFilter[Person](1000, 0.01, xbloomfilter.MurmurHasher)

	p1 := Person{ID: "1", Name: "Steve"}
	p2 := Person{ID: "2", Name: "John"}
	p3 := Person{ID: "3", Name: "Martha"}
	p4 := Person{ID: "4", Name: "Kelly"}

	bloomfilter.Add(p1)
	bloomfilter.Add(p2)
	bloomfilter.Add(p3)

	r1, _ := bloomfilter.MightContain(p1)
	r2, _ := bloomfilter.MightContain(p2)
	r3, _ := bloomfilter.MightContain(p3)
	r4, _ := bloomfilter.MightContain(p4)

	log.Println("Person ", p1, " in filter: ", r1)
	log.Println("Person ", p2, " in filter: ", r2)
	log.Println("Person ", p3, " in filter: ", r3)
	log.Println("Person ", p4, " in filter: ", r4)
}
