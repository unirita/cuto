package converter

// Element is a node in the flow description.
// The parallel block (a block which enclosed in bracket) count as one Element.
type Element interface {
	SetNext(Element)
	Next() Element
}

// Implementation struct of Element interface.
type element struct {
	next Element
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

	// Job name
	Name string
}

// Create new Job object with name.
func NewJob(name string) *Job {
	j := new(Job)
	j.Name = name
	return j
}

// Parallel is a block which enclosed in bracket in the flow description.
// Parallel has paths which consists of jobs.
type Parallel struct {
	element

	// Head job of paths.
	PathHeads []*Job
}

// Create new Parallel object with no path.
func NewParallel() *Parallel {
	p := new(Parallel)
	p.PathHeads = make([]*Job, 0)
	return p
}

// Add path to Parallel.
func (p *Parallel) addPathHead(head *Job) {
	p.PathHeads = append(p.PathHeads, head)
}
