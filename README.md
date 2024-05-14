# keexp

keexp is a small program for exporting data from KeePass databases to environment variables.

## Usage

* Create a config file based on `config_example.json`.
  * An entry's UUID can be viewed on its properties tab.
  * Custom fields can be used just the same as default fields like `Password` and `UserName`. Keep in mind that these are case-sensitive, though.
* Override the `keexp` command so the exported variables will be set in the shell you're running it in:

      # Bash
      keexp() { eval $(command keexp "$@"); }
      # fish
      function keexp; /usr/local/bin/keexp $argv | source; end

* Run keexp:

      keexp /path/to/database.kdbx /path/to/config.json
