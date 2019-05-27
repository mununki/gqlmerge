package command

func Usage() string {
	return helpMsg + flagIndentMsg
}

const helpMsg = `👋 'gqlmerge' is the tool to merge & stitch GraphQL files and generate a GraphQL schema
Author : Woonki Moon <woonki.moon@gmail.com>

Usage:	gqlmerge [FLAG ...] [PATH ...] [OUTPUT]

e.g.

	gqlmerge ./schema schema.graphql

Flags:

	-v	: check the version
	-h	: help
`

const flagIndentMsg = `
	-indent	: (default=2s) defines the padding

	It follows the next pattern: indent={n}{i},

		* n - amount of indents
		* i - indent ("t" for tabs and "s" for spaces)

	If "n" is not stated 1 will be used, 
	so "--indent=1t" is equal to "--indent=t"`
