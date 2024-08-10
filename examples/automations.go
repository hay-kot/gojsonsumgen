// Package test contains test files used for the gosumtype cli tool
package examples

// ActionType is a sum type discriminator for ActionDef
//
// gosumtype: ActionDef
//
//	opt:tag:  type
//	opt:name: Type
//
//	'action-shell'  : ActionShell
//	'action-python' : ActionPython
//	'action-http'   : ActionHTTP
type ActionType string

type ActionShell struct {
	Type   ActionType `json:"type"`
	ValueA string     `json:"value_a"`
}

type ActionPython struct {
	Type   ActionType `json:"type"`
	ValueB int        `json:"value_b"`
}

type ActionHTTP struct {
	Type   ActionType `json:"type"`
	ValueC bool       `json:"value_c"`
}

type (
	// ShellType is a sum type discriminator for ShellDef
	//
	// gosumtype: ShellDef
	//
	//  opt:tag:  type
	//  opt:name: Type
	//
	//  'shell-bash' : ShellBash
	//  'shell-zsh'  : ShellZsh
	//  'shell-fish' : ShellFish
	ShellType string

	ShellBash struct {
		Type  ShellType `json:"type"`
		Value string    `json:"value"`
	}

	ShellZsh struct {
		Type  ShellType `json:"type"`
		Value string    `json:"value"`
	}

	ShellFish struct {
		Type  ShellType `json:"type"`
		Value string    `json:"value"`
	}
)
