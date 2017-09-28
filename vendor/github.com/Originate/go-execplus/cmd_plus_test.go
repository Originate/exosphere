package execplus_test

import (
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	execplus "github.com/Originate/go-execplus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Process", func() {
	It("returns no errors when the process succeeds", func() {
		cmdPlus := execplus.NewCmdPlus("./test_executables/passing")
		err := cmdPlus.Run()
		Expect(err).To(BeNil())
	})

	It("returns errors when the process fails", func() {
		cmdPlus := execplus.NewCmdPlus("./test_executables/failing")
		err := cmdPlus.Run()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("exit status 1"))
	})

	It("captures output", func() {
		cmdPlus := execplus.NewCmdPlus("./test_executables/passing")
		err := cmdPlus.Run()
		Expect(err).To(BeNil())
		Expect(cmdPlus.GetOutput()).To(Equal("output"))
	})

	It("allows settings of the current working directory", func() {
		cwd, err := os.Getwd()
		customDir := path.Join(cwd, "test_executables")
		Expect(err).To(BeNil())
		cmdPlus := execplus.NewCmdPlus("./print_cwd")
		cmdPlus.SetDir(customDir)
		err = cmdPlus.Run()
		Expect(err).To(BeNil())
		Expect(cmdPlus.GetOutput()).To(Equal(customDir))
	})

	It("allows settings of the env variables", func() {
		cmdPlus := execplus.NewCmdPlus("./test_executables/print_env")
		cmdPlus.AppendEnv([]string{"MY_VAR=special"})
		err := cmdPlus.Run()
		Expect(err).To(BeNil())
		Expect(cmdPlus.GetOutput()).To(Equal("special"))
	})

	It("allows killing long running processes", func() {
		cmdPlus := execplus.NewCmdPlus("./test_executables/output_chunks")
		err := cmdPlus.Start()
		Expect(err).To(BeNil())
		time.Sleep(time.Second)
		err = cmdPlus.Kill()
		Expect(err).To(BeNil())
		Expect(cmdPlus.GetOutput()).NotTo(ContainSubstring("late chunk 4"))
	})

	It("allows waiting for long running processes", func() {
		cmdPlus := execplus.NewCmdPlus("./test_executables/output_chunks")
		err := cmdPlus.Start()
		Expect(err).To(BeNil())
		err = cmdPlus.Wait()
		Expect(err).To(BeNil())
		Expect(cmdPlus.GetOutput()).To(ContainSubstring("late chunk 4"))
	})

	Describe("output channel", func() {
		It("allows access to output chunks (separated by newlines) via a channel", func() {
			cmdPlus := execplus.NewCmdPlus("./test_executables/output_chunks")
			outputChannel, _ := cmdPlus.GetOutputChannel()
			err := cmdPlus.Start()
			Expect(err).To(BeNil())
			chunks := []execplus.OutputChunk{
				<-outputChannel,
				<-outputChannel,
				<-outputChannel,
				<-outputChannel,
				<-outputChannel,
			}
			Expect(chunks).To(Equal([]execplus.OutputChunk{
				{Chunk: "", Full: ""},
				{Chunk: "chunk 1", Full: "chunk 1"},
				{Chunk: "special chunk 2", Full: "chunk 1\nspecial chunk 2"},
				{Chunk: "chunk 3", Full: "chunk 1\nspecial chunk 2\nchunk 3"},
				{Chunk: "late chunk 4", Full: "chunk 1\nspecial chunk 2\nchunk 3\nlate chunk 4"},
			}))
			err = cmdPlus.Kill()
			Expect(err).To(BeNil())
		})

		It("sends the current status whenever the channel is added", func() {
			cmdPlus := execplus.NewCmdPlus("./test_executables/output_chunks")
			err := cmdPlus.Start()
			Expect(err).To(BeNil())
			time.Sleep(time.Second)
			outputChannel, _ := cmdPlus.GetOutputChannel()
			chunk := <-outputChannel
			Expect(chunk).To(Equal(execplus.OutputChunk{Chunk: "", Full: "chunk 1\nspecial chunk 2\nchunk 3"}))
			err = cmdPlus.Kill()
			Expect(err).To(BeNil())
		})

		It("works with stderr", func() {
			cmdPlus := execplus.NewCmdPlus("./test_executables/print_to_stderr")
			err := cmdPlus.Start()
			Expect(err).To(BeNil())
			outputChannel, _ := cmdPlus.GetOutputChannel()
			chunks := []execplus.OutputChunk{
				<-outputChannel,
				<-outputChannel,
				<-outputChannel,
				<-outputChannel,
			}
			Expect(chunks).To(Equal([]execplus.OutputChunk{
				{Chunk: "", Full: ""},
				{
					Chunk: "stdout chunk 1",
					Full:  "stdout chunk 1",
				},
				{
					Chunk: "stderr chunk",
					Full:  "stdout chunk 1\nstderr chunk",
				},
				{
					Chunk: "stdout chunk 2",
					Full:  "stdout chunk 1\nstderr chunk\nstdout chunk 2",
				},
			}))
			err = cmdPlus.Kill()
			Expect(err).To(BeNil())
		})
	})

	Describe("waitForCondition", func() {
		It("returns nil if the condition passes within the timeout", func() {
			cmdPlus := execplus.NewCmdPlus("./test_executables/output_chunks")
			err := cmdPlus.Start()
			Expect(err).To(BeNil())
			err = cmdPlus.WaitForCondition(func(chunk, full string) bool {
				return strings.Contains(chunk, "special")
			}, time.Second*2)
			Expect(err).To(BeNil())
			err = cmdPlus.Kill()
			Expect(err).To(BeNil())
		})

		It("returns error if the text is not seen within the timeout", func() {
			cmdPlus := execplus.NewCmdPlus("./test_executables/output_chunks")
			err := cmdPlus.Start()
			Expect(err).To(BeNil())
			err = cmdPlus.WaitForCondition(func(chunk, full string) bool {
				return strings.Contains(chunk, "other")
			}, time.Second*2)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Timed out after 2s, full output:\nchunk 1\nspecial chunk 2\nchunk 3"))
			err = cmdPlus.Kill()
			Expect(err).To(BeNil())
		})

		It("works for sequential waits", func() {
			cmdPlus := execplus.NewCmdPlus("./test_executables/output_chunks")
			err := cmdPlus.Start()
			Expect(err).To(BeNil())
			err = cmdPlus.WaitForCondition(func(chunk, full string) bool {
				return chunk == "chunk 1"
			}, time.Second*2)
			Expect(err).To(BeNil())
			err = cmdPlus.WaitForCondition(func(chunk, full string) bool {
				return chunk == "late chunk 4"
			}, time.Second*8)
			Expect(err).To(BeNil())
			err = cmdPlus.Kill()
			Expect(err).To(BeNil())
		})
	})

	Describe("waitForRegexp", func() {
		It("returns nil if the text is seen within the timeout", func() {
			cmdPlus := execplus.NewCmdPlus("./test_executables/output_chunks")
			err := cmdPlus.Start()
			Expect(err).To(BeNil())
			isChunk := regexp.MustCompile(`special chunk \d`)
			err = cmdPlus.WaitForRegexp(isChunk, time.Second*2)
			Expect(err).To(BeNil())
			err = cmdPlus.Kill()
			Expect(err).To(BeNil())
		})

		It("returns error if the text is not seen within the timeout", func() {
			cmdPlus := execplus.NewCmdPlus("./test_executables/output_chunks")
			err := cmdPlus.Start()
			Expect(err).To(BeNil())
			isChunk := regexp.MustCompile(`other chunk \d`)
			err = cmdPlus.WaitForRegexp(isChunk, time.Second*2)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Timed out after 2s, full output:\nchunk 1\nspecial chunk 2\nchunk 3"))
			err = cmdPlus.Kill()
			Expect(err).To(BeNil())
		})
	})

	Describe("waitForText", func() {
		It("returns nil if the text is seen within the timeout", func() {
			cmdPlus := execplus.NewCmdPlus("./test_executables/output_chunks")
			err := cmdPlus.Start()
			Expect(err).To(BeNil())
			err = cmdPlus.WaitForText("chunk 3", time.Second*2)
			Expect(err).To(BeNil())
			err = cmdPlus.Kill()
			Expect(err).To(BeNil())
		})

		It("returns error if the text is not seen within the timeout", func() {
			cmdPlus := execplus.NewCmdPlus("./test_executables/output_chunks")
			err := cmdPlus.Start()
			Expect(err).To(BeNil())
			err = cmdPlus.WaitForText("chunk 4", time.Second*2)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Timed out after 2s, full output:\nchunk 1\nspecial chunk 2\nchunk 3"))
			err = cmdPlus.Kill()
			Expect(err).To(BeNil())
		})

		It("works for prompts (text that ends with a colon followed by a space)", func() {
			cmdPlus := execplus.NewCmdPlus("./test_executables/prompt")
			err := cmdPlus.Start()
			Expect(err).To(BeNil())
			err = cmdPlus.WaitForText("prompt: ", time.Second*2)
			Expect(err).To(BeNil())
			err = cmdPlus.Kill()
			Expect(err).To(BeNil())
		})
	})
})
