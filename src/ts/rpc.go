package ts

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"strconv"
)

// Add your RPC definitions here.
type RPCArgs struct {
	CurGeneration   int
	IsNewChromExist bool
	NewChrom        Chrom
}

type RPCReply struct {
	TaskType       string
	CurGeneration  int
	CourseList     []Course
	TimeSlotList   []TimeSlot
	RoomList       []Room
	PrevGeneration []Chrom
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the master.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func masterSock() string {
	s := "/var/tmp/824-ts-"
	s += strconv.Itoa(os.Getuid())
	return s
}
