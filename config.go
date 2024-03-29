package kairos

// Config 是一个结构体，包含一个 Callback 类型的字段
// Config is a struct that contains a field of type Callback
type Config struct {
	// callback 是一个 Callback 类型的字段，用于设置任务的回调函数
	// callback is a field of type Callback, used to set the callback functions for the task
	callback Callback
}

// NewConfig 是一个函数，用于创建一个新的 Config 实例
// NewConfig is a function used to create a new instance of Config
func NewConfig() *Config {
	// 返回一个新的 Config 实例，其中 callback 字段被设置为一个新的空任务回调
	// Return a new instance of Config, where the callback field is set to a new empty task callback
	return &Config{
		callback: NewEmptyTaskCallback(),
	}
}

// DefaultConfig 是一个函数，用于获取默认的 Config 实例
// DefaultConfig is a function used to get the default instance of Config
func DefaultConfig() *Config {
	// 返回一个新的 Config 实例
	// Return a new instance of Config
	return NewConfig()
}

// WithCallback 是 Config 的一个方法，用于设置 Config 的 callback 字段
// WithCallback is a method of Config, used to set the callback field of Config
func (c *Config) WithCallback(callback Callback) {
	// 设置 Config 的 callback 字段为传入的 callback 参数
	// Set the callback field of Config to the passed-in callback parameter
	c.callback = callback
}

// isConfigValid 是一个函数，用于检查 Config 实例是否有效
// isConfigValid is a function used to check if the instance of Config is valid
func isConfigValid(conf *Config) *Config {
	// 如果 conf 不为 nil
	// If conf is not nil
	if conf != nil {
		// 如果 conf 的 callback 字段为 nil
		// If the callback field of conf is nil
		if conf.callback == nil {
			// 设置 conf 的 callback 字段为一个新的空任务回调
			// Set the callback field of conf to a new empty task callback
			conf.callback = NewEmptyTaskCallback()
		}
	} else {
		// 如果 conf 为 nil，设置 conf 为默认的 Config 实例
		// If conf is nil, set conf to the default instance of Config
		conf = DefaultConfig()
	}

	// 返回 conf
	// Return conf
	return conf
}
