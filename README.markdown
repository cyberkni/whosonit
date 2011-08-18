Who's On It
===========

Who's On It is a lightweight email tracking system written in Go for use on AppEngine. It utilizes the AppEngine Email services to receive and send email.

The workflow of the system is the following:

1. An email is received by the AppEngine incoming email feature and passed to this application.
2. The system records this message and sends a notice to a configured notification email address.
3. An interested person comes along to act on the email. To indicate this intent they click the "Accept" button which marks the email as owned by them.
4. The owner of the email comes back later after handling the email and marks the email as "Closed" which hides it from the list.

Each of these state changes will generate an optional notification email to notification address.

What Works?
-----------
* Index - email listing page
* Individual email view
* Accept event
* Close event
* Fake Email test form and action to handle form

What Is Left To Do?
-------------------
* User authentication - make sure only the right people will be allowed to see the emails
* Use real usernames as owners
* Fix int64 timestamps to look like real times
* Add email handler to process real incoming emails
* Make it prettier
* Email notifications