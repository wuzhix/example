package trap

import "time"

/**
 * 迭代器中使用线程的陷阱
 */
func GoParam() {
	println("直接使用参数")
	for i := 0; i < 5; i++ {
		go func() {
			println(i)
		}()
	}
	time.Sleep(3 * time.Second)
	println("传入参数")
	for i := 0; i < 5; i++ {
		go func(index int) {
			println(index)
		}(i)
	}
	time.Sleep(3 * time.Second)
}
