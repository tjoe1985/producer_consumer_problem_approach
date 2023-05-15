package main

import (
	"fmt"
	"math/rand"
	"time"
)

const NumberOfPizzas = 10

var pizzsasMade, pizzasFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func main() {
	// seed the random number generator
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// print out message
	fmt.Println("We are open now.......")
	fmt.Println("...........................")
	// crete a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}
	//run the producer in the background
	go Pizzeria(pizzaJob)
	// create and run consumer
	for i := range pizzaJob.data {
		// try to make a pizza if i is less than the NumberOfPizzas.
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				fmt.Println(i.message)
				fmt.Printf("Order #%d is out for delivery!", i.pizzaNumber)

			} else {
				fmt.Println(i.message)
				fmt.Println("The client is hangry")
			}
		} else {
			fmt.Println("Done making pizzas")
			//close our channels
			err := pizzaJob.Close()
			if err != nil {
				fmt.Println("*** Error closing channel.", err)
			}
		}
	}
	// print out the ending message
	fmt.Println("We are now closed.......")
	fmt.Println("...........................")
	fmt.Printf("Toda we made %d pizzas but failed to make %d, with a total of %d attempts.", pizzsasMade, pizzasFailed, total)

}
func Pizzeria(pizzaMaker *Producer) {
	// keep track of which pizza we are preparing
	var count = 0

	// run until we get quit signal
	// try to produce pizza
	for {
		currentPizza := makePizza(count)
		if currentPizza != nil {
			count = currentPizza.pizzaNumber
			select {
			// we tried to make a pizza (only means we sent info to the chanel
			case pizzaMaker.data <- *currentPizza:

			case quitChan := <-pizzaMaker.quit:
				// close channel
				close(pizzaMaker.data)
				//close quitChan
				close(quitChan)
				//exit out of the go routine
				return
			}
		}
		// decision

	}
}

// Close will be available to all producer it returns a channel with an error wish will be empty if there is an error or contain the actual error.
func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func makePizza(count int) *PizzaOrder {
	//increment the ammount of pizzas
	count++
	//make sure we don't do more thant the limit of 10 pizzas
	if count <= NumberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Println("Received order #", count)
		// generate a random number with a chance for something to go wrong
		rnd := rand.Intn(12) + 1
		msg := ""
		success := false
		// rnd is less than 5 our pizza failed.
		if rnd <= 5 {
			pizzasFailed++
		} else {
			pizzsasMade++
		}
		// increment total no matter of what since we are keeping track of all atempts
		total++

		fmt.Printf("Making pizza #%d. it will take %d seconds...\n ", count, delay)
		//delay as needed to mimic production
		time.Sleep(time.Duration(delay) * time.Second)
		if rnd <= 2 {
			msg = fmt.Sprintf("*** we ran out of ingredients for pizza #%d ! ", count)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** the cook got on a fight while making pizza #%d ! ", count)
		} else {
			success = true
			msg = fmt.Sprintf("*** Pizza order #%d is hot and ready to go! ", count)
		}
		p := PizzaOrder{
			pizzaNumber: count,
			message:     msg,
			success:     success,
		}
		return &p
	}
	return &PizzaOrder{
		pizzaNumber: count,
	}
}
