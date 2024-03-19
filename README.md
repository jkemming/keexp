# Usage

* Create a config file based on `config_example.json`.
	* An entry's UUID can be viewed on its properties tab.
	* Custom fields can be used just the same as default fields like `Password` and `UserName`. Keep in mind that these are case-sensitive, though.
* Move your config file to a location where keexp will find it:
	* If environment variable `XDG_CONFIG_HOME` is set, the path is `${XDG_CONFIG_HOME}/keexp/config.json`
	* By default, the path is `~/.config/keexp/config.json`
* Set an alias so the exported variables will actually be set in the shell you're running `keexp` in:

		```shell
		# Bash
		alias keexp='keexp | source /dev/stdin'
		# fish
		alias keexp '/usr/local/bin/keexp | source'
		```
