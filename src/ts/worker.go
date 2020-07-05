package ts

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

func Worker() {

	args := RPCArgs{}
	args.IsNewChromExist = false
	args.CurGeneration = -1
	reply := RPCReply{}

	for reply.TaskType != "done" {
		time.Sleep(3 * time.Millisecond)
		call("Master.RPCHandler", &args, &reply)
		fmt.Println(reply.CurGeneration)
		if reply.TaskType == "reassign" {
			args.CurGeneration = reply.CurGeneration
			fmt.Println("reassigned")
		} else if reply.TaskType == "wait" {

		} else {
			fmt.Println("asdfasd")
			args.NewChrom = CreateChromForNextGeneration(reply.PrevGeneration, reply.CourseList, reply.TimeSlotList, reply.RoomList)
			args.IsNewChromExist = true
			PrintChrom(args.NewChrom)
		}
	}
	// call("Master.RPCHandler", &args, &reply)
	// //fmt.Println(reply.CourseList)
	// nextGen := CreateNextGeneration(reply.PrevGeneration, reply.CourseList, reply.TimeSlotList, reply.RoomList)

	// bestFitValue := float64(0)
	// bestChromId := 0

	// for _, chrom := range nextGen {
	// 	if chrom.FitnessScore > bestFitValue {
	// 		bestChromId = chrom.Id
	// 		bestFitValue = chrom.FitnessScore
	// 	}
	// }

	// fmt.Println("------------------------------")
	// PrintGeneration(nextGen)
	// fmt.Println("bestFitValue in initial generation is ", bestFitValue, " chrom id is ", bestChromId)

	// fmt.Println(PrevGeneration)
}

//
// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := masterSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
