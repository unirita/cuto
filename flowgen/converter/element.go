package converter

import "fmt"

// Element is a node in the flow description.
// The gateway block (a block which enclosed in bracket) count as one Element.
type Element interface {
	ID() string
	SetNext(Element)
	Next() Element
}

// Implementation struct of Element interface.
type element struct {
	id   string
	next Element
}

// Get the element id.
func (e *element) ID() string {
	return e.id
}

// Get the next element.
func (e *element) Next() Element {
	return e.next
}

// Set n as a next element.
func (e *element) SetNext(n Element) {
	e.next = n
}

// Job is a job element in the flow description.
type Job struct {
	element
	name string
}

const jPrefix = "job"

var jIndex int = 1

// Create new Job object with name.
func NewJob(name string) *Job {
	j := new(Job)
	j.id = fmt.Sprintf("%s%d", jPrefix, jIndex)
	j.name = name
	jIndex++
	return j
}

// Get the job name.
func (j *Job) Name() string {
	return j.name
}

const gPrefix = "gw"

var gIndex int = 1

// Gateway is a block which enclosed in bracket in the flow description.
// Gateway has paths which consists of jobs.
type Gateway struct {
	element

	// Head job of paths.
	PathHeads []*Job
}

// Create new Gateway object with no path.
func NewGateway() *Gateway {
	g := new(Gateway)
	g.id = fmt.Sprintf("%s%d", gPrefix, gIndex)
	g.PathHeads = make([]*Job, 0)
	gIndex++
	return g
}

// Add path to Gateway.
func (g *Gateway) AddPathHead(head *Job) {
	g.PathHeads = append(g.PathHeads, head)
}
