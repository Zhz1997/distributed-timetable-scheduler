package ts

import (
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"time"
)

type Master struct {
	NReduce int
	NMap    int
}

var courseList []Course
var roomList []Room
var timeSlotList []TimeSlot
var prevGeneration []Chrom
var wipGen []Chrom

var genSize int = 20
var numberOfGen int = 3000
var curGenCount int = 0
var curGenChromCount int = 0

func (m *Master) RPCHandler(args *RPCArgs, reply *RPCReply) error {
	if curGenCount == numberOfGen {
		reply.TaskType = "done"

		return nil
	}
	// accept chrom from worker and reassign
	if args.CurGeneration == curGenCount {
		// fmt.Println("----------------------------------")
		// PrintGeneration(wipGen)
		// fmt.Println("----------------------------------")
		reply.TaskType = "train"

		reply.CourseList = courseList
		reply.RoomList = roomList
		reply.TimeSlotList = timeSlotList
		reply.PrevGeneration = prevGeneration
		reply.CurGeneration = curGenCount

		if args.IsNewChromExist && len(wipGen) < genSize {
			wipGen = append(wipGen, args.NewChrom)
			curGenChromCount = curGenChromCount + 1
		}

		return nil

	} else { //reject the chrom from worker, and tell worker cur gen count
		reply.TaskType = "reassign"
		reply.CurGeneration = curGenCount

		return nil
	}

	return nil
}

func Monitor() {
	for true {
		time.Sleep(1 * time.Millisecond)
		if len(wipGen) >= genSize {
			wipGen = wipGen[:genSize]
			bestChromInPrevGen := GetBestChromFromGen(prevGeneration)
			wipGen = SetID(wipGen)
			prevGeneration = wipGen
			PrintBestInPrev()
			wipGen = nil
			wipGen = append(wipGen, bestChromInPrevGen)
			curGenCount = curGenCount + 1
		}
	}
}

func PrintBestInPrev() {

	bestFitValue := float64(0)
	bestChromId := 0

	for _, chrom := range prevGeneration {
		if chrom.FitnessScore > bestFitValue {
			bestChromId = chrom.Id
			bestFitValue = chrom.FitnessScore
		}
	}

	fmt.Println("------------------------------")
	PrintGeneration(prevGeneration)
	fmt.Println("bestFitValue in generation ", curGenCount, " ", "is ", bestFitValue, " chrom id is ", bestChromId)
}

func SetID(gen []Chrom) []Chrom {
	var result []Chrom
	result = gen
	for i := range result {
		result[i].Id = i + 1
	}

	return result
}

//
// start a thread that listens for RPCs from worker.go
//
func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/tsmaster.go calls Done() periodically to find out
// if the entire job has finished.
//
func (m *Master) Done() bool {
	ret := false

	return ret
}

//
// create a Master.
// main/tsmaster.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeMaster(filePath string, nReduce int) *Master {

	WriteDataIntoLists(filePath)

	// fmt.Println(courseList)
	// fmt.Println(roomList)
	// fmt.Println(timeSlotList)

	var firstGeneration []Chrom
	firstGeneration = CreateFirstGeneration(courseList, timeSlotList, roomList)
	//fmt.Println(len(firstGeneration))
	//PrintGeneration(firstGeneration)
	// for _, chrom := range firstGeneration {
	// 	for _, gene := range chrom.genes {
	// 		PrintGene(gene)
	// 	}
	// }

	bestFitValue := float64(0)
	bestChromId := 0

	for _, chrom := range firstGeneration {
		if chrom.FitnessScore > bestFitValue {
			bestChromId = chrom.Id
			bestFitValue = chrom.FitnessScore
		}
	}

	prevGeneration = firstGeneration

	bestChromInPrevGen := GetBestChromFromGen(prevGeneration)
	wipGen = append(wipGen, bestChromInPrevGen)
	fmt.Println("------------------------------")
	PrintGeneration(firstGeneration)
	fmt.Println("bestFitValue in initial generation is ", bestFitValue, " chrom id is ", bestChromId)

	m := Master{}

	go Monitor()
	m.server()
	return &m
}

func WriteDataIntoLists(filePath string) {
	courses, err := ReadCsv(filePath + "/courses.csv")
	if err != nil {
		panic(err)
	}
	courses = courses[1:]
	SetCourses(courses)

	rooms, err := ReadCsv(filePath + "/rooms.csv")
	if err != nil {
		panic(err)
	}
	rooms = rooms[1:]
	SetRooms(rooms)

	timeSlots, err := ReadCsv(filePath + "/timeSlots.csv")
	if err != nil {
		panic(err)
	}
	timeSlots = timeSlots[1:]
	SetTimeSlots(timeSlots)
}

// ReadCsv accepts a file and returns its content as a multi-dimentional type
// with lines and each column. Only parses to string type.
func ReadCsv(filename string) ([][]string, error) {

	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}

func SetCourses(lines [][]string) {
	for _, line := range lines {
		curID, err := strconv.Atoi(line[0])
		if err != nil {
			panic(err)
		}

		curCapLevel, err := strconv.Atoi(line[1])
		if err != nil {
			panic(err)
		}

		curDuration, err := strconv.Atoi(line[2])
		if err != nil {
			panic(err)
		}

		curTimePref, err := strconv.Atoi(line[3])
		if err != nil {
			panic(err)
		}

		data := Course{
			Id:       curID,
			CapLevel: curCapLevel,
			Duration: curDuration,
			TimePref: curTimePref,
		}
		courseList = append(courseList, data)
	}
}

func SetRooms(lines [][]string) {
	for _, line := range lines {
		curID, err := strconv.Atoi(line[0])
		if err != nil {
			panic(err)
		}

		curCapLevel, err := strconv.Atoi(line[1])
		if err != nil {
			panic(err)
		}

		data := Room{
			Id:       curID,
			CapLevel: curCapLevel,
		}
		roomList = append(roomList, data)
	}
}

func SetTimeSlots(lines [][]string) {
	for _, line := range lines {
		curID, err := strconv.Atoi(line[0])
		if err != nil {
			panic(err)
		}

		temp := strings.Split(line[2], "&")

		curStartTime, err := strconv.Atoi(temp[0])
		if err != nil {
			panic(err)
		}

		data := TimeSlot{
			Id:         curID,
			StartTime:  curStartTime,
			Duration:   len(temp),
			IsOccupied: false,
		}
		timeSlotList = append(timeSlotList, data)
	}
}
