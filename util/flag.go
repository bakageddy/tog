package util

import "flag"

type AddFileCommand struct {
	Path string
}

func (a AddFileCommand) Setup() {
	AddFileFlags.StringVar(&a.Path, "file", ".", "Add File to be Managed")
}

type RemoveFileCommand struct {
	Path string
}

func (a RemoveFileCommand) Setup() {
	RemoveFileFlags.StringVar(&a.Path, "file", ".", "Remove File from being Managed")
}

type SearchFileCommand struct {
	PathGlob string
}

func (a SearchFileCommand) Setup() {
	SearchFileFlags.StringVar(&a.PathGlob, "search", "", "Search file(s)")
}

type ListFileCommand struct {
	Execute bool
}

func (a ListFileCommand) Setup() {
	ListFileFlags.BoolVar(&a.Execute, "list", true, "List all managed file(s)")
}

type AddTagCommand struct {
	TagName string
	TagDesc string
}

func (a AddTagCommand) Setup() {
	AddTagFlags.StringVar(&a.TagName, "name", "[default]", "Name of the tag")
	AddTagFlags.StringVar(&a.TagDesc, "desc", "None", "Description of the tag")
}

type RemoveTagCommand struct {
	TagName string
}

func (a RemoveTagCommand) Setup() {
	RemoveTagFlags.StringVar(&a.TagName, "name", "", "Name of the tag to remove")
}

type SearchTagCommand struct {
	TagGlob string
}

func (a SearchTagCommand) Setup() {
	SearchTagFlags.StringVar(&a.TagGlob, "search", "", "Search Tag(s)")
}

type ListTagCommand struct {
	Execute bool
}

func (a ListTagCommand) Setup() {
	ListTagFlags.BoolVar(&a.Execute, "list", true, "List Tag(s)")
}

type ManagedPaths []string

type AssociateTagCommand struct {
	Path    ManagedPaths
	TagName string
}

func (a AssociateTagCommand) Setup() {
	AssociateTagFlags
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
