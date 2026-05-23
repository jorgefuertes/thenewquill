package condact

// OpCode is a single byte representing a condact in the bytecode.
type OpCode byte

const (
	NOOP OpCode = iota
	NOT         // NOT: negates the next condition

	// --- Conditions ---

	// arithmetic comparison
	EQ   // EQ <var|value|property> <var|value|property>: equal
	GT   // GT <var|value|property> <var|value|property>: greater than
	LT   // LT <var|value|property> <var|value|property>: less than
	GTE  // GTE <var|value|property> <var|value|property>: greater than or equal
	LTE  // LTE <var|value|property> <var|value|property>: less than or equal
	ZERO // ZERO <var|value|property>: is zero

	// probability
	CHANCE // CHANCE number[%]: succeeds number% of the time

	// entity location and state
	ISHERE       // ISHERE entity: entity is in same location as human
	ISAT         // ISAT entity loc: entity is at specific location
	ISCARRIED    // ISCARRIED item [who]: item is in who's inventory (default: human)
	ISWORN       // ISWORN item [who]: item is worn by who (default: human)
	ISPRESENT    // ISPRESENT entity: HERE or HAS or INSIDE open container (accessible to human)
	ISINSIDE     // ISINSIDE item item_container: item is inside container
	ISEMPTY      // ISEMPTY container|var: container has no items, or var is empty/zero
	ISOVERWEIGHT // ISOVERWEIGHT item [who]: adding item to who would exceed weight limit (default: human)

	// --- Actions ---

	// inventory management
	GET     // GET item [who]: pick up item (PRESENT) into human inventory
	DROP    // DROP item [who]: drop item at current location
	WEAR    // WEAR item [who]: equip item (CARRIED) as worn
	UNWEAR  // UNWEAR item [who]: unequip worn item back to inventory
	PUTIN   // PUTIN item container: put item inside container
	TAKEOUT // TAKEOUT item container: remove item from container
	GIVE    // GIVE item who: move item from human inventory to who's inventory

	// world manipulation
	MOVE    // MOVE entity loc: force move entity to location
	CREATE  // CREATE item [loc]: spawn item (default: current location)
	DESTROY // DESTROY item: remove item from game until CREATE is called again
	SWAP    // SWAP entity entity: exchange entities of same type
	GOTO    // GOTO loc: move human to location (fires location events)

	// variables
	SET    // SET var value: assign value to variable or property
	INC    // INC var [n]: increment by n (default: 1)
	DEC    // DEC var [n]: decrement by n (default: 1)
	RANDOM // RANDOM var min max: assign random value in [min, max]

	// terminal
	WINDOW // WINDOW id: open or select window
	AT     // AT row col: set cursor position
	PRINT  // PRINT msg|prop: print message label or property value
	CURSOR // CURSOR on|off: show or hide cursor
	CLEAR  // CLEAR [window]: clear screen or window
	BG     // BG color: set background color
	FG     // FG color: set foreground color

	// media
	PLAY    // PLAY blob: play audio
	PICTURE // PICTURE blob: display image

	// input / timing
	INPUT  // INPUT var prompt: read text input into variable
	WAIT   // WAIT ms: pause execution for milliseconds
	ANYKEY // ANYKEY: pause until any key is pressed

	// persistence
	SAVE // SAVE [slot]: save game state
	LOAD // LOAD [slot]: load game state

	// game control
	PARSE   // PARSE str: re-parse string as player input
	ENGAGE  // ENGAGE npc: start NPC dialogue dispatch
	RESTART // RESTART: restart adventure from the beginning
	QUIT    // QUIT: exit game

	// process control
	DONE // end table cycle; continue with next SUBSL or SL
	OK   // print ok message; end table cycle; continue with next SUBSL or SL
	END  // end current table; continue with next table in cycle
	NOOK // print nook message; discard all remaining SLs and SUBSLs
)
