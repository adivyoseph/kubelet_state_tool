# kubelet_state_tool

purpose to read a kubelet cpu_manager_state found in /var/lib/kubelet/cpu_manager_state and map the PODs/containers for fixed/static to AMD server topology


# Usage:

./kubelet_state_tool -help

  -s int      Sockets (default 1)  

  -n int      NUMA nodes per Sockets (default 1)

  -l int      L3Groups per NUMA node (default 4)

  -c int      Cores per L3Group (default 8) 
  
  -t int      CPUs per Core (default 2)

  -f string      cpu_manager_state path/name (default "./cpu_manager_state")

  -help       Help

 
 
 # example

` 
index: POD name
      - index: container name  (cpus )
     0: "a9ab02d4-c08d-42b2-b5d0-474c43d1a871"
     -   0: "eric-pc-up-data-plane"  ("13-19,109-115")
     1: "f66add72-58f1-4676-a437-b13fc2d2f83b"
     -   0: "mon"  ("4-5,100-101")
     -   1: "chown-container-data-dir"  ("4-5,100-101")
     -   2: "init-mon-fs"  ("4-5,100-101")
     2: "ef4ff370-0cf7-42c3-9952-32c033257ce7"
     -   0: "eric-pc-up-data-plane"  ("6-12,102-108")
     3: "8c0ebe9f-131c-4af6-aef4-48fd223cfdce"
     -   0: "eric-pc-up-data-plane"  ("27-33,123-129")
     4: "bf4e1603-21ec-467f-b54e-1af32e93ba01"
     -   0: "eric-pc-up-data-plane"  ("72-78,168-174")
     5: "47525978-b2ca-49be-a8f2-06e54c927e11"
     -   0: "chown-container-data-dir"  ("49-51,145-147")
     -   1: "activate"  ("49-51,145-147")
     -   2: "osd"  ("49-51,145-147")
     6: "4656f088-0379-49af-a414-52d5d6002a47"
     -   0: "eric-pc-up-data-plane"  ("41-47,137-143")
     7: "7776d501-650f-4501-93d4-509c23430cd0"
     -   0: "eric-pc-up-data-plane"  ("20-26,116-122")
     8: "048ce3dc-6bed-42ad-9cb6-78d09337faed"
     -   0: "eric-pc-up-data-plane"  ("65-71,161-167")
     9: "17f878fd-fb73-42d5-9d1e-3ffdacf2ab07"
     -   0: "eric-pc-up-data-plane"  ("58-64,154-160")
    10: "3367b904-b4f8-4966-a311-2b3e146813f3"
     -   0: "eric-pc-up-data-plane"  ("48,52-57,144,148-153")
    11: "3ee6b330-0b3d-49d0-a9f1-400708b17f54"
     -   0: "eric-pc-up-data-plane"  ("34-40,130-136")

Map:
Socket[0]
Node[0]
L3group[0]
    CPU:    0   1   2   3   4   5   6   7
    CPU:   96  97  98  99 100 101 102 103
    type:                   S   S   S   S
    Con:                    2   2   0   0
    POD:                    1   1   2   2
L3group[1]
    CPU:    8   9  10  11  12  13  14  15
    CPU:  104 105 106 107 108 109 110 111
    type:   S   S   S   S   S   S   S   S
    Con:    0   0   0   0   0   0   0   0
    POD:    2   2   2   2   2   0   0   0
L3group[2]
    CPU:   16  17  18  19  20  21  22  23
    CPU:  112 113 114 115 116 117 118 119
    type:   S   S   S   S   S   S   S   S
    Con:    0   0   0   0   0   0   0   0
    POD:    0   0   0   0   7   7   7   7
L3group[3]
    CPU:   24  25  26  27  28  29  30  31
    CPU:  120 121 122 123 124 125 126 127
    type:   S   S   S   S   S   S   S   S
    Con:    0   0   0   0   0   0   0   0
    POD:    7   7   7   3   3   3   3   3
L3group[4]
    CPU:   32  33  34  35  36  37  38  39
    CPU:  128 129 130 131 132 133 134 135
    type:   S   S   S   S   S   S   S   S
    Con:    0   0   0   0   0   0   0   0
    POD:    3   3  11  11  11  11  11  11
L3group[5]
    CPU:   40  41  42  43  44  45  46  47
    CPU:  136 137 138 139 140 141 142 143
    type:   S   S   S   S   S   S   S   S
    Con:    0   0   0   0   0   0   0   0
    POD:   11   6   6   6   6   6   6   6
L3group[6]
    CPU:   48  49  50  51  52  53  54  55
    CPU:  144 145 146 147 148 149 150 151
    type:   S   S   S   S   S   S   S   S
    Con:    0   2   2   2   0   0   0   0
    POD:   10   5   5   5  10  10  10  10
L3group[7]
    CPU:   56  57  58  59  60  61  62  63
    CPU:  152 153 154 155 156 157 158 159
    type:   S   S   S   S   S   S   S   S
    Con:    0   0   0   0   0   0   0   0
    POD:   10  10   9   9   9   9   9   9
L3group[8]
    CPU:   64  65  66  67  68  69  70  71
    CPU:  160 161 162 163 164 165 166 167
    type:   S   S   S   S   S   S   S   S
    Con:    0   0   0   0   0   0   0   0
    POD:    9   8   8   8   8   8   8   8
L3group[9]
    CPU:   72  73  74  75  76  77  78  79
    CPU:  168 169 170 171 172 173 174 175
    type:   S   S   S   S   S   S   S
    Con:    0   0   0   0   0   0   0
    POD:    4   4   4   4   4   4   4
L3group[10]
    CPU:   80  81  82  83  84  85  86  87
    CPU:  176 177 178 179 180 181 182 183
    type:
    Con:
    POD:
L3group[11]
    CPU:   88  89  90  91  92  93  94  95
    CPU:  184 185 186 187 188 189 190 191
    type:
    Con:
    POD:
Key:
 S Static POD

`

