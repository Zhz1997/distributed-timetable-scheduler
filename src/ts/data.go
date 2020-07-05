package ts

type Room struct {
	Id       int
	CapLevel int
}

type TimeSlot struct {
	Id         int
	StartTime  int
	Duration   int
	IsOccupied bool
}

type Course struct {
	Id       int
	CapLevel int
	Duration int
	TimePref int
}

type Gene struct {
	CourseID     int
	RoomID       int
	TimeSlotID   int
	DurationPref int
}

type Chrom struct {
	Id           int
	Genes        []Gene
	FitnessScore float64
	Scv          int
	Hcv          int
}
