package util

import "flag"

type AddFileCommand struct {
	Path string
}

type RemoveFileCommand struct {
	Path string
}

type SearchFileCommand struct {
	PathGlob string
}

type ListFileCommand struct{}

type AddTagCommand struct {
	TagName string
	TagDesc string
}

type RemoveTagCommand struct {
	TagName string
}

type SearchTagCommand struct {
	TagGlob string
}

type ListTagCommand struct{}

type AssociateTagCommand struct {
	Path    []string
	TagName string
}

type DisassociateTagCommand struct {
	Path    string
	TagName string
}

type FetchFileCommand struct {
	TagName string
}

var (
	AddFileFlags         flag.FlagSet
	RemoveFileFlags      flag.FlagSet
	SearchFileFlags      flag.FlagSet
	ListFileFlags        flag.FlagSet
	AddTagFlags          flag.FlagSet
	RemoveTagFlags       flag.FlagSet
	SearchTagFlags       flag.FlagSet
	ListTagFlags         flag.FlagSet
	AssociateTagFlags    flag.FlagSet
	DisassociateTagFlags flag.FlagSet
	FetchFileFlags       flag.FlagSet
	FetchTagsFlags       flag.FlagSet
)
