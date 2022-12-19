package flightpath

import (
    "fmt"
    "encoding/json"
    mapset "github.com/deckarep/golang-set/v2"
)

type Flight struct {
    src, dst string
}

func (f *Flight) MarshalJSON() ([]byte, error) {
    if f.src == "" || f.dst == "" {
        return nil, fmt.Errorf("incomplete flight: [\"%s\", \"%s\"]", f.src, f.dst)
    }
    return []byte(fmt.Sprintf("[\"%s\",\"%s\"]", f.src, f.dst)), nil
}

func (f *Flight) UnmarshalJSON(bs []byte) (err error) {
    var ss []string
    if err = json.Unmarshal(bs, &ss); err == nil {
        if l := len(ss); l != 2 {
            err = fmt.Errorf("flight has more than 2 airports: %d %s", l, ss)
        } else if ss[0] == "" || ss[1] == "" {
            err = fmt.Errorf("one or more airport in flight is empty")
        } else {
            f.src, f.dst = ss[0], ss[1]
        }
    }
    return
}

func (f *Flight) String() string {
    return fmt.Sprintf("[%s->%s]", f.src, f.dst)
}


type FlightPath struct {
    Flights []*Flight `json:"flights"`
    Chain []string `json:"flightChain"`
    Layovers []string `json:"layovers"`
    Start string `json:"start"`
    End string `json:"end"`
}

func (fp *FlightPath) UnmarshalJSON(bs []byte) error {
    return json.Unmarshal(bs, &fp.Flights)
}


func (fp *FlightPath) Add(src, dst string) {
    fp.Flights = append(fp.Flights, &Flight{src: src, dst: dst})
}

func (fp *FlightPath) BuildChain() error {
    if len(fp.Flights) < 1 {
        return fmt.Errorf("No flights passed")
    }

    srcSet := mapset.NewSet[string]()
    dstSet := mapset.NewSet[string]()

    // build sets of departing and arriving airports
    // check if airports are departed from or arrived at more than once
    for _, v := range fp.Flights {
        if srcSet.Contains(v.src) {
            return fmt.Errorf("multiple flights departing from %s", v.src)
        }
        if dstSet.Contains(v.dst) {
            return fmt.Errorf("multiple flights arriving at %s", v.dst)
        }
        srcSet.Add(v.src)
        dstSet.Add(v.dst)
    }


    // unique arriving airports == unique departing airports == total flight count
    if srcSet.Cardinality() != len(fp.Flights) {
        return fmt.Errorf("unique departing airports (%d) doesn't match total flight count (%d)", srcSet.Cardinality(), len(fp.Flights))
    } else if dstSet.Cardinality() != len(fp.Flights) {
        return fmt.Errorf("unique arriving airports (%d) doesn't match total flight count (%d)", dstSet.Cardinality(), len(fp.Flights))
    }

    route := make(map[string]string)
    for _, v := range fp.Flights {
        route[v.src] = v.dst
    }

    // If the set of flights truly forms a chain, then the set difference between source
    // airports and destination airports (and vice versa) should be singleton sets with
    // the overall total start/end of the flight chain. However, there could still be loops
    // in the flights (i.e. [PDX, LAX], [LAX, ATL], [ATL, PDX]). 
    s := srcSet.Difference(dstSet)
    d := dstSet.Difference(srcSet)

    switch c := s.Cardinality(); c {
        case 0:
            loop := []string{}
            if cur, ok := srcSet.Pop(); ok {
                ini := cur
                loop = []string{ini}
                for {
                    if nxt, ok := route[cur]; ok {
                        loop = append(loop, nxt)
                        if nxt == ini {
                            break
                        }
                        cur = nxt
                    } else {
                        break
                    }
                }
            }
            return fmt.Errorf("no starting airports due to loop: %s", loop)
        case 1:
        default:
            return fmt.Errorf("too many starting airports: %v", s.ToSlice())
    }

    if d.Cardinality() != 1 {
        return fmt.Errorf("too many ending airports: %v", d.ToSlice())
    }

    arrived := mapset.NewSet[string]() // to detect loops
    cur, _ := s.Pop()
    chain := []string{cur}
    arrived.Add(cur)

    // Build the flight chain, check for loops
    for {
        if nxt, ok := route[cur]; ok {
            if arrived.Contains(nxt) {
                return fmt.Errorf("loop: arriving at %s, but already departed %s", nxt, nxt)
            } else {
                arrived.Add(nxt)
                chain = append(chain, nxt)
                cur = nxt
            }
        } else {
            break // end of chain
        }
    }

    if len(chain) == len(fp.Flights) + 1 {
        fp.Chain = chain
        if len(chain) > 2 {
            fp.Layovers = chain[1:len(chain)-1]
        }
        fp.Start = chain[0]
        fp.End = chain[len(chain)-1]
        return nil
    } else { // incomplete chain, likely due to a loop
        loop := []string{}
        if ini, ok := dstSet.Difference(arrived).Pop(); ok {
            cur = ini
            loop = []string{cur}
            for {
                if nxt, ok := route[cur]; ok {
                    loop = append(loop, nxt)
                    if nxt == ini {
                        break
                    }
                    cur = nxt
                }
            }
        } else {
            return fmt.Errorf("incomplete chain")
        }
        return fmt.Errorf("incomplete chain: %v, at least one loop: %v", chain, loop)
    }
}
