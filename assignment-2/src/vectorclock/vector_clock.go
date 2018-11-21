package vectorclock

import (
	"bytes"
	"fmt"
)

type VectorClock struct {
	Id          int
	clockVector map[int]int
	chans       []chan map[int]int
}

func NewVectorClock(newIdIdx int, ids []int, chans []chan map[int]int) (vc VectorClock, err error) {
	//Make sure the index is valid
	if newIdIdx >= len(ids) || newIdIdx < 0 {
		err = fmt.Errorf("Error: Index value %d must be between 0 and %d", newIdIdx, len(ids))
		return
	}

	vc.clockVector = make(map[int]int)
	vc.Id = ids[newIdIdx]
	vc.chans = chans

	//Start all clocks at 0
	for _, id := range ids {
		vc.clockVector[id] = 0
	}

	return
}

func (this *VectorClock) GetClock(id int) (clock int, err error) {
	clock, exists := this.clockVector[id]
	if !exists {
		err = fmt.Errorf("Error: Id %d does not exist", id)
	}

	return
}

func (this *VectorClock) Receive(otherId int) {
	otherClockVector := <-this.chans[this.Id]
	for id, _ := range this.clockVector {
		this.clockVector[id] = max(this.clockVector[id], otherClockVector[id])
	}
}

func (this *VectorClock) Send(otherId int) {
	vectorCopy := make(map[int]int)
	for key, value := range this.clockVector {
		vectorCopy[key] = value
	}

	this.chans[otherId] <- vectorCopy
}

func (this *VectorClock) Inc() {
	this.clockVector[this.Id]++
}

func (this *VectorClock) GetClockString() string {
	var b bytes.Buffer
	clocks := make([]int, len(this.clockVector))

	for id, clock := range this.clockVector {
		clocks[id] = clock
	}

	for id := range clocks {
		b.WriteString(fmt.Sprintf("%d:%d ", id, clocks[id]))
	}
	return b.String()
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
