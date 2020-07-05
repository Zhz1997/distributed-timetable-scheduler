package ts

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

type DPVPair struct {
	gene  Gene
	index int
}

type ByFitness []Chrom

func (a ByFitness) Len() int           { return len(a) }
func (a ByFitness) Less(i, j int) bool { return a[i].FitnessScore < a[j].FitnessScore }
func (a ByFitness) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// print helper functions
func PrintGene(gene Gene) {
	fmt.Println("gene: courseID = ", gene.CourseID,
		" roomID = ", gene.RoomID,
		" timeSlotID = ", gene.TimeSlotID,
		" durationPref = ", gene.DurationPref)
}

func PrintChrom(chrom Chrom) {
	fmt.Println("chrom: ID = ", chrom.Id,
		" fitness = ", chrom.FitnessScore,
		" scv = ", chrom.Scv,
		" hcv = ", chrom.Hcv)
}

func PrintGeneration(generation []Chrom) {
	for _, chrom := range generation {
		PrintChrom(chrom)
	}
}

func CreateGene(courseID int, courseList []Course, timeSlotList []TimeSlot) Gene {
	var curGene Gene
	curGene.CourseID = courseID
	rand.Seed(time.Now().UTC().UnixNano())
	curGene.RoomID = rand.Intn(4-1) + 1
	rand.Seed(time.Now().UTC().UnixNano())
	curGene.TimeSlotID = rand.Intn(len(timeSlotList)) + 1
	curGene.DurationPref = courseList[courseID-1].Duration

	return curGene
}

func CreateChrom(id int, courseList []Course, timeSlotList []TimeSlot, roomList []Room) Chrom {
	var curChrom Chrom
	curChrom.Id = id
	var genes []Gene
	for i := 0; i < len(courseList); i++ {
		gene := CreateGene(i+1, courseList, timeSlotList)
		genes = append(genes, gene)
	}
	curChrom.Genes = genes
	curChrom.Hcv, curChrom.Scv, curChrom.FitnessScore = CalculateFitness(curChrom, timeSlotList, courseList, roomList)

	//PrintChrom(curChrom)

	return curChrom
}

func GetBestChromFromGen(gen []Chrom) Chrom {
	sort.Sort(sort.Reverse(ByFitness(gen)))
	result := gen[0]
	result.Id = 1
	return result
}

func CalculateFitness(chrom Chrom, timeSlotList []TimeSlot, courseList []Course, roomList []Room) (int, int, float64) {
	hcv := 0
	scv := 0
	// hard constraint violations
	roomCapViolation := 0
	timeViolation := 0
	var timeSlotConflictArray [3][][]int
	for i := 0; i < 3; i++ {
		timeSlotConflictArray[i] = make([][]int, len(timeSlotList))
	}

	// soft constraint violations - WIP (Phase 2?)
	durationPrefViolation := 0

	// loop through genes of the chrom to compute violations
	for i := 0; i < len(chrom.Genes); i++ {
		// room cap violation
		// fmt.Println(courseList)
		requiredRoomCap := courseList[chrom.Genes[i].CourseID-1].CapLevel
		actualRoomCap := roomList[chrom.Genes[i].RoomID-1].CapLevel

		if requiredRoomCap > actualRoomCap {
			roomCapViolation = roomCapViolation + 1
		}

		// time conflict violation
		if len(timeSlotConflictArray[chrom.Genes[i].RoomID-1][chrom.Genes[i].TimeSlotID-1]) != 0 {
			timeViolation = timeViolation + 1
			timeSlotConflictArray[chrom.Genes[i].RoomID-1][chrom.Genes[i].TimeSlotID-1] =
				append(timeSlotConflictArray[chrom.Genes[i].RoomID-1][chrom.Genes[i].TimeSlotID-1], chrom.Genes[i].CourseID)
		} else {
			var tempL []int
			tempL = append(tempL, chrom.Genes[i].CourseID)
			timeSlotConflictArray[chrom.Genes[i].RoomID-1][chrom.Genes[i].TimeSlotID-1] = tempL
		}

		// compute duration pref violations
		requiredDuration := chrom.Genes[i].DurationPref
		actualDuration := timeSlotList[chrom.Genes[i].TimeSlotID-1].Duration
		if requiredDuration != actualDuration {
			durationPrefViolation = durationPrefViolation + 1
		}

	}

	// fmt.Println(chrom.id, " : roomVio = ", roomCapViolation, ", timeVio = ", timeViolation)
	hcv = roomCapViolation + timeViolation
	scv = durationPrefViolation

	var fitness float64
	penalty := hcv*21 + scv
	fitness = (1 / float64(penalty)) * (1 / float64(penalty))

	return hcv, scv, fitness
}

func CreateFirstGeneration(courseList []Course, timeSlotList []TimeSlot, roomList []Room) []Chrom {
	var firstGen []Chrom
	for i := 0; i < genSize; i++ {
		chrom := CreateChrom(i+1, courseList, timeSlotList, roomList)
		firstGen = append(firstGen, chrom)
	}
	return firstGen
}

func CreateChromForNextGeneration(prevGen []Chrom, courseList []Course, timeSlotList []TimeSlot, roomList []Room) Chrom {

	sort.Sort(sort.Reverse(ByFitness(prevGen)))

	// choose parents
	fitnessT := float64(0)

	for _, chrom := range prevGen {
		fitnessT = fitnessT + chrom.FitnessScore
	}

	var l []int
	for i := 0; i < len(prevGen); i++ {
		numOfIns := math.Floor((prevGen[i].FitnessScore / fitnessT) * 100)
		for j := 0; j < int(numOfIns); j++ {
			l = append(l, i)
		}
	}

	// fmt.Println(l)

	var newChrom Chrom
	parentOne := prevGen[PickOne(l)]
	parentTwo := prevGen[PickOne(l)]

	// crossover
	newChrom = CrossOver(parentOne, parentTwo, 0, courseList, timeSlotList, roomList)

	// ifMutate := rand.Intn(11-1) + 1
	// if ifMutate < 4 {
	// 	newChrom = CreateChrom(newChrom.Id, courseList, timeSlotList, roomList)
	// }

	// return newChrom

	// evolve
	evolvedChrom := Evolve(newChrom, courseList, timeSlotList, roomList)

	// mutate
	rand.Seed(time.Now().UTC().UnixNano())
	ifMutate := rand.Intn(11-1) + 1
	if ifMutate < 4 {
		evolvedChrom = CreateChrom(evolvedChrom.Id, courseList, timeSlotList, roomList)
	}

	return evolvedChrom
}

func PickOne(l []int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return l[rand.Intn(len(l)-1)+1]
}

func DPVPairPickOne(l []DPVPair) (DPVPair, int) {
	rand.Seed(time.Now().UTC().UnixNano())
	removeIndex := rand.Intn(len(l)-1) + 1
	return l[removeIndex], removeIndex
}

func CrossOver(parentOne Chrom, parentTwo Chrom, id int, courseList []Course, timeSlotList []TimeSlot, roomList []Room) Chrom {
	midpoint := 0
	rand.Seed(time.Now().UTC().UnixNano())
	midpoint = rand.Intn(len(courseList)-1) + 1
	var newGenes []Gene
	for i := 0; i < midpoint; i++ {
		newGenes = append(newGenes, parentOne.Genes[i])
	}
	for i := midpoint; i < len(courseList); i++ {
		newGenes = append(newGenes, parentTwo.Genes[i])
	}

	var newChrom Chrom
	newChrom.Id = id
	newChrom.Genes = newGenes
	newChrom.Hcv, newChrom.Scv, newChrom.FitnessScore = CalculateFitness(newChrom, timeSlotList, courseList, roomList)

	return newChrom
}

func Evolve(newChrom Chrom, courseList []Course, timeSlotList []TimeSlot, roomList []Room) Chrom {
	// evolve room cap
	evolvedChrom := newChrom
	if newChrom.Hcv > 0 {
		for i := range evolvedChrom.Genes {
			if courseList[evolvedChrom.Genes[i].CourseID-1].CapLevel > roomList[evolvedChrom.Genes[i].RoomID-1].CapLevel {
				evolvedChrom.Genes[i].RoomID = evolvedChrom.Genes[i].RoomID + 1
			}
		}
	}

	return evolvedChrom
}
