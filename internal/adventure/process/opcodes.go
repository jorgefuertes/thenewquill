package process

// OpCode is a single byte representing a condact in the bytecode.
type OpCode byte

const (
	NOOP OpCode = 0
	NOT  OpCode = 1 // NOT: negates the next condition

	// --- Conditions ---

	// arithmetic comparison
	EQ   OpCode = 2 // EQ <var|value|property> <var|value|property>: equal
	GT   OpCode = 3 // GT <var|value|property> <var|value|property>: greater than
	LT   OpCode = 4 // LT <var|value|property> <var|value|property>: less than
	GTE  OpCode = 5 // GTE <var|value|property> <var|value|property>: greater than or equal
	LTE  OpCode = 6 // LTE <var|value|property> <var|value|property>: less than or equal
	ZERO OpCode = 7 // ZERO <var|value|property>: is zero

	// probability
	CHANCE OpCode = 8 // CHANCE number[%]: succeeds number% of the time

	// entity location and state
	ISHERE       OpCode = 9  // ISHERE entity: entity is in same location as human
	ISAT         OpCode = 10 // ISAT entity loc: entity is at specific location
	ISCARRIED    OpCode = 11 // ISCARRIED item [who]: item is in who's inventory (default: human)
	ISWORN       OpCode = 12 // ISWORN item [who]: item is worn by who (default: human)
	ISPRESENT    OpCode = 13 // ISPRESENT entity: HERE or HAS or INSIDE open container (accessible to human)
	ISINSIDE     OpCode = 14 // ISINSIDE item item_container: item is inside container
	ISEMPTY      OpCode = 15 // ISEMPTY container|var: container has no items, or var is empty/zero
	ISOVERWEIGHT OpCode = 16 // ISOVERWEIGHT item [who]: adding item to who would exceed weight limit (default: human)

	// --- Actions ---

	// inventory management
	GET     OpCode = 50 // GET item [who]: pick up item (PRESENT) into human inventory
	DROP    OpCode = 51 // DROP item [who]: drop item at current location
	WEAR    OpCode = 52 // WEAR item [who]: equip item (CARRIED) as worn
	UNWEAR  OpCode = 53 // UNWEAR item [who]: unequip worn item back to inventory
	PUTIN   OpCode = 54 // PUTIN item container: put item inside container
	TAKEOUT OpCode = 55 // TAKEOUT item container: remove item from container
	GIVE    OpCode = 56 // GIVE item who: move item from human inventory to who's inventory

	// world manipulation
	MOVE    OpCode = 60 // MOVE entity loc: force move entity to location
	CREATE  OpCode = 61 // CREATE item [loc]: spawn item (default: current location)
	DESTROY OpCode = 62 // DESTROY item: remove item from game until CREATE is called again
	SWAP    OpCode = 63 // SWAP entity entity: exchange entities of same type
	GOTO    OpCode = 64 // GOTO loc: move human to location (fires location events)

	// variables
	SET    OpCode = 70 // SET var value: assign value to variable or property
	INC    OpCode = 71 // INC var [n]: increment by n (default: 1)
	DEC    OpCode = 72 // DEC var [n]: decrement by n (default: 1)
	RANDOM OpCode = 73 // RANDOM var min max: assign random value in [min, max]

	// terminal
	WINDOW OpCode = 80 // WINDOW id: open or select window
	CLOSE  OpCode = 81 // CLOSE id: close a window
	AT     OpCode = 82 // AT row col: set cursor position
	PRINT  OpCode = 83 // PRINT msg|prop: print message label or property value
	CURSOR OpCode = 84 // CURSOR on|off: show or hide cursor
	CLEAR  OpCode = 85 // CLEAR [window]: clear screen or window
	BG     OpCode = 86 // BG color: set background color
	FG     OpCode = 87 // FG color: set foreground color

	// gfx & sfx
	PLAY    OpCode = 90 // PLAY blob: play audio
	PICTURE OpCode = 91 // PICTURE blob: display image

	// input / timing
	INPUT  OpCode = 101 // INPUT var prompt: read text input into variable
	WAIT   OpCode = 102 // WAIT ms: pause execution for milliseconds
	ANYKEY OpCode = 103 // ANYKEY: pause until any key is pressed

	// persistence
	SAVE OpCode = 110 // SAVE [slot]: save game state
	LOAD OpCode = 111 // LOAD [slot]: load game state

	// game control
	PARSE   OpCode = 120 // PARSE str: re-parse string as player input
	ENGAGE  OpCode = 121 // ENGAGE npc: start NPC dialogue dispatch
	RESTART OpCode = 122 // RESTART: restart adventure from the beginning
	QUIT    OpCode = 123 // QUIT: exit game

	// process control
	DONE OpCode = 124 // end table cycle; continue with next SUBSL or SL
	OK   OpCode = 125 // print ok message; end table cycle; continue with next SUBSL or SL
	END  OpCode = 126 // end current table; continue with next table in cycle
	NOOK OpCode = 127 // print nook message; discard all remaining SLs and SUBSLs
)
