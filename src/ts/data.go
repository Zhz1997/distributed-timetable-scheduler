package ts

type Room struct {
	id       int
	capLevel int
}

type TimeSlot struct {
	id         int
	startTime  int
	duration   int
	isOccupied bool
}

type Course struct {
	id       int
	capLevel int
	duration int
	timePref int
}

type Gene struct {
	courseID     int
	roomID       int
	timeSlotID   int
	durationPref int
}

type Chrom struct {
	id           int
	genes        []Gene
	fitnessScore int
	scv          int
	hcv          int
}
