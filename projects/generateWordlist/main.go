package main

import (
	"fmt"
	"log"
	"regexp"
	"sync"
)

type global struct {
	expressions   []expression
	numGenerators int
	numValidators int
}

var G *global

var AllCharacters []byte = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890~!@#$%^&*()_+=-./<>?,;: ")

func generator(vCh <-chan []byte, gCh chan<- []byte, wg *sync.WaitGroup) {
	for validated := range vCh {
		go func() {
			defer wg.Done()
			for _, b := range AllCharacters {
				generated := make([]byte, 0, len(validated)+1)
				generated = append(generated, validated...)
				generated = append(generated, b)
				wg.Add(1)
				gCh <- generated
			}
		}()
	}
}

func validator(gCh <-chan []byte, vCh chan<- []byte, resultCh chan<- string, wg *sync.WaitGroup) {

	for generated := range gCh {
		go func() {
			defer wg.Done()
			for _, ex := range G.expressions {
				status := validate(generated, ex)
				if status.matched {
					if status.result {
						resultCh <- string(generated)
					}
					if status.toBeContinued {
						wg.Add(1)
						vCh <- generated
					}
				}
			}
		}()
	}
}

type input struct {
	expression string
	sample     string
	minLength  int
}

func main() {
	G = &global{}
	inputs := []input{
		{
			expression: `linn\d{3}@gmail\.com`,
			sample:     "linn221@gmail.com",
		},
	}
	for _, input := range inputs {
		ex := expression{
			re:        regexp.MustCompile(input.expression),
			sample:    []byte(input.sample),
			minLength: input.minLength,
		}
		G.expressions = append(G.expressions, ex)
	}
	G.numGenerators = 100
	G.numValidators = 100
	var wg sync.WaitGroup
	// channels
	validated := make(chan []byte)
	generated := make(chan []byte)
	resultCh := make(chan string)

	for range G.numGenerators {
		go generator(validated, generated, &wg)
	}
	for range G.numValidators {
		go validator(generated, validated, resultCh, &wg)
	}

	wg.Add(1)
	validated <- []byte{}

	go func() {
		wg.Wait()
		close(resultCh)
	}()
	for result := range resultCh {
		fmt.Println(result)
	}
	log.Println("done")
}
