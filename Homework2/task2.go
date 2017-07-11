// task2
package main

func ConcurrentRetryExecutor(tasks []func() string, concurrentLimit int, retryLimit int) <-chan struct {
	index  int
	result string
} {

	TheChannel := make(chan struct {
		index  int
		result string
	})

	Used := make(chan struct{}, concurrentLimit)
	done := make(chan struct{})

	go func() {
		for index, fnc := range tasks {

			Used <- struct{}{}

			go func(limit int, ind int, fnc func() string) {

				for times := limit; times > 0; times-- {

					res := fnc()

					TheChannel <- struct {
						index  int
						result string
					}{ind, res}

					if res != "" {

						<-Used
						done <- struct{}{}
						return

					}
				}

				<-Used
				done <- struct{}{}

			}(retryLimit, index, fnc)

		}
	}()

	go func() {

		for range tasks {

			<-done
		}

		close(TheChannel)

	}()

	return TheChannel
}
