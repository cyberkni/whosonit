package whosonit

import (
    "template"
    )

const root_template = `
<html>
<body>
<h1>Who is on it?</h1>
<h2>Current Emails</h2>
<ul>
{.repeated section @}
  <li><a href="/show?sender={Sender}&date={RecieptDate}">{Sender} - {Subject} - {RecieptDate}</a></li>
{.or}
Nothing to do. Good job.
{.end}
</ul>
</body>
</html>
`

var rootTemplate = template.MustParse(root_template, nil)

const show_template = `
<html>
<body>
{.repeated section @}
<h1>{Subject|html}</h1>
<h2>{Sender}</h2>
<h3>Status</h3>
<p>Received: {RecieptDate}</p>
<p>Owner: {.section Owner}{Owner}</p>
<p>Accepted On: {OwnerDate}{.or}Nobody{.end}</p>
{.section ClosedDate}
<p>Closed On: {ClosedDate}</p>
{.end}
<h3>Message Content</h3>
<p>{Body}</p>
{.section Owner}
{.section ClosedDate}
<h3>Actions</h3>
<form action="/close">
<input type="hidden" name="sender" value="{Sender}">
<input type="hidden" name="date" value="{RecieptDate}">
<input type="submit" name="Close" value="Close">
</form>
{.or}
{.end}
{.or}
<h3>Actions</h3>
<form action="/accept">
<input type="hidden" name="sender" value="{Sender}">
<input type="hidden" name="date" value="{RecieptDate}">
<input type="submit" name="Accept" value="Accept">
</form>
{.end}
{.end}
</body>
</html>`
var showTemplate = template.MustParse(show_template, nil)

const test_form = `
<html>
<body>
<form action="/test_form">
<p>Sender: <input name="Sender"></p>
<p>Message: <input name="Body"></p>
<p><input type="submit"></p>
</body>
</html>`


