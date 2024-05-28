package log

import (
	"sync"

	"go.uber.org/zap/zapcore"
)

type zapCoreWithCustomStack struct {
	core     zapcore.Core
	cloned   *zapCoreWithCustomStack
	stackMtx sync.Mutex
	stack    []byte
}

type CoreCloneObserver func(*zapCoreWithCustomStack)

// Sets stack and locks to ensure that nobody change stack till it is printed
func (c *zapCoreWithCustomStack) SetStackAndLock(stack []byte) {
	c.stackMtx.Lock()
	c.stack = stack
}

// Unlocks stack
func (c *zapCoreWithCustomStack) UnlockStack() {
	c.stackMtx.Unlock()
}

func (c *zapCoreWithCustomStack) Enabled(lvl zapcore.Level) bool {
	return c.core.Enabled(lvl)
}

// With adds structured context to the Core.
func (c *zapCoreWithCustomStack) With(fields []zapcore.Field) zapcore.Core {
	c.cloned = &zapCoreWithCustomStack{
		core:  c.core.With(fields),
		stack: nil,
	}
	return c.cloned
}

func (c *zapCoreWithCustomStack) GetCloned() *zapCoreWithCustomStack {
	return c.cloned
}

// Check determines whether the supplied Entry should be logged (using the
// embedded LevelEnabler and possibly some extra logic). If the entry
// should be logged, the Core adds itself to the CheckedEntry and returns
// the result.
//
// Callers must use Check before calling Write.
func (c *zapCoreWithCustomStack) Check(entry zapcore.Entry, checkedEntry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	chkEntry := c.core.Check(entry, checkedEntry)

	if chkEntry == nil {
		return nil
	}

	var newChkEntry *zapcore.CheckedEntry
	newChkEntry = newChkEntry.AddCore(chkEntry.Entry, c)

	return newChkEntry
}

// Write serializes the Entry and any Fields supplied at the log site and
// writes them to their destination.
//
// If called, Write should always log the Entry and Fields; it should not
// replicate the logic of Check.
func (c *zapCoreWithCustomStack) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	if c.stack != nil && len(c.stack) > 0 {
		if c.stack[0] != '\n' {
			entry.Stack = string(c.stack)
		} else {
			entry.Stack = string(c.stack[1:])
		}
		c.stack = nil
	}
	return c.core.Write(entry, fields)
}

// Sync flushes buffered logs (if any).
func (c *zapCoreWithCustomStack) Sync() error {
	return c.core.Sync()
}
