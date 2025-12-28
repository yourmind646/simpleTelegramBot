package fsm

var availableStates map[string][]string = map[string][]string{
	"MainMenu":  {"main"},
	"AdminMenu": {"main"},
	"AdminFiles": {"main", "inputFileKey", "inputNewFileValue"},
}
