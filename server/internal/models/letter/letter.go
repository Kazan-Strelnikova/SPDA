package letter

import "fmt"

type Change struct {
	Name       string
	OldContent string
	NewContent string
}

func NewChange(name, old, new string) Change {
	return Change{
		Name:       name,
		OldContent: old,
		NewContent: new,
	}
}

func NewUpdateNotification(name, title string, changes []Change) string {
	text := fmt.Sprintf(`
	<div>
		<h1> Hello, dear %s </h1>
		<p> Details of the event %s have been changed </p>
		<ul>
	`,
		name,
		title,
	)

	for _, change := range changes {
		text += fmt.Sprintf("\n\t\t<li>%s changed from:<br/>%s<br/>to:<br/>%s<br/></li>",
			change.Name,
			change.OldContent,
			change.NewContent,
		)
	}

	text += "\n\t</ul>\n</div>"

	return text
}
