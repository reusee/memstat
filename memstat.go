package memstat

import (
	"fmt"
	"runtime"
	"sort"
	"strings"
)

type Stat struct {
	Pos   string
	InUse int64
}

type StatSorter []Stat

func (s StatSorter) Len() int           { return len(s) }
func (s StatSorter) Less(i, j int) bool { return s[i].InUse > s[j].InUse }
func (s StatSorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func Print() {
	records := make([]runtime.MemProfileRecord, 16)
	n, ok := runtime.MemProfile(records, false)
	for !ok {
		records = append(records, make([]runtime.MemProfileRecord, len(records))...)
		n, ok = runtime.MemProfile(records, false)
	}
	records = records[:n]
	inuseStat := make(map[string]int64)
	for _, record := range records {
		index := 0
	nextIndex:
		pc := record.Stack0[index]
		if pc == 0 {
			continue
		}
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		if strings.HasPrefix(fn.Name(), "runtime.") {
			index++
			goto nextIndex
		}
		file, line := fn.FileLine(pc)
		key := fmt.Sprintf("%s:%d", file, line)
		inuseStat[key] += record.InUseBytes()
	}
	var stats []Stat
	for pos, inuse := range inuseStat {
		stats = append(stats, Stat{pos, inuse})
	}
	sort.Sort(StatSorter(stats))
	for _, stat := range stats {
		fmt.Printf("%s %d\n", stat.Pos, stat.InUse)
	}
}
