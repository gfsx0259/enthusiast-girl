package help

type Command struct{}

func (c Command) String() string {
	return ""
}

func New() Command {
	return Command{}
}

func (c Command) Run() (string, error) {
	return `
	Here's how to build and deliver your code:
	Start by creating a release in Jira with your associated tasks.

	Tell the build bot to build your release and return the release candidate tag using the following command:
    <b>/build#{RELEASE ID}</b>

	Use the release candidate tag to trigger the image building process with this command:
    <b>/image build {APP}#{RC TAG}</b>

	Once the image is built, you can deploy it to the stage environment using this command: /deploy stage {APP} {RC TAG}
    <b>/deploy stage {APP}#{RC TAG}</b>

	After testing is complete, finalize the release build by pushing the final tag to the repository with this command: /build {RELEASE ID}
	<b>/build#{RELEASE ID}</b>

	Additionally, update the image registry with the final tag using this command:
	<b>/image release {APP}#{RC TAG}</b>

	Finally, deliver the image to the production environment with this command:
	<b>/deploy prod {APP}#{TAG}</b>
	`, nil
}
