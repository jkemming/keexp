# Usage

* Copy `config_example.json` to `~/.config/keexp/config.json` (or `${XDG_CONFIG_HOME}/keexp/config.json`) and adapt it to your database.
	* An entry's UUID can be viewed on its properties tab.
	* Custom fields can be used just the same as default fields like `Password` and `UserName`. Keep in mind that these are case-sensitive, though.
* Set an alias so the exported variables will actually be set in the shell you're running `keexp` in:

		```shell
		# Bash
		alias keexp='keexp | source /dev/stdin'
		# fish
		alias keexp 'keexp | source'
		```
