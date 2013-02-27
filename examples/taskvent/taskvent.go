//
//  Task ventilator.
//  Binds PUSH socket to tcp://localhost:5557
//  Sends batch of tasks to workers via that socket
//
package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq3"
	"math/rand"
	"time"
)

func main() {
	context, _ := zmq.NewContext()

	//  Socket to send messages on
	sender, _ := context.NewSocket(zmq.PUSH)
	sender.Bind("tcp://*:5557")

	//  Socket to send start of batch message on
	sink, _ := context.NewSocket(zmq.PUSH)
	sink.Connect("tcp://localhost:5558")

	fmt.Print("Press Enter when the workers are ready: ")
	var line string
	fmt.Scanln(&line)
	fmt.Println("Sending tasks to workers…")

	//  The first message is "0" and signals start of batch
	sink.Send("0", 0)

	//  Initialize random number generator
	rand.Seed(time.Now().Unix())

	//  Send 100 tasks
	total_msec := 0
	for task_nbr := 0; task_nbr < 100; task_nbr++ {
		//  Random workload from 1 to 100msecs
		workload := rand.Intn(100) + 1
		total_msec += workload
		s := fmt.Sprintf("%d", workload)
		sender.Send(s, 0)
	}
	fmt.Println("Total expected cost: " +  (time.Duration(total_msec) * time.Millisecond).String())
	time.Sleep(time.Second) //  Give 0MQ time to deliver

}
