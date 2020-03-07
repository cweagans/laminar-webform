# laminar webform

This is a web UI to kick off [laminar](https://laminar.ohwg.net) jobs using specific
params.

documentation help is welcome.

## Configuration

Example config:

```toml
[general]
# this is used for the <title> element on the page
title = "laminar-webform"
# this should be the URL where you can actually see the laminar UI.
# used for redirects
laminar_url = "http://localhost:8080"
# this will output additional information in the logs
debug = true

# start by defining a new form block
[forms.bulk_action]
# the title and description will be used to describe the form
title = "Run bulk action"
description = "Perform an action against all sites on a specified environment."
# this is the job in laminar that this form should queue.
job = "hello"

# one or more fields can be added to a form. if you just want a "build now" button
# with no arguments, you don't need this.
[[forms.bulk_action.fields]]
# this is the field label
title = "Site environment"
# this is the actual name of the argument that you would normally pass to laminarc
name = "environment"
# a user-friendly description of the field.
description = "The environment that you want to run the action against. Actions against live need signoff from one other person."
# this can either be "select" or "text". "select" requires options as seen below.
type = "select"
# options will be shown in a select dropdown. do not use spaces or special characters here.
options = [
  "dev",
  "test",
  "earlylive",
  "live",
]

[[forms.bulk_action.fields]]
title = "Action"
name = "action"
description = "The action that you want to perform."
# text fields are single line text entry boxes
type = "text"
# the regex defined here must match the contents of the text field or the form
# will not be considered valid.
filter = "^[a-zA-Z0-8]+$"

```


## Building

`make`

## Releasing

Push a tag to the repo and Github CI will take over, build a release, and upload
the artifacts to the release page.
