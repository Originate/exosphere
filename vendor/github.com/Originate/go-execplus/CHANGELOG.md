# 0.8.0 (2017-08-30)

* fix scanning of output chunks in stderr

# 0.7.0 (2017-08-30)

* fix possible race condition of stdout file descriptor when starting a command

# 0.6.0 (2017-08-23)

* fix possible race condition in `Wait*` funcs

# 0.5.0 (2017-08-17)

* Stop sending output chunks in their own goroutine
* Convert Mutex to RWMutex
