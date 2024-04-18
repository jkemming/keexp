# Usage

* Create a config file based on `config_example.json`.
	* An entry's UUID can be viewed on its properties tab.
	* Custom fields can be used just the same as default fields like `Password` and `UserName`. Keep in mind that these are case-sensitive, though.
* Override the `keexp` command so the exported variables will be set in the shell you're running it in:

		# Bash
		keexp() { keexp "$1" "$2" | source /dev/stdin; }
		# fish
		function keexp; /usr/local/bin/keexp $argv[1] $argv[2] | source; end

* Run keexp:

		keexp /path/to/database.kdbx /path/to/config.json
