package util

import "errors"

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
	FetchTags

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
fetch-tags                             Fetch tag(s) under a file
`

func Mux(cmd string) CmdType {
	switch cmd {
	case "add", "add-file":
		return AddFile
	case "rm", "remove", "remove-file":
		return RemoveFile
	case "search", "search-file":
		return SearchFile
	case "list":
		return ListFile
	case "add-tag":
		return AddTag
	case "remove-tag":
		return RemoveTag
	case "search-tag":
		return SearchTag
	case "list-tag":
		return ListTag
	case "associate", "a":
		return AssociateTag
	case "disassociate", "d":
		return DisassociateTag
	case "fetch":
		return FetchFile
	default:
		return UnknownCommand
	}
}

func Parse(t CmdType, args []string) error {
	switch t {
	case AddFile:
		return AddFileFlags.Parse(args)
	case RemoveFile:
		return RemoveFileFlags.Parse(args)
	case SearchFile:
		return SearchFileFlags.Parse(args)
	case ListFile:
		return ListFileFlags.Parse(args)
	case AddTag:
		return AddFileFlags.Parse(args)
	case RemoveTag:
		return RemoveTagFlags.Parse(args)
	case SearchTag:
		return SearchTagFlags.Parse(args)
	case ListTag:
		return ListTagFlags.Parse(args)
	case AssociateTag:
		return AssociateTagFlags.Parse(args)
	case DisassociateTag:
		return DisassociateTagFlags.Parse(args)
	case FetchFile:
		return FetchFileFlags.Parse(args)
	case FetchTags:
		return FetchTagsFlags.Parse(args)
	default:
		return errors.New("Tog Unknown command")
	}
}
