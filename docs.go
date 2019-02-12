package indexr

import "sort"

type Docs []int

// Implement sort.Interface
func (d Docs) Len() int           { return len(d) }
func (d Docs) Swap(a, b int)      { d[a], d[b] = d[b], d[a] }
func (d Docs) Less(a, b int) bool { return d[a] < d[b] }

func sortDocs(d Docs) Docs {
	sort.Sort(d)
	return d
}

// mapToDocs converts a map to a document array
func mapToDocs(m map[int]bool) Docs {
	result := Docs{}
	for k := range m {
		result = append(result, k)
	}
	sort.Sort(result)
	return result
}

// intersect returns the intersection between two arrays
func intersect(a, b Docs) Docs {
	result := Docs{}
	i, j := 0, 0

	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			result = append(result, a[i])
			i++
			j++
			if i >= len(a) || j >= len(b) {
				return result
			}
		}

		for a[i] < b[j] {
			i++
			if i >= len(a) {
				return result
			}
		}

		for j < len(b) && a[i] > b[j] {
			j++
			if j >= len(b) {
				return result
			}
		}
	}

	return result
}
