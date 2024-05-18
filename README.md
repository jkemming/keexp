# keexp

keexp is a small program for exporting data from KeePass databases to environment variables.

## Usage

* Create a config file based on `config_example.json`.
  * An entry's UUID can be viewed on its properties tab.
  * Custom fields can be used just the same as default fields like `Password` and `UserName`. Keep in mind that these are case-sensitive, though.
* Run `keexp` and `eval` the output so exported variables will be set in the shell:
  ```shell
  eval $(/path/to/keexp /path/to/database.kdbx /path/to/config.json)
  ```
