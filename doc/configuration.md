# OpenSAPS configuration

There is no hardcoded place for OpenSAPS configuration. You **should**
provide path to configuration file via ``-config`` parameter.

# Example configuration

Example can be viewed in opensaps.example.yaml, which is stored in root
directory of this repository.

# Configuration values.

Here we will go thru configuration values available. Nesting shows
nesting level in configuration file.

* ``slackhandler`` - namespace for configuring Slack API HTTP listener.

  * ``listener`` - namespace for configuring HTTP listener itself.

    * ``address`` - IP address and port we will listen on. Defaulting
    to ``127.0.0.1:39231``.

* ``webhooks`` - namespace for webhooks configuration. Here you should
define webhook name (**should be unique!**) and some parameters.

  * ``gitea_to_matrix`` - example webhook name. Should be unique and
  can be anything you can imagine (in text, of course).

  **WARNING:** multiline webhook names wasn't tested! Try to keep
  your text in single line!
    * ``slack`` - namespace for configuring Slack API parameters.
    URL for Slack webhook looks like:

    ```
    http(s)://server.tls/services/T12345678/B87654321/24charslongstring
    ```

    Where ``12345678`` is a random 8-char string (all caps) and
    ``24charslongstring`` is a random 24-char string.

    Next variables configures these strings.

      * ``random1`` - first 8-char random string (``T12345678``).

      * ``random2`` - second 8-char random string (``B87654321``).

      * ``longrandom`` - 24-char random string.

    * ``remote`` - configures pusher this webhook should use and
    connection name for it.

      * ``pusher`` - what pusher this webhook should use.

      * ``push_to`` - connection name for this pusher.

* ``matrix`` - configures Matrix pusher connections available.

  * ``matrix_test`` - connection name. Should be unique and
  can be anything you can imagine (in text, of course).

  **WARNING:** multiline webhook names wasn't tested! Try to keep
  your text in single line!
    * ``api_root`` - API root for Matrix connection. For example,
    ``https://localhost:8448/_matrix/client/r0``.

    * ``user`` - Matrix user.

    * ``password`` - password for Matrix user.

    * ``room`` - room ID to use. If Matrix user isn't in that room
    while OpenSAPS logging in - OpenSAPS will try to join this room.
