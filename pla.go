package main

import "math/rand"
import "time"
import "strconv"
import "fmt"
import "os"
import "flag"

type sign int8

type point struct {
	x float32
	y float32
}

type samplePoint struct {
	p point
	s sign
}

type sample struct {
	p1, p2  point
	ws      []float32
	inputs  []samplePoint
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		os.Exit(1)
	}
	times, _ := strconv.Atoi(args[0])
	num, _   := strconv.Atoi(args[1])
	fmt.Println(bigTest(times, num))
	os.Exit(0)
}

func bigTest(times, count int) (float32, float32) {
	rand.Seed(time.Now().UnixNano())
	cnum := 0
	mnum := 0
	for i := 0; i < times; i++ {
		s := newSample(count)
		cnum += s.correct()
		mnum += s.misclassify()
	}
	return float32(cnum)/float32(times), float32(mnum) / float32(times)
}

func (s *sample) correct() (n int) {
	for {
		fault := []samplePoint{}
		for i, in := range s.inputs {
			if s.hypoIsMatched(i) == false {
				fault = append(fault, in)
			}
		}
		if length := len(fault); length > 0 {
			i := rand.Int31n(int32(length))
			s.hypoCorr(fault[i])
			n++
			continue
		}
		break
	}
	return 
}

func (s *sample) misclassify() int {
	p := randPoint()
	mark := s.origTest(p)
	if s.hypoTest(p) != mark {
		return 1
	}
	return 0
}

func newSample(n int) *sample {
	s := sample{
		p1 : randPoint(),
		p2 : randPoint(),
		ws : make([]float32, 3),
		inputs  : make([]samplePoint, n),
	}
	for i, _ := range s.inputs {
		p := randPoint()
		t := s.origTest(p)
		s.inputs[i] = samplePoint{p, t}
	}
	return &s
}

func randPoint() point {
	return point{rand.Float32()*2-1, rand.Float32()*2-1}
}

func (s *sample) hypoCorr(sp samplePoint) {
	s.ws[0] += float32(sp.s)
	s.ws[1] += sp.p.x * float32(sp.s)
	s.ws[2] += sp.p.y * float32(sp.s)
}

func (s *sample) hypoIsMatched(i int) (t bool) {
	switch m := s.ws[0] + s.ws[1]*s.inputs[i].p.x + s.ws[2]*s.inputs[i].p.y; {
	case s.inputs[i].s == 1  && m >= 0: t = true
	case s.inputs[i].s == -1 && m < 0: t = true
	default: t = false
	}
	return
}

func (s *sample) hypoTest(p point) sign {
	if s.ws[0] + s.ws[1]*p.x + s.ws[2]*p.y >= 0 {
		return 1
	} 
	return -1
}

func (s *sample) origTest(p point) sign {
	if (s.p2.x - s.p1.x)*(p.y - s.p1.y) - (s.p2.y - s.p1.y)*(p.x - s.p1.x) >= 0 {
		return 1
	}
	return -1
}