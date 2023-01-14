package DeepMerge

import (
	"fmt"
)

type Config struct {
	// PreserveUnmergeables set to true to skip any unmergeable elements from source
	PreserveUnmergeables bool

	// KnockoutPrefix set to string value to signify prefix which deletes elements from existing element
	KnockoutPrefix *string

	// OverwriteArrays set to true if you want to avoid merging arrays
	OverwriteArrays bool

	// ExtendExistingArrays set to true to add src elements to existing array rather than overwriting
	ExtendExistingArrays bool

	// SortMergedArrays set to true to sort all arrays that are merged together
	SortMergedArrays bool

	// UnpackArrays set to string value to run "Array::join" then "String::split" against all arrays
	UnpackArrays *string

	// MergeHashArrays set to true to merge hashes within arrays
	MergeHashArrays bool

	// KeepArrayDuplicates set to true to preserve duplicate array entries
	KeepArrayDuplicates bool

	// MergeNilValues set to true to merge empty source values rather than skipping them (the default)
	MergeNilValues bool

	// Debug set to true to get console output of merge process for debugging
	Debug bool

	// DebugIndent set to customize indentation level of debug output
	DebugIndent string
}

func NewConfig() *Config {
	return &Config{
		PreserveUnmergeables: false,
		KnockoutPrefix:       nil,
		OverwriteArrays:      false,
		SortMergedArrays:     false,
		UnpackArrays:         nil,
		ExtendExistingArrays: false,
		MergeHashArrays:      false,
		MergeNilValues:       false,
		KeepArrayDuplicates:  false,
		Debug:                false,
		DebugIndent:          "",
	}
}

func NewConfigDeeperMergeKO() *Config {
	return NewConfig().WithKnockout("--").OverwriteUnmergeable()
}

func NewConfigDeeperMergeBang() *Config {
	return NewConfig().OverwriteUnmergeable()
}

func NewConfigDeeperMerge() *Config {
	return NewConfig().PreserveUnmergeable()
}

func (c *Config) WithDebug(d bool) *Config {
	c.Debug = d
	return c
}

func (c *Config) EnableDebug() *Config {
	return c.WithDebug(true)
}

func (c *Config) PreserveUnmergeable() *Config {
	c.PreserveUnmergeables = true
	return c
}

func (c *Config) OverwriteUnmergeable() *Config {
	c.PreserveUnmergeables = false
	return c
}

func (c *Config) WithKnockout(prefix string) *Config {
	c.KnockoutPrefix = &prefix
	return c
}

func (c *Config) UnpackArraysWith(sep string) *Config {
	c.UnpackArrays = &sep
	return c
}

// writeDebug conditionally writes a formatted message only when o.Debug is true
func (c *Config) writeDebug(f string, a ...interface{}) {
	if c.Debug {
		fmt.Println(c.DebugIndent + fmt.Sprintf(f, a...))
	}
}

func (c *Config) copyWithIncreasedDebugIndent() *Config {
	cc := *c
	cc.DebugIndent = "  " + c.DebugIndent
	return &cc
}
