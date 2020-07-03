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
)

type Master struct {
	NReduce int
	NMap    int
}

var courseList []Course
var roomList []Room
var timeSlotList []TimeSlot

func (m *Master) RPCHandler(args *RPCArgs, reply *RPCReply) error {

	return nil
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
	fmt.Println(filePath)
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

	fmt.Println(timeSlotList)

	m := Master{}

	m.server()
	return &m
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
			id:       curID,
			capLevel: curCapLevel,
			duration: curDuration,
			timePref: curTimePref,
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
			id:       curID,
			capLevel: curCapLevel,
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
			id:         curID,
			startTime:  curStartTime,
			duration:   len(temp),
			isOccupied: false,
		}
		timeSlotList = append(timeSlotList, data)
	}
}
