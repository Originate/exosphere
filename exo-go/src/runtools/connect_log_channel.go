package runtools

import execplus "github.com/Originate/go-execplus"

// ConnectLogChannel connects a log channel that wants to receive only each new output
func ConnectLogChannel(cmdPlus *execplus.CmdPlus, logChannel chan string) {
	outputChannel, _ := cmdPlus.GetOutputChannel()
	go func() {
		for {
			outputChunk := <-outputChannel
			if outputChunk.Chunk != "" {
				logChannel <- outputChunk.Chunk
			}
		}
	}()
}
