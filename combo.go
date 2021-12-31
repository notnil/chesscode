package chesscode

import (
	"sort"

	"github.com/notnil/chess"
)

var (
	char2RevLookup  []string
	char2Lookup     map[string]int
	char3RevLookup  []string
	char3Lookup     map[string]int
	combo2by2Lookup [][]combo2by2Result
)

func init() {
	char2RevLookup, char2Lookup = char2Combos()
	char3RevLookup, char3Lookup = char3Combos()
	combo2by2Lookup = combo2by2(52)
}

func combo2(n int) [][]int {
	results := [][]int{}
	for i1 := 0; i1 < n; i1++ {
		results = append(results, []int{i1})
		for i2 := i1 + 1; i2 < n; i2++ {
			results = append(results, []int{i1, i2})
		}
	}
	return results
}

func combo4(n int) [][]int {
	results := [][]int{}
	for i1 := 0; i1 < n; i1++ {
		results = append(results, []int{i1})
		for i2 := i1 + 1; i2 < n; i2++ {
			results = append(results, []int{i1, i2})
			for i3 := i2 + 1; i3 < n; i3++ {
				results = append(results, []int{i1, i2, i3})
				for i4 := i3 + 1; i4 < n; i4++ {
					results = append(results, []int{i1, i2, i3, i4})
				}
			}
		}
	}
	return results
}

type combo2by2Result struct {
	idx   int
	piece chess.Piece
}

func combo2by2(n int) [][]combo2by2Result {
	r1 := combo2(n)
	r2 := combo2(n)
	results := [][]combo2by2Result{}
	for _, c1 := range r1 {
		for _, c2 := range r2 {
			r := []combo2by2Result{}
			m := map[int]bool{}
			for _, i := range c1 {
				r = append(r, combo2by2Result{idx: i, piece: chess.WhiteKnight})
				m[i] = true
			}
			for _, i := range c2 {
				r = append(r, combo2by2Result{idx: i, piece: chess.BlackKnight})
				m[i] = true
			}
			if len(m) != len(r) {
				continue
			}
			sort.Slice(r, func(i, j int) bool {
				return r[i].idx < r[j].idx
			})
			results = append(results, r)
		}
	}
	return results
}

func char2Combos() ([]string, map[string]int) {
	a := []string{}
	m := map[string]int{}
	for _, c1 := range charset {
		for _, c2 := range charset {
			s := string(c1) + string(c2)
			a = append(a, s)
			m[s] = len(a) - 1
		}
	}
	return a, m
}

func char3Combos() ([]string, map[string]int) {
	a := []string{}
	m := map[string]int{}
	for _, c1 := range charset {
		for _, c2 := range charset {
			for _, c3 := range charset {
				s := string(c1) + string(c2) + string(c3)
				a = append(a, s)
				m[s] = len(a) - 1
			}
		}
	}
	return a, m
}
