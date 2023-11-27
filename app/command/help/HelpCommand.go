package help

type Command struct{}

func New() Command {
	return Command{}
}

func (c Command) Run() (string, error) {
	return `
	Hello baby! I'll tell you how to build and deliver your code.
	
	Firstly, you should create release at jira with your tasks.
	Next, tell the build bot to build your release and return release candidate tag:
    <b>/build#{RELEASE ID}</b>

	Use release candidate tag to trigger image building:
    <b>/image build {APP}#{RC TAG}</b>
	Next, you can use image to deploy it on stage environment:
    <b>/deploy stage {APP}#{RC TAG}</b>

	When testing is finished, complete the release build by placing the final tag in the repository:
	<b>/build#{RELEASE ID}</b>
	Also put the final tag to image register:
	<b>/image release {APP}#{RC TAG}</b>

	Last step, deliver the image to the production environment:
	<b>/deploy prod {APP}#{TAG}</b>
	`, nil
}
