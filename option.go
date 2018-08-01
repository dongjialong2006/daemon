package daemon

type option interface {
	Get(key string) interface{}
}

type element struct {
	name  string
	value interface{}
}

func (o element) Get(key string) interface{} {
	if o.name == key {
		return o.value
	}

	return nil
}

func new(name string, value interface{}) *element {
	return &element{
		name:  name,
		value: value,
	}
}

const (
	process_num  = "process_num"
	process_envs = "process_envs"
	process_args = "process_args"
)

func WithProcessNum(num int32) option {
	return new(process_num, num)
}

func WithEnvs(envs []string) option {
	return new(process_envs, envs)
}

func WithArgs(args []string) option {
	return new(process_args, args)
}
