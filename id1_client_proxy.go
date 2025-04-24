package id1_client

type Id1ClientProxy struct {
	client         Id1Client
	preprocessors  []func(cmd *Command) error
	postprocessors []func(data []byte, err error) ([]byte, error)
}

func NewId1ClientProxy(client Id1Client) Id1ClientProxy {
	return Id1ClientProxy{
		client:         client,
		preprocessors:  []func(cmd *Command) error{},
		postprocessors: []func(data []byte, err error) ([]byte, error){},
	}
}

func (t *Id1ClientProxy) AddPreprocessor(fn func(cmd *Command) error) {
	t.preprocessors = append(t.preprocessors, fn)
}

func (t *Id1ClientProxy) AddPostprocessor(fn func(data []byte, err error) ([]byte, error)) {
	t.postprocessors = append(t.postprocessors, fn)
}

func (t *Id1ClientProxy) preprocess(cmd *Command) error {
	if cmd.Args == nil {
		cmd.Args = map[string]string{}
	}
	for _, preprocessor := range t.preprocessors {
		if err := preprocessor(cmd); err != nil {
			return err
		}
	}
	return nil
}

func (t *Id1ClientProxy) postprocess(data []byte, err error) ([]byte, error) {
	for _, postprocessor := range t.postprocessors {
		if err != nil {
			return []byte{}, err
		}
		data, err = postprocessor(data, err)
	}
	return data, err
}

func (t *Id1ClientProxy) Authenticate(id string, privateKey string) error {
	return t.client.Authenticate(id, privateKey)
}

func (t *Id1ClientProxy) Connect() (chan bool, error) {
	return t.client.Connect()
}

func (t *Id1ClientProxy) Close() {
	t.client.Close()
}

func (t *Id1ClientProxy) AddListener(listener func(cmd Command), listenerId string) string {
	return t.client.AddListener(listener, listenerId)
}

func (t *Id1ClientProxy) RemoveListener(listenerId string) {
	t.client.RemoveListener(listenerId)
}

func (t *Id1ClientProxy) Send(cmd Command) error {
	if err := t.preprocess(&cmd); err != nil {
		return err
	} else {
		return t.client.Send(cmd)
	}
}

func (t *Id1ClientProxy) Exec(cmd Command) ([]byte, error) {
	if err := t.preprocess(&cmd); err != nil {
		return []byte{}, err
	} else {
		data, err := t.client.Exec(cmd)
		return t.postprocess(data, err)
	}
}

func (t *Id1ClientProxy) Get(key Id1Key) ([]byte, error) {
	cmd := &Command{Op: Get, Key: key}
	if err := t.preprocess(cmd); err != nil {
		return []byte{}, err
	} else {
		data, err := t.client.Get(cmd.Key)
		return t.postprocess(data, err)
	}
}

func (t *Id1ClientProxy) Del(key Id1Key) error {
	cmd := &Command{Op: Del, Key: key}
	if err := t.preprocess(cmd); err != nil {
		return err
	} else {
		return t.client.Del(cmd.Key)
	}
}
func (t *Id1ClientProxy) Set(key Id1Key, data []byte) error {
	cmd := &Command{Op: Set, Key: key, Data: data}
	if err := t.preprocess(cmd); err != nil {
		return err
	} else {
		return t.client.Set(cmd.Key, cmd.Data)
	}
}

func (t *Id1ClientProxy) Add(key Id1Key, data []byte) error {
	cmd := &Command{Op: Add, Key: key, Data: data}
	if err := t.preprocess(cmd); err != nil {
		return err
	} else {
		return t.client.Add(cmd.Key, cmd.Data)
	}
}

func (t *Id1ClientProxy) Mov(key, tgtKey Id1Key) error {
	cmd := &Command{Op: Mov, Key: key, Data: []byte(tgtKey.String())}
	if err := t.preprocess(cmd); err != nil {
		return err
	} else {
		return t.client.Mov(cmd.Key, K(string(cmd.Data)))
	}
}

func (t *Id1ClientProxy) List(key Id1Key, options ListOptions) (map[string][]byte, error) {
	cmd := &Command{Op: Get, Key: key, Args: options.Map()}
	if err := t.preprocess(cmd); err != nil {
		return map[string][]byte{}, err
	} else {
		options := &ListOptions{}
		options.Parse(cmd.Args)
		if list, err := t.client.List(cmd.Key, *options); err != nil {
			return map[string][]byte{}, err
		} else {
			processedList := map[string][]byte{}
			for k, v := range list {
				if data, err := t.postprocess(v, nil); err == nil {
					processedList[k] = data
				}
			}
			return processedList, nil
		}
	}
}
