package business

func Main(c chan *[]byte) {
	// todo : kafka config
	consumer(c, nil)
}
