package main

// Channel is a channel.
type Channel struct {
	ch chan Value
}

// NewChannel news a channel.
func NewChannel(bufSize Value) *Channel {
	if !bufSize.isNumber() || bufSize.number() < 0 {
		panic(NewTypeError("Channel: buffer size must be a number greater than zero"))
	}
	return &Channel{
		ch: make(chan Value, bufSize.number()),
	}
}

// GetKey implements KeyGetter.
func (c *Channel) GetKey(key string) Value {
	if fn, ok := _channelMethods[key]; ok {
		return ValueFromBuiltin(c, key, fn)
	}
	return ValueFromNil()
}

// SetKey implements KeyIndexer.
func (c *Channel) SetKey(key string, val Value) {
	panic(NewNotAssignableError(ValueFromObject(c)))
}

// Read reads a value from channel.
func (c *Channel) Read() Value {
	return <-c.ch
}

// Write writes a value into channel.
func (c *Channel) Write(value Value) {
	c.ch <- value
}

// Close closes the channel.
func (c *Channel) Close() {
	close(c.ch)
}

var _channelMethods map[string]BuiltinFunction

func init() {
	_channelMethods = map[string]BuiltinFunction{
		"read":  _channelRead,
		"write": _channelWrite,
		"close": _channelClose,
	}
}

func _channelRead(this interface{}, ctx *Context, args *Values) Value {
	channel := this.(*Channel)
	return channel.Read()
}

func _channelWrite(this interface{}, ctx *Context, args *Values) Value {
	channel := this.(*Channel)
	for _, arg := range args.values {
		channel.Write(arg)
	}
	return ValueFromNil()
}

func _channelClose(this interface{}, ctx *Context, args *Values) Value {
	channel := this.(*Channel)
	channel.Close()
	return ValueFromNil()
}
