package main

import (
    "flag"
    "fmt"
	"io/ioutil"
	"errors"
	
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)



type coreStruct struct{
	cpu [2]int
	containerType string
	podIndex int
	containerIndex int
	splitL3Pod int
}
type l3GroupStruct struct {
	ID   int
	numCores int
	cores []coreStruct
}
type numaStruct struct {
	ID 	int
	numL3Groups int
	l3Groups []l3GroupStruct
}
type socketStruct struct {
	ID int
	numNumas	int

	nodes []numaStruct
}
type topologyStruct struct {
		smtOn 		int
		numSockets int
		sockets []socketStruct
}


type containerEntry struct {
	Name   string
	CpuSet string
}

type podListEntry struct {
	Name       string //pod key
	Containers []containerEntry
}

type NodeState struct {
	Policy     string                       `yaml:"policyName"`
	DefaultSet string                       `yaml:"defaultCpuSet"`
	Entries    map[string]map[string]string `yaml:"entries"`

	Checksun int64 `yaml:"checksum"`
}

var pods []podListEntry
var contCnt = 0
var podCnt = 0

type cpuLookupStruct struct {
	socket 	   int
	node       int
	l3Group    int
	core       int
}

var cpuLookup [512]cpuLookupStruct


func main() {
	var topo *topologyStruct
    sockets    := flag.Int("s", 1, "Sockets")
	numaPer    := flag.Int("n", 1, "NUMA nodes per Sockets")
	l3PerNode  := flag.Int("l", 4, "L3Groups per NUMA node")
	coresPerL3 := flag.Int("c", 8, "Cores per L3Group")
	cpusPerCore  := flag.Int("t", 2, "CPUs per Core")
	stateFile  := flag.String("f", "./cpu_manager_state", "cpu_manager_state path/name")
    help := flag.Bool("help", false, "Help")

	flag.Parse()

    if *help {
        flag.PrintDefaults()
    } else {
		fmt.Printf("Sockets 		= %d\n", *sockets)
		fmt.Printf("Nodes per 		= %d\n", *numaPer)
		fmt.Printf("L3Groups per \t	= %d\n", *l3PerNode)
		fmt.Printf("Cores per 		= %d\n", *coresPerL3)
		fmt.Printf("CPUs per 		= %d\n", *cpusPerCore)
		fmt.Printf("file    		= %v\n", *stateFile)
		topo = buildTopology(*sockets, *numaPer, *l3PerNode, *coresPerL3, *cpusPerCore)
		fmt.Printf("smt %d\n", topo.smtOn)
		err := readStateFile(topo, *stateFile)
		if(err != nil){
			return
		}
		printMap(topo)





	}
}

func printMap(topo *topologyStruct){
	smtOn := topo.smtOn



	fmt.Printf("index: POD name\n")
	fmt.Printf("      - index: container name  (cpus )\n")
	for p := 0; p < len(pods); p++ {

		fmt.Printf("   %3d: %q\n", p, pods[p].Name)
		for c := 0; c < len(pods[p].Containers); c++ {
		fmt.Printf("     - %3d: %q  (%q)\n", c, pods[p].Containers[c].Name, pods[p].Containers[c].CpuSet)
		}
	}
	fmt.Printf("\nMap:\n")
	for s := 0; s < len(topo.sockets); s++{
		fmt.Printf("Socket[%v]\n",s)
		for n := 0; n < len(topo.sockets[s].nodes); n++ {
		fmt.Printf("Node[%v]\n", n)
			for l := 0; l < len(topo.sockets[s].nodes[n].l3Groups); l++ {
				fmt.Printf("L3group[%v]\n    CPU:  ", l)
				if smtOn == 2 {
					
					coresPerL3 := len(topo.sockets[s].nodes[n].l3Groups[l].cores)
					for i := 0; i < coresPerL3; i++ {
						fmt.Printf("%3d ", topo.sockets[s].nodes[n].l3Groups[l].cores[i].cpu[0])
					}
					fmt.Printf("\n    CPU:  ")
					for i := 0; i < coresPerL3; i++  {
						fmt.Printf("%3d ", topo.sockets[s].nodes[n].l3Groups[l].cores[i].cpu[1])
					}
					fmt.Printf("\n    type: ")
					for i := 0; i < coresPerL3; i++ {
						fmt.Printf("%3s ", topo.sockets[s].nodes[n].l3Groups[l].cores[i].containerType)
					}

					fmt.Printf("\n    Con:  ")
					for i := 0; i < coresPerL3; i++  {
						output := topo.sockets[s].nodes[n].l3Groups[l].cores[i].containerType
						if output != "" {
							if output == "S" {
								fmt.Printf("%3d ", topo.sockets[s].nodes[n].l3Groups[l].cores[i].containerIndex)
							} else {
								fmt.Printf("    " )
							}
	
						} else {
							fmt.Printf("    ")
						}
					}

					fmt.Printf("\n    POD:  ")
					for i := 0; i < coresPerL3; i++  {
						output := topo.sockets[s].nodes[n].l3Groups[l].cores[i].containerType
						if output != "" {
							if output == "S" {
								if topo.sockets[s].nodes[n].l3Groups[l].cores[i].splitL3Pod > 0 {
									fmt.Printf("\x1B[31m%3d\x1B[0m ", topo.sockets[s].nodes[n].l3Groups[l].cores[i].podIndex)
								}  else {
									fmt.Printf("%3d ", topo.sockets[s].nodes[n].l3Groups[l].cores[i].podIndex)
								}
								
							} else {
								fmt.Printf("    " )
							}
	
						} else {
							fmt.Printf("    ")
						}
					}
					fmt.Printf("\n")

				} else {

				}

			}

		}

	}

	fmt.Printf("Key:\n S Static POD\n")
}

func htopoSetState(topo *topologyStruct, socket int, node int, llc int, index int, state string, pod int, container int) {
	topo.sockets[socket].nodes[node].l3Groups[llc].cores[index].containerType = state
	topo.sockets[socket].nodes[node].l3Groups[llc].cores[index].podIndex = pod
	topo.sockets[socket].nodes[node].l3Groups[llc].cores[index].containerIndex = container
	
}

func topoSetCpuState(topo *topologyStruct, cpu int, state string, pod int, container int) {

	htopoSetState(topo, cpuLookup[cpu].socket, cpuLookup[cpu].node, cpuLookup[cpu].l3Group, cpuLookup[cpu].core, state, pod, container)
	//fmt.Printf("topoSetState(%v) %q\n", cpu, topoGetCpuState(cpu))
}

func setCpuSet(topo *topologyStruct, cpuSet string, pod int, container int) {

	s := strings.Split(cpuSet, ",")
	for i := 0; i < len(s); i++ {
		c := strings.Split(s[i], "-")
		if len(c) < 2 {
			//just on value
			cn, _ := strconv.Atoi(c[0])
			topoSetCpuState(topo, cn, "S", pod, container)

		} else {
			x, _ := strconv.Atoi(c[0]) // start
			y, _ := strconv.Atoi(c[1]) // end
			for ; x <= y; x++ {
				topoSetCpuState(topo, x, "S", pod, container)
			}
		}

	}
}


func readStateFile(topo *topologyStruct, file string) error {
	var nodeState NodeState

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v \n", err)
		return errors.New("yamlFile.Get err")
	}

	err = yaml.Unmarshal(yamlFile, &nodeState)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
		return errors.New("Unmarshal")
	}
	if nodeState.Policy == "static" {
		fmt.Printf("Policy \x1B[32m%q\x1B[0m\n", nodeState.Policy)
	} else {
		fmt.Printf("Policy \x1B[31m%q\x1B[0m\n", nodeState.Policy)
	}

	for pod, entry := range nodeState.Entries {
		//fmt.Printf("pod %q entry %q\n", pod, entry)
		newPod := podListEntry{}
		newPod.Name = pod
		for container, cpus := range entry {
			//fmt.Printf("   container %q cpus %q\n", container, cpus)
			newContainer := containerEntry{}
			newContainer.CpuSet = string(cpus)
			newContainer.Name = container
			newPod.Containers = append(newPod.Containers, newContainer)

			//fmt.Printf("\tcontainer %q cpuset %q\n", containers[count].Name, containers[count].CpuSet) ///"2-4,18-20"
			//fmt.Printf("\tcontainer %q cpuset %q\n", ekey, pod)                                        ///"2-4,18-20"
		}
		pods = append(pods, newPod)
	}
	for p := 0; p < len(pods); p++ {
		//fmt.Printf("    %3d   %q\n", p, pods[p].Name)
		for c := 0; c < len(pods[p].Containers); c++ {
			//fmt.Printf("      -%3d   %q\n", p, pods[p].Containers[c].Name)
			setCpuSet(topo, pods[p].Containers[c].CpuSet, p, c)
		}
		//fmt.Printf(" %d  %q : %q\n", i, containers[i].Name, containers[i].CpuSet)
	}
	for p := 0; p < len(pods); p++ {
		lastL3Id := -1									
		splitFound := 0
		for s := 0; s < len(topo.sockets); s++{
			for n := 0; n < len(topo.sockets[s].nodes); n++ {
				for l := 0; l < len(topo.sockets[s].nodes[n].l3Groups); l++ {
					for c := 0; c < len(topo.sockets[s].nodes[n].l3Groups[l].cores); c++ {
						if topo.sockets[s].nodes[n].l3Groups[l].cores[c].podIndex == p {
							if lastL3Id < 0 {
								lastL3Id = l
							} else {
								if lastL3Id != l{
									//lastL3Id = l
									splitFound++
								}
							}

						}
						
					}
				}
			}
		}
		if splitFound > 0 {
			for s := 0; s < len(topo.sockets); s++{
				for n := 0; n < len(topo.sockets[s].nodes); n++ {
					for l := 0; l < len(topo.sockets[s].nodes[n].l3Groups); l++ {
						for c := 0; c < len(topo.sockets[s].nodes[n].l3Groups[l].cores); c++ {
							if topo.sockets[s].nodes[n].l3Groups[l].cores[c].podIndex == p {
								topo.sockets[s].nodes[n].l3Groups[l].cores[c].splitL3Pod = 1
	
							}
							
						}
					}
				}
			}
		}
		
	}

	return nil

}

func buildTopology(socks int, numa int, l3 int, coresPerL3 int, smt int) *topologyStruct {
	topo := topologyStruct{}
	topo.smtOn = smt
	totalCpus := 0
	//calculate total cpus
	totalCpus = smt * coresPerL3 * l3 *numa *socks
	fmt.Printf("totalCpus = %d\n", totalCpus )
	topo.numSockets = socks
	core := 0
	for s := 0; s < socks; s++{
		newSocket := []socketStruct{socketStruct{}}
 		topo.sockets = append(topo.sockets, newSocket...)
		topo.sockets[s].ID = s
		topo.sockets[s].numNumas = numa
		for n := 0; n < numa; n++{
			newNode := []numaStruct{numaStruct{}}
 			topo.sockets[s].nodes = append(topo.sockets[s].nodes, newNode...)
			topo.sockets[s].nodes[n].ID = n
			topo.sockets[s].nodes[n].numL3Groups = l3
			for l := 0; l < l3; l++{
				newL3 := []l3GroupStruct{l3GroupStruct{}}
	 			topo.sockets[s].nodes[n].l3Groups = append(topo.sockets[s].nodes[n].l3Groups, newL3...)
				topo.sockets[s].nodes[n].l3Groups[l].ID = l
				topo.sockets[s].nodes[n].l3Groups[l].numCores = coresPerL3
				for c := 0; c < coresPerL3; c++{
					newCore := []coreStruct{coreStruct{}}
		 			topo.sockets[s].nodes[n].l3Groups[l].cores = append(topo.sockets[s].nodes[n].l3Groups[l].cores, newCore...)
					topo.sockets[s].nodes[n].l3Groups[l].cores[c].cpu[0] = core
					cpuLookup[core].socket = s
					cpuLookup[core].node = n
					cpuLookup[core].l3Group = l
					cpuLookup[core].core = c
					
					if smt == 2 {
						cpu := core + totalCpus/2
						topo.sockets[s].nodes[n].l3Groups[l].cores[c].cpu[1] = cpu
						cpuLookup[cpu].socket = s
						cpuLookup[cpu].node = n
						cpuLookup[cpu].l3Group = l
						cpuLookup[cpu].core = c

					} else {
						topo.sockets[s].nodes[n].l3Groups[l].cores[c].cpu[1] = -1
					}
					core++

				}
	
			}

		}
	}
/*
	for s := 0; s < len(topo.sockets); s++{
		fmt.Printf("Socket[%v]\n",s)
		for n := 0; n < len(topo.sockets[s].nodes); n++ {
			fmt.Printf("Node[%v]\n", n)
			for l := 0; l < len(topo.sockets[s].nodes[n].l3Groups); l++ {
				fmt.Printf("  l3Group[%v]\n", l)
				for c := 0; c < len(topo.sockets[s].nodes[n].l3Groups[l].cores); c++ {
					fmt.Printf("  Core[%v]  %v", l, topo.sockets[s].nodes[n].l3Groups[l].cores[c].cpu[0])
					if topo.smtOn == 2 {
						fmt.Printf(". %v\n", topo.sockets[s].nodes[n].l3Groups[l].cores[c].cpu[1])
					} else {
						fmt.Printf("\n")
					}
				}
			}
		}
	}

*/

	return &topo
}