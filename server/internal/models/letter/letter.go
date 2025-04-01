package letter

import (
	"fmt"
)

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

func NewReminderNotification(name, title, location, time string) string {
	text := fmt.Sprintf(`
		<div> 
			<h1> Hello, dear %s </h1>
			<p> We remind you of an upcoming event </p>
			<h1> %s </h1>
			<h3> The event will take place at <b> %s </b></h3>
			<h3> The time of the event is <b> %s </b></h3>
		</div>
		`,
		name, 
		title,
		location,
		time,
	)

	return text
}
