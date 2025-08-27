package util

type CmdType uint8

const (
	UnknownCommand CmdType = iota

	AddFile
	RemoveFile
	SearchFile
	ListFile

	AddTag
	RemoveTag
	SearchTag
	ListTag

	AssociateTag
	DisassociateTag
	FetchFile

	// TODO: Implement these operation
	// JsonDump
	// LoadJson
)

const CommandDescription string = `Command Options:
add, add-file                          Add File to be managed (default)
rm, remove, remove-file                Remove File from being managed
s, search                              Search managed files
list                                   List all managed files

add-tag                                Create new tag
remove-tag                             Remove tag
search-tag                             Search among tag(s)
list-tag                               List all tags

a, associate                           Associate file(s) with tag
d, disassociate                        Disassociate file(s) with tag
fetch                                  Fetch file(s) under a tag
`

func Mux(cmd string) CmdType {
	switch cmd {
	case "add", "add-file":
		return AddFile
	case "rm", "remove", "remove-file":
		return RemoveFile
	case "search", "search-file":
		return SearchFile
	case "add-tag":
		return AddTag
	case "remove-tag":
		return RemoveTag
	case "search-tag":
		return SearchTag
	case "associate", "a":
		return AssociateTag
	case "disassociate", "d":
		return DisassociateTag
	default:
		return UnknownCommand
	}
}
