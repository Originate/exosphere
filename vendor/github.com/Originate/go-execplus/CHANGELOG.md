# 0.9.2 (2017-10-05)

* prevent send on closed channel

# 0.9.1 (2017-09-26)

* fix random failure in `cmdPlus.Kill()`

# 0.9.0 (2017-09-26)

* replace `SetEnv` with `AppendEnv` which retains the current process environment

# 0.8.1 (2017-09-18)

* fix possible race condition when sending the initial chunk

# 0.8.0 (2017-08-30)

* fix scanning of output chunks in stderr

# 0.7.0 (2017-08-30)

* fix possible race condition of stdout file descriptor when starting a command

# 0.6.0 (2017-08-23)

* fix possible race condition in `Wait*` funcs

# 0.5.0 (2017-08-17)

* Stop sending output chunks in their own goroutine
* Convert Mutex to RWMutex
